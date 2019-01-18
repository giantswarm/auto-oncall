// Package server provides a server implementation to connect network transport
// protocols and service business logic by defining server endpoints.
//
package server

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"

	"github.com/giantswarm/microerror"
	microserver "github.com/giantswarm/microkit/server"
	microvalidator "github.com/giantswarm/microkit/validator"
	"github.com/giantswarm/micrologger"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/spf13/viper"

	"github.com/giantswarm/auto-oncall/flag"
	"github.com/giantswarm/auto-oncall/server/endpoint/webhook"
	"github.com/giantswarm/auto-oncall/service"
)

// Config represents the configuration used to create a new server object.
type Config struct {
	Flag    *flag.Flag
	Logger  micrologger.Logger
	Service *service.Service
	Viper   *viper.Viper

	ProjectName string
}

type Server struct {
	// Dependencies.
	logger micrologger.Logger

	// Internals.
	bootOnce     sync.Once
	config       microserver.Config
	shutdownOnce sync.Once
}

// New creates a new configured server object.
func New(config Config) (*Server, error) {
	var err error

	var endpointCollection *endpoint.Endpoint
	{
		endpointConfig := endpoint.Config{
			Flag:    config.Flag,
			Logger:  config.Logger,
			Service: config.Service,
			Viper:   config.Viper,
		}
		endpointCollection, err = endpoint.New(endpointConfig)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	s := &Server{
		// Dependencies.
		logger: config.Logger,

		// Internals.
		bootOnce: sync.Once{},
		config: microserver.Config{
			Logger:      config.Logger,
			ServiceName: config.ProjectName,
			Viper:       config.Viper,

			Endpoints: []microserver.Endpoint{
				endpointCollection.Version,
			},
		},
		shutdownOnce: sync.Once{},
	}

	return s, nil
}

func (s *Server) Boot() {
	s.bootOnce.Do(func() {
		// Here goes your custom boot logic for your server/endpoint/middleware, if
		// any.
	})
}

func (s *Server) Config() microserver.Config {
	return s.config
}

func (s *Server) Shutdown() {
	s.shutdownOnce.Do(func() {
		// Here goes your custom shutdown logic for your server/endpoint/middleware,
		// if any.
	})
}
