package middleware

import (
	"github.com/giantswarm/micrologger"

	"github.com/giantswarm/auto-oncall/service"
)

// Config represents the configuration used to create a middleware.
type Config struct {
	// Dependencies.
	Logger  micrologger.Logger
	Service *service.Service
}

// New creates a new configured middleware.
func New(config Config) (*Middleware, error) {
	newMiddleware := &Middleware{}

	return newMiddleware, nil
}

// Middleware is middleware collection.
type Middleware struct {
}
