// Package service implements business logic of the micro service.
package service

import (
	"strings"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/spf13/viper"

	"github.com/giantswarm/opsctl/service/opsgenie"

	"github.com/giantswarm/auto-oncall/flag"
	"github.com/giantswarm/auto-oncall/service/version"
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
		c := version.Config{
			Description: config.Description,
			GitCommit:   config.GitCommit,
			Name:        config.Name,
			Source:      config.Source,
		}

		versionService, err = version.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var opsgenieService *opsgenie.OpsGenie
	{
		c := opsgenie.Config{
			Logger:    config.Logger,
			AuthToken: config.Viper.GetString(config.Flag.Service.Oncall.OpsgenieToken),
		}
		opsgenieService, err = opsgenie.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var webhookService *webhook.Service
	{
		users := make(map[string]string)
		userList := config.Viper.GetString(config.Flag.Service.Oncall.Users)
		for _, user := range strings.Split(userList, ",") {
			kv := strings.Split(user, ":")
			users[kv[0]] = kv[1]
		}

		webhookConfig := webhook.Config{
			Logger: config.Logger,

			Opsgenie:      opsgenieService,
			Users:         users,
			WebhookSecret: config.Viper.GetString(config.Flag.Service.Oncall.WebhookSecret),
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
