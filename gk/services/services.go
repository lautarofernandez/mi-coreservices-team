package services

import (
	"database/sql"
	"fmt"
	"net/url"
	"strings"
	"time"

	// The mysql driver needs to be initialized implicitly so that it
	// can hook itself into the database/sql API.
	_ "github.com/go-sql-driver/mysql"

	"github.com/mercadolibre/coreservices-team/libs/go/server"
	bq "github.com/mercadolibre/go-meli-toolkit/gobigqueue"
	ds "github.com/mercadolibre/go-meli-toolkit/godsclient"
	kvs "github.com/mercadolibre/go-meli-toolkit/gokvsclient"
	lock "github.com/mercadolibre/go-meli-toolkit/golockclient"
	cache "github.com/mercadolibre/go-meli-toolkit/gomemcached"
	os "github.com/mercadolibre/go-meli-toolkit/goosclient"
)

// Services ...
type Services struct {
	ctx      server.ApplicationContext
	services map[string]service
}

// NewWithFile parses the given configuration file and returns a Services struct.
func NewWithFile(file string, ctx server.ApplicationContext) (*Services, error) {
	s, err := parseYAML(file, ctx.Environment)
	if err != nil {
		return nil, err
	}

	return &Services{ctx, s}, nil
}

// New parses a configuration file and returns a Services struct.
func New(ctx server.ApplicationContext) (*Services, error) {
	return NewWithFile("config.yml", ctx)
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

// Lock returns and initializes a lock client with the correct configuration for
// the given environment, or error if something goes wrong.
func (s *Services) Lock(name string, config lock.LockClientConfig) (lock.Client, error) {
	svc, err := s.service(name)
	if err != nil {
		return nil, err
	}

	if !svc.HasRole(s.ctx.Role) {
		return nil, nil
	}

	if svc.Type != TypeLock {
		return nil, fmt.Errorf("service %s is of type %s, not KVS", name, svc.Type)
	}

	// If config is not given then use default config values.
	if config == nil {
		config = lock.MakeLockClientConfig()
	}

	// If a service name is given as part of the KVS config, then use that to
	// initialize the KVS client.
	if svcName, ok := svc.SvcParams["service"]; ok {
		return lock.MakeLockClient(svcName, config), nil
	}

	return nil, fmt.Errorf("missing params for initializing lock namespace")
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

// OS returns and initializes a Object Storage client with the correct configuration for
// the given environment, or error if something goes wrong.
func (s *Services) OS(name string, configRead, configWrite os.OsClientConfig) (os.Client, error) {
	svc, err := s.service(name)
	if err != nil {
		return nil, err
	}

	if !svc.HasRole(s.ctx.Role) {
		return nil, nil
	}

	if svc.Type != TypeObjectStorage {
		return nil, fmt.Errorf("service %s is of type %s, not Object Storage", name, svc.Type)
	}

	if configRead == nil {
		configRead = os.MakeOSClientConfigRead()
	}

	if configWrite == nil {
		configWrite = os.MakeOSClientConfigRead()
	}

	if mapContains(svc.SvcParams, "service") {
		return os.MakeOsClient(svc.SvcParams["service"], configRead, configWrite), nil
	}

	return nil, fmt.Errorf("missing params for initializing object storage container")
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

// Cache returns and initializes a memcached client with the correct configuration for
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

	if !mapContains(svc.SvcParams, "endpoints") {
		return nil, fmt.Errorf("missing params for initializing Cache service")
	}

	servers := strings.Split(svc.SvcParams["endpoints"], " ")
	if len(servers) < 1 {
		return nil, fmt.Errorf("no servers found for cache service")
	}

	cache.RegisterCluster(name, servers...)

	return cache.NewClient(name)
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
