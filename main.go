package main

import (
	"fmt"

	"github.com/giantswarm/microkit/command"
	microserver "github.com/giantswarm/microkit/server"
	"github.com/giantswarm/micrologger"
	"github.com/spf13/viper"

	"github.com/giantswarm/auto-oncall/flag"
	"github.com/giantswarm/auto-oncall/server"
	"github.com/giantswarm/auto-oncall/service"
)

var (
	description string     = "This is the webhook handler application for creating opsgenie routing rules from github deployment events."
	f           *flag.Flag = flag.New()
	gitCommit   string     = "n/a"
	name        string     = "auto-oncall"
	source      string     = "https://github.com/giantswarm/auto-oncall"
)

func main() {
	var err error

	// Create a new logger which is used by all packages.
	var newLogger micrologger.Logger
	{
		c := micrologger.Config{}

		newLogger, err = micrologger.New(c)
		if err != nil {
			panic(fmt.Sprintf("%#v", err))
		}
	}

	// We define a server factory to create the custom server once all command
	// line flags are parsed and all microservice configuration is storted out.
	newServerFactory := func(v *viper.Viper) microserver.Server {
		var newService *service.Service
		{
			c := service.Config{
				Logger: newLogger,

				Description: description,
				Flag:        f,
				GitCommit:   gitCommit,
				Name:        name,
				Source:      source,
				Viper:       v,
			}

			newService, err = service.New(c)
			if err != nil {
				panic(fmt.Sprintf("%#v", err))
			}
		}

		var newServer microserver.Server
		{
			c := server.Config{
				Flag:    f,
				Logger:  newLogger,
				Service: newService,
				Viper:   v,

				ProjectName: name,
			}

			newServer, err = server.New(c)
			if err != nil {
				panic(fmt.Sprintf("%#v", err))
			}
		}

		return newServer
	}

	// Create a new microkit command which manages our custom microservice.
	var newCommand command.Command
	{
		c := command.Config{
			Logger:        newLogger,
			ServerFactory: newServerFactory,

			Description: description,
			GitCommit:   gitCommit,
			Name:        name,
			Source:      source,
		}

		newCommand, err = command.New(c)
		if err != nil {
			panic(fmt.Sprintf("%#v", err))
		}
	}

	daemonCommand := newCommand.DaemonCommand().CobraCommand()

	daemonCommand.PersistentFlags().String(f.Service.Oncall.OpsgenieToken, "", "Opsgenie API token.")
	daemonCommand.PersistentFlags().StringToString(f.Service.Oncall.Users, map[string]string{}, "github_id:opsgenie_id mapppings, separated by comma.")
	daemonCommand.PersistentFlags().String(f.Service.Oncall.WebhookSecret, "", "Github organization webhook secret.")

	newCommand.CobraCommand().Execute()
}
