package server

import (
	"fmt"

	"github.com/atarantini/ginrequestid"
	"github.com/gin-gonic/gin"
	"github.com/mercadolibre/go-meli-toolkit/gingonic/mlhandlers"
)

// MercadoPagoCoreServicesGroupPreffix is the preffix added to every exposed url
const MercadoPagoCoreServicesGroupPreffix = "/mpcs"

// settings contains the relevant information that our server may use for initiating.
type settings struct {
	LogLevel    string
	PushMetrics bool
	Debug       bool
	AuthScopes  []string
}

// Map with default server settings for each possible scope.
var envSettings = map[Environment]settings{
	EnvDevelop:     {LogLevel: "DEBUG", PushMetrics: true, Debug: true, AuthScopes: []string{}},
	EnvTest:        {LogLevel: "INFO", PushMetrics: true, Debug: true, AuthScopes: []string{}},
	EnvIntegration: {LogLevel: "INFO", PushMetrics: false, Debug: true, AuthScopes: []string{}},
	EnvProduction:  {LogLevel: "INFO", PushMetrics: true, Debug: false, AuthScopes: []string{}},
}

// RoutingGroup is a map of urls and functions for a given role.
type RoutingGroup map[Role]func(*gin.RouterGroup)

// Opt is a function for Server, it's used for optional modifiers used in package constructor
type Opt func(*Server)

// Server is our application main struct, it's basically a wrapper around a
// gin.Engine instance, with some functionality hidden for easier usage.
type Server struct {
	*gin.Engine
	Context ApplicationContext

	settings settings
}

// NewEngine configures the underlying gin.Engine struct of Server with a given fury scope, a
// RoutingGroup (exposed urls mapped to a valid Role) and accepts a list of options for
// specifying configuration options outside of the defaults for a given environment.
func NewEngine(scope string, routes RoutingGroup, opts ...Opt) (*Server, error) {
	// Infer application context from fury scope
	ctx, err := ContextFromScopeString(scope)
	if err != nil {
		return nil, fmt.Errorf("error infering context from fury scope: %v", err)
	}

	// Check if the given routes are valid for the current application role
	if _, ok := routes[ctx.Role]; !ok {
		return nil, fmt.Errorf("given routes do not contain endpoints for the current application role")
	}

	// Create server with default configuration for current environment
	server := &Server{
		Context: ctx,

		settings: envSettings[ctx.Environment],
	}

	// Call option functions on instance before instantiating the server so that custom
	// options are taken into consideration.
	for _, opt := range opts {
		opt(server)
	}

	// Create a gin engine with debug or release config depending on the given settings
	if server.settings.Debug {
		gin.SetMode(gin.DebugMode)
		server.Engine = gin.Default()
	} else {
		gin.SetMode(gin.ReleaseMode)
		server.Engine = gin.New()
	}

	// Global server configuration, common to all environments
	server.NoRoute(NoRouteHandler())
	server.RedirectFixedPath = false
	server.RedirectTrailingSlash = false

	// Setup health check handler
	server.GET("/ping", HealthCheckHandler)

	// Call the current Role group function with the current group as param
	// so that it loads the active urls.
	group := server.Group(MercadoPagoCoreServicesGroupPreffix)

	group.Use(ginrequestid.RequestId())

	// Add authentication middleware, but only if not on indexer role.
	// When on indexer role, requests are being called by BigQ, and in this
	// scenario there's not authentication present (no caller ID nor scopes).
	// If this middleware is run under this conditions, all requests would fail.
	if ctx.Role != RoleIndexer {
		group.Use(Auth())
	}

	if server.settings.PushMetrics {
		group.Use(mlhandlers.NewRelic())
		group.Use(mlhandlers.Datadog())
		group.Use(RenameNewRelicTransaction())
	}

	if auth := server.settings.AuthScopes; len(auth) > 0 {
		group.Use(mlhandlers.MLAuth(auth))
	}

	fn := routes[ctx.Role]
	fn(group)

	return server, nil
}

// WithAuthScopes func sets up required application authentication scopes.
func WithAuthScopes(authScopes []string) Opt {
	return func(s *Server) {
		s.settings.AuthScopes = authScopes
	}
}

// WithLogLevel func sets application log level.
func WithLogLevel(logLevel string) Opt {
	return func(s *Server) {
		s.settings.LogLevel = logLevel
	}
}

// WithPushMetrics func sets push metrics options.
func WithPushMetrics(pushMetrics bool) Opt {
	return func(s *Server) {
		s.settings.PushMetrics = pushMetrics
	}
}

// WithDebug func sets debugging for webserver.
func WithDebug(debug bool) Opt {
	return func(s *Server) {
		s.settings.Debug = debug
	}
}
