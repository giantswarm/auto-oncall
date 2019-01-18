package endpoint

import (
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/spf13/viper"

	"github.com/giantswarm/auto-oncall/flag"
	"github.com/giantswarm/auto-oncall/server/endpoint/version"
	"github.com/giantswarm/auto-oncall/server/endpoint/webhook"
	"github.com/giantswarm/auto-oncall/server/middleware"
	"github.com/giantswarm/auto-oncall/service"
)

// Config represents the configuration used to create a endpoint.
type Config struct {
	// Dependencies.
	Flag       *flag.Flag
	Logger     micrologger.Logger
	Middleware *middleware.Middleware
	Service    *service.Service
	Viper      *viper.Viper
}

// Endpoint is the endpoint collection.
type Endpoint struct {
	Version *version.Endpoint
	Webhook *webhook.Endpoint
}

// New creates a new configured endpoint.
func New(config Config) (*Endpoint, error) {
	if config.Flag == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Flag must not be empty", config)
	}
	if config.Viper == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Viper must not be empty", config)
	}

	var err error

	var versionEndpoint *version.Endpoint
	{
		versionConfig := version.Config{
			Logger:     config.Logger,
			Middleware: config.Middleware,
			Service:    config.Service,
		}
		versionEndpoint, err = version.New(versionConfig)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var webhookEndpoint *webhook.Endpoint
	{
		webhookConfig := webhook.Config{
			Logger:     config.Logger,
			Middleware: config.Middleware,
			Service:    config.Service,
		}
		webhookEndpoint, err = webhook.New(webhookConfig)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	newEndpoint := &Endpoint{
		Version: versionEndpoint,
		Webhook: webhookEndpoint,
	}

	return newEndpoint, nil
}
