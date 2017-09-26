# MercadoPago CoreServices Server

MPCS server initialization and Fury scope parsing abstraction. 

## Concept

MPCS uses Fury scopes for defining which functionality of an application a given scaling group should respond to.

For example, given a fury scope `production-search` the application should bootstrap itself so that it only accept requests pertaining `search` functionality and respond `404` for every other endpoint; using server defaults for a `production` environment.

The Fury scope format used by MPCS is: `environment`-`role`-`tag` (the `tag` segment is optional).

This package provides the functionality needed for initializing a web server with sane settings based on a given Fury scope.

The following table show the default server configuration  for each valid environment:

| Environment | Log Level | Push Metrics | Debug | Auth Scopes |
|:-:|:-:|:-:|:-:|:-:|
| Develop | DEBUG | ✔ | ✔ | [] |
| Test | INFO | ✔ | ✔ | [] |
| Integration | INFO | ❌ | ✔ | [] |
| Production | INFO | ✔ | ❌ | [] |

> *NOTE: Every option listed above can be overridden when creating a new server, by using the With\* opt functions provided by the package.*

## Installation

```bash
go get github.com/mercadolibre/coreservices-team/libs/go/server
```

## Usage

This package is though for using when bootstraping the `main` package of a new application.

```go
package main

import (
    "github.com/mercadolibre/coreservices-team/libs/go/server"
)

routes := server.RoutingGroup{
    server.RoleIndexer: func(g *gin.RouterGroup) {
        g.POST("/indexer", func(c *gin.Context) {})
    },
    server.RoleWrite: func(g *gin.RouterGroup) {
        g.POST("/writer", func(c *gin.Context) {})
    },
}

srv, err := server.NewEngine("test-indexer-feature-branch", routes, server.WithDebug(false), server.WithPushMetrics(true))
if err != nil {
    log.Fatal(err)
}

srv.Run(":8080")
```

Alternatively you can use a controller with a method that returns a function with signature `func(*gin.RouterGroup)` and then use that to register the appropriate routes. 

In effect you'll be delegating the knowledge of the routes to the module, instead of bundling it with your main program.

```go
searchCtrl := search.NewController(db)

routes := server.RoutingGroup{
    server.RoleSearch: searchCtrl.RegisterRoutes(),
    // ...
}

// ...
```

## Changelog

####  2017-09-25:
First release.

## Author
Core Services Team 
- `coreservices@mercadolibre.com`
- [Slack](https://mercadopago-team.slack.com/messages/C45S2LB5K)
