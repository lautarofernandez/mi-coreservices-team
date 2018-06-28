package services

import (
	"database/sql"
	"fmt"
	"net/url"
	"time"

	// The mysql driver needs to be initialized implicitly so that it
	// can hook itself into the database/sql API.
	_ "github.com/go-sql-driver/mysql"

	"github.com/mercadolibre/coreservices-team/libs/go/server"
	bq "github.com/mercadolibre/go-meli-toolkit/gobigqueue"
	ds "github.com/mercadolibre/go-meli-toolkit/godsclient"
	kvs "github.com/mercadolibre/go-meli-toolkit/gokvsclient"
	cache "github.com/mercadolibre/go-meli-toolkit/gomemcached"
)

// Services ...
type Services struct {
	ctx      server.ApplicationContext
	services map[string]service
}

// New parses a configuration file and returns a Services struct.
func New(ctx server.ApplicationContext) (*Services, error) {
	s, err := parseYAML(ctx.Environment)
	if err != nil {
		return nil, err
	}

	return &Services{ctx, s}, nil
}

func (s *Services) service(name string) (service, error) {
	svc, exists := s.services[name]
	if !exists {
		return svc, fmt.Errorf("service %s not found", name)
	}

	return svc, nil
}

// KVS returns and initializes KVS client with the correct configuration for
// the given environment, or error if something goes wrong.
func (s *Services) KVS(name string, config kvs.KvsClientConfig) (kvs.Client, error) {
	svc, err := s.service(name)
	if err != nil {
		return nil, err
	}

	if !svc.HasRole(s.ctx.Role) {
		return nil, nil
	}

	if svc.Type != TypeKVS {
		return nil, fmt.Errorf("service %s is of type %s, not KVS", name, svc.Type)
	}

	if config == nil {
		config = kvs.MakeKvsConfig()
		config.SetReadMaxIdleConnections(50)
		config.SetWriteMaxIdleConnections(50)
		config.SetReadTimeout(300 * time.Millisecond)
		config.SetWriteTimeout(300 * time.Millisecond)
	}

	// If a service name is given as part of the KVS config, then use that to
	// initialize the KVS client.
	if svcName, ok := svc.SvcParams["service"]; ok {
		return kvs.MakeKvsClient(svcName, config), nil
	}

	// If there's no service name, then we'll try to initialize the KVS using the read and write endpoints
	if mapContains(svc.SvcParams, "endpoint_read", "endpoint_write", "container_name") {
		config.SetContainerName(svc.SvcParams["container_name"])
		config.SetReadEndpoint(svc.SvcParams["endpoint_read"])
		config.SetWriteEndpoint(svc.SvcParams["endpoint_write"])

		return kvs.MakeKvsClient(svc.SvcParams["container_name"], config), nil
	}

	return nil, fmt.Errorf("missing params for initializing KVS container")
}

// DS returns and initializes a DS client with the correct configuration for
// the given environment, or error if something goes wrong.
func (s *Services) DS(name string, config *ds.DsClientConfig) (ds.Client, error) {
	svc, err := s.service(name)
	if err != nil {
		return nil, err
	}

	if !svc.HasRole(s.ctx.Role) {
		return nil, nil
	}

	if svc.Type != TypeDS {
		return nil, fmt.Errorf("service %s is of type %s, not DS", name, svc.Type)
	}

	if config == nil {
		config = ds.NewDsClientConfig()
	}

	// If a service name is given as part of the KVS config, then use that to
	// initialize the KVS client.
	if svcName, ok := svc.SvcParams["service"]; ok {
		config = config.WithServiceName(svcName)

		return ds.NewEntityClient(config), nil
	}

	// If not, we'll have to manually initialize it using the read and write endpoints
	if mapContains(svc.SvcParams, "namespace", "entity", "read_endpoint", "write_endpoint") {
		config = config.
			WithNamespace(svc.SvcParams["namespace"]).
			WithEntity(svc.SvcParams["entity"]).
			WithReadEndpoint(svc.SvcParams["read_endpoint"]).
			WithWriteEndpoint(svc.SvcParams["write_endpoint"])

		return ds.NewEntityClient(config), nil
	}

	return nil, fmt.Errorf("missing params for initializing DS container")
}

// Publisher returns and initializes a BigQ publisher with the correct configuration
// for the given environment, or error if something goes wrong.
func (s *Services) Publisher(name string) (bq.Publisher, error) {
	svc, err := s.service(name)
	if err != nil {
		return nil, err
	}

	if !svc.HasRole(s.ctx.Role) {
		return nil, nil
	}

	if svc.Type != TypeQueueTopic {
		return nil, fmt.Errorf("service %s is of type %s, not Topic", name, svc.Type)
	}

	if mapContains(svc.SvcParams, "topic", "cluster") {
		return bq.NewPublisher(svc.SvcParams["cluster"], []string{svc.SvcParams["topic"]}), nil
	}

	return nil, fmt.Errorf("missing params for initializing BigQ publisher")

}

// Cache returns and initializes a DS client with the correct configuration for
// the given environment, or error if something goes wrong.
func (s *Services) Cache(name string) (cache.Client, error) {
	svc, err := s.service(name)
	if err != nil {
		return nil, err
	}

	if !svc.HasRole(s.ctx.Role) {
		return nil, nil
	}

	if svc.Type != TypeCache {
		return nil, fmt.Errorf("service %s is of type %s, not Cache", name, svc.Type)
	}

	if mapContains(svc.SvcParams, "endpoint") {
		cache.RegisterCluster(name, svc.SvcParams["endpoint"])

		return cache.NewClient(name)
	}

	return nil, fmt.Errorf("missing params for initializing Cache service")
}

// DB returns and initializes a SQL DB client with the correct configuration for
// the given environment, or error if something goes wrong.
func (s *Services) DB(name string) (*sql.DB, error) {
	svc, err := s.service(name)
	if err != nil {
		return nil, err
	}

	if !svc.HasRole(s.ctx.Role) {
		return nil, nil
	}

	if svc.Type != TypeDatabase {
		return nil, fmt.Errorf("service %s is of type %s, not Database", name, svc.Type)
	}

	if mapContains(svc.SvcParams, "database", "username", "password", "host") {
		qs := url.Values{}
		qs.Add("collation", "utf8mb4_general_ci")
		qs.Add("parseTime", "true")
		qs.Add("timeout", "1.5s")
		qs.Add("readTimeout", "1s")
		qs.Add("writeTimeout", "1s")

		connectionString := fmt.Sprintf("%s:%s@tcp(%s)/%s?%s", svc.SvcParams["username"], svc.SvcParams["password"], svc.SvcParams["host"], svc.SvcParams["database"], qs.Encode())
		db, err := sql.Open("mysql", connectionString)
		if err != nil {
			return nil, err
		}

		return db, nil
	}

	return nil, fmt.Errorf("missing params for initializing Database service")
}

func mapContains(m map[string]string, keys ...string) bool {
	for _, k := range keys {
		if _, ok := m[k]; !ok {
			return false
		}
	}

	return true
}