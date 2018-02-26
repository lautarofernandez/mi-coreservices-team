# Gordik

This package contains parts of the work in progress Gordik Framework.

## Context

In order to use `gk.Context` package, you should define your HTTP handlers as follow:

```go
func ControllerHandler(c *gin.Context, ctx *gk.Context) {
    // ...
}
```

When hooking tour handler to GIN you must encapsulate it in `gk.Handler` function.

```go
v1 := g.Group("/v1")

v1.POST("/", gk.Handler(ControllerHandler))
```

The gordik context gives you typed access to authentication information, request data, and a configured logger that will automatically add the `request_id` to any logged lines.

You'll also find helper methods for creating NewRelic segments and measuring database operations.

## Services

> :warning: This package is WIP and should be used carefully.

Services is a library that let's you use a configuration file for defining Fury services, and setting different environment values depending on the SCOPE in which the application bootstrapped. The library handles the initialization of each service using `go-meli-toolkit` SDK.
