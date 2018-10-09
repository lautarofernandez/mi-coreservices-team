package services

import (
	"fmt"
	"io/ioutil"

	"github.com/mercadolibre/coreservices-team/libs/go/server"
	"github.com/mercadolibre/go-meli-toolkit/gomelipass"
	yaml "gopkg.in/yaml.v2"
)

// serviceType is the type used for allowed config service names
type serviceType string

const (
	// TypeCache ...
	TypeCache serviceType = "cache"

	// TypeDatabase ...
	TypeDatabase serviceType = "database"

	// TypeQueueTopic ...
	TypeQueueTopic serviceType = "topic"

	// TypeKVS ...
	TypeKVS serviceType = "kvs"

	// TypeDS ...
	TypeDS serviceType = "ds"

	// TypeLock ...
	TypeLock serviceType = "lock"

	// TypeObjectStorage ...
	TypeObjectStorage serviceType = "storage"
)

var validServices = []serviceType{
	TypeCache,
	TypeDatabase,
	TypeQueueTopic,
	TypeKVS,
	TypeDS,
	TypeLock,
	TypeObjectStorage,
}

// service struct represents a Fury service with a given name and a
// given configuration for different application environments.
type service struct {
	Name      string
	Type      serviceType
	Roles     []string
	SvcParams map[string]string
}

// HasRole returns whether the given service is enabled for a specific role or not.
func (s service) HasRole(role server.Role) bool {
	for _, r := range s.Roles {
		if r == string(role) {
			return true
		}
	}

	return false
}

func parseYAML(filename string, environment server.Environment) (services map[string]service, err error) {
	services = map[string]service{}

	// Add default panic handler. Given that we cast a lot of interface{}
	// types to ones we know we should have, if the user gives us a
	// type that's not what we expected then a panic will occur.
	defer func() {
		if r := recover(); r != nil {
			var ok bool
			err, ok = r.(error)
			if !ok {
				err = fmt.Errorf("panic parsing config.yml")
			}
		}
	}()

	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("error reading contents of config.yml: %v", err)
	}

	m := map[interface{}]interface{}{}
	if err := yaml.Unmarshal(b, &m); err != nil {
		return nil, fmt.Errorf("error unmarshalling contents of config.yml: %v", err)
	}

	// Check that the root "services" key is present
	s, ok := m["services"]
	if !ok {
		return nil, fmt.Errorf("unable to find root services key")
	}

	// Iterate through each service, and parse it to a concrete Service struct. If
	// any expected convertion fails, then a panic will be risen and catched.
	for _, v := range s.([]interface{}) {
		svc := v.(map[interface{}]interface{})
		s := service{
			Name:      svc["name"].(string),
			Type:      serviceType(svc["type"].(string)),
			Roles:     []string{},
			SvcParams: map[string]string{},
		}

		for k, v := range svc {
			if k := k.(string); k == "name" || k == "type" {
				continue
			}

			if k := k.(string); k == "roles" {
				roles := v.([]interface{})
				for _, role := range roles {
					s.Roles = append(s.Roles, role.(string))
				}
			}

			// Given environment is not the current one, discard it.
			if k.(string) != string(environment) {
				continue
			}

			for k, v := range v.(map[interface{}]interface{}) {
				varName := k.(string)
				varValue := v.(string)

				s.SvcParams[varName] = varValue

				// The given param value might be a global variable that Fury will
				// inject with the correct value. We try to read from the key from
				// the env, and if something is found we replace it in the map.
				if envValue := gomelipass.GetEnv(varValue); envValue != "" {
					s.SvcParams[varName] = envValue
				}
			}
		}

		// Check that the current parsed service has at least 1 param for the given environment
		if len(s.SvcParams) == 0 {
			return nil, fmt.Errorf("service %s has 0 parameters for environment %s", s.Name, environment)
		}

		// Check whether the given service name was already loaded. If this is the case
		// then we fail, is better to fail for repeated services, that it is to
		// replace them and have the user see unexpected issues.
		if _, ok := services[s.Name]; ok {
			return nil, fmt.Errorf("duplicated service name found")
		}

		services[s.Name] = s
	}

	return services, nil
}
