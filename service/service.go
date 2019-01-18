// Package service implements business logic of the micro service.
package service

import (
	kitlog "github.com/go-kit/kit/log"

	"github.com/giantswarm/microendpoint/service/version"
	"github.com/giantswarm/microerror"

	"github.com/giantswarm/auto-oncall/service"
	"github.com/giantswarm/auto-oncall/service/webhook"
)

// Config represents the configuration used to create a new service.
type Config struct {
	Logger micrologger.Logger

	Description string
	Flag        *flag.Flag
	GitCommit   string
	Name        string
	Source      string
	Viper       *viper.Viper
}

// Service bundles other services.
type Service struct {
	// Dependencies
	Version *version.Service
	Webhook *webhook.Service

	// Settings
	Flag  *flag.Flag
	Viper *viper.Viper
}

// New creates a new configured service object.
func New(config Config) (*Service, error) {
	var err error
	var versionService *version.Service
	{
		versionConfig := version.Config{
			Description: config.Description,
			GitCommit:   config.GitCommit,
			Name:        config.Name,
			Source:      config.Source,
		}

		versionService, err = version.New(versionConfig)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var webhookService *webhook.Service
	{
		webhookConfig := webhook.Config{
			Logger: config.Logger,

			ConfigFilePath: config.Viper.GetString(config.Flag.Service.Oncall.ConfigFilePath),
			OpsgenieToken:  config.Viper.GetString(config.Flag.Service.Oncall.OpsgenieToken),
			WebhookSecret:  config.Viper.GetString(config.Flag.Service.Oncall.WebhookSecret),
		}

		webhookService, err = webhook.New(webhookConfig)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	newService := &Service{
		Version: versionService,
		Webhook: webhookService,
	}

	return newService, nil
}
