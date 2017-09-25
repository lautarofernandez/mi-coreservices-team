package server_test

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/mercadolibre/coreservices-team/libs/go/server"
)

func ExampleNewEngine() {
	routes := server.RoutingGroup{
		server.RoleIndexer: func(g *gin.RouterGroup) {
			g.POST("/indexer", func(c *gin.Context) {})
		},
		server.RoleWrite: func(g *gin.RouterGroup) {
			g.POST("/writer", func(c *gin.Context) {})
		},
	}

	srv, err := server.NewEngine(
		"test-indexer-feature-branch",
		routes,

		// Optional configuration override
		server.WithDebug(false),
		server.WithPushMetrics(true),
	)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Env: %s, Role: %s, Tag: %s", srv.Context.Environment, srv.Context.Role, srv.Context.Tag)

	//
	// srv.Run(":8080")

	//output: Env: test, Role: indexer, Tag: feature-branch
}
