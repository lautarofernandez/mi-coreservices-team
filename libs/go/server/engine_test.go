package server

import (
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
)

func init() {
	// We need to disable logging on this package because gin outputs debug
	// information using log that makes it hard to understand tests output.
	log.SetOutput(ioutil.Discard)
}

func okHandler(c *gin.Context) {
	c.String(http.StatusOK, "OK")
}

var routes = RoutingGroup{
	RoleIndexer: func(g *gin.RouterGroup) {
		g.GET("/test-indexer", okHandler)
	},
	RoleRead: func(g *gin.RouterGroup) {
		g.GET("/test-read", okHandler)
	},
	RoleSearch: func(g *gin.RouterGroup) {
		g.GET("/test-search", okHandler)
	},
	RoleWrite: func(g *gin.RouterGroup) {
		g.GET("/test-write", okHandler)
	},
	RoleMiddleEnd: func(g *gin.RouterGroup) {
		g.GET("/test-middleend", okHandler)
	},
}

func TestNewEngine(t *testing.T) {
	tt := []struct {
		Name              string
		Scope             string
		ExpectedEndpoints []string
	}{
		{"Indexer Role", "test-indexer", []string{"/ping", "/mpcs/test-indexer"}},
		{"Read Role", "test-read", []string{"/ping", "/mpcs/test-read"}},
		{"Search Role", "test-search", []string{"/ping", "/mpcs/test-search"}},
		{"Write Role", "test-write", []string{"/ping", "/mpcs/test-write"}},
		{"Middleend Role", "test-middleend", []string{"/ping", "/mpcs/test-middleend"}},
	}

	for _, tc := range tt {
		t.Run(tc.Scope, func(t *testing.T) {
			s, err := NewEngine(tc.Scope, routes)

			if err != nil {
				t.Fatalf("Error not expected, received: %v", err)
			}

			r := s.Engine.Routes()
			if len(r) != len(tc.ExpectedEndpoints) {
				t.Fatalf("Expected engine to have %v rules, got: %v", len(tc.ExpectedEndpoints), len(r))
			}

			for _, expected := range tc.ExpectedEndpoints {
				func() {
					for _, route := range r {
						if expected == route.Path {
							return
						}
					}

					t.Fatalf("Could not find route %v in router", expected)
				}()
			}
		})
	}
}

func TestNewEngine_InvalidScope(t *testing.T) {
	_, err := NewEngine("invalid-scope", routes)

	if err == nil {
		t.Fatal("Expected an error when initializing server")
	}
}

func TestNewEngine_NoRoutes(t *testing.T) {
	_, err := NewEngine("test-indexer", RoutingGroup{})

	if err == nil {
		t.Fatal("Expected an error when initializing server")
	}
}

func TestNewEngine_ServerMode(t *testing.T) {
	tt := []struct {
		Name     string
		Debug    bool
		Expected string
	}{
		{"Release mode", false, gin.ReleaseMode},
		{"Debug mode", true, gin.DebugMode},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			_, err := NewEngine("test-indexer", routes, WithDebug(tc.Debug))
			if err != nil {
				t.Fatalf("Error not expected, received: %v", err)
			}

			if gin.Mode() != tc.Expected {
				t.Fatalf("Expected gin mode to be %v, got: %v", tc.Expected, gin.Mode())
			}
		})
	}
}

func TestWithAuthScopes(t *testing.T) {
	scopes := []string{"admin", "write"}

	s, err := NewEngine("test-indexer", routes, WithAuthScopes(scopes))
	if err != nil {
		t.Fatalf("Error not expected, received: %v", err)
	}

	if len(scopes) != len(s.settings.AuthScopes) {
		t.Fatalf("Expected %v scopes, got: %v", len(scopes), len(s.settings.AuthScopes))
	}

	if !reflect.DeepEqual(scopes, s.settings.AuthScopes) {
		t.Fatal("Expected auth scopes and scopes in server differ.")
	}
}

func TestWithLogLevel(t *testing.T) {
	tt := []struct {
		Name     string
		LogLevel string
	}{
		{"DEBUG mode", "DEBUG"},
		{"INFO mode", "INFO"},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			s, err := NewEngine("test-indexer", routes, WithLogLevel(tc.LogLevel))
			if err != nil {
				t.Fatalf("Error not expected, received: %v", err)
			}

			if s.settings.LogLevel != tc.LogLevel {
				t.Fatalf("Expected log level mode to be %v, got: %v", tc.LogLevel, s.settings.LogLevel)
			}
		})
	}
}

func TestWithPushMetrics(t *testing.T) {
	tt := []struct {
		Name        string
		PushMetrics bool
	}{
		{"PushMetrics ON", true},
		{"PushMetrics OFF", false},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			s, err := NewEngine("test-indexer", routes, WithPushMetrics(tc.PushMetrics))
			if err != nil {
				t.Fatalf("Error not expected, received: %v", err)
			}

			if s.settings.PushMetrics != tc.PushMetrics {
				t.Fatalf("Expected log level mode to be %v, got: %v", tc.PushMetrics, s.settings.PushMetrics)
			}
		})
	}
}
