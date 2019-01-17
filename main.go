package main

import (
	"flag"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/giantswarm/micrologger"

	"github.com/giantswarm/auto-oncall/endpoint"
)

const (
	listenAddres           = 8080
	opsgenieTokenEnv       = "OPSGENIE_TOKEN"
	githubWebhookSecretEnv = "GITHUB_WEBHOOK_SECRET"
)

var configpath = flag.String("config", "/etc/oncall/config.yaml", "path to configuration file")
var help = flag.Bool("help", false, "show help for this tool")

func main() {
	flag.Set("logtostderr", "true")
	flag.Parse()

	if *help {
		flag.Usage()
		return
	}

	logger, err := micrologger.New(micrologger.Config{})
	if err != nil {
		panic(fmt.Sprintf("%#v", err))
	}

	var config endpoint.Config
	{
		yamlFile, err := ioutil.ReadFile(*configpath)
		if err != nil {
			panic(fmt.Sprintf("#%v", err))
		}

		err = yaml.Unmarshal(yamlFile, &config)
		if err != nil {
			panic(fmt.Sprintf("%v", err))
		}

		config.Logger = logger
		config.OpsgenieToken = os.Getenv(opsgenieTokenEnv)
		config.WebhookSecret = os.Getenv(githubWebhookSecretEnv)
	}

	var oncall endpoint.Oncall
	{
		oncall, err = endpoint.New(config)
		if err != nil {
			panic(fmt.Sprintf("%v", err))
		}
	}

	log.Printf("Running oncall server...")
	http.ListenAndServe(fmt.Sprintf(":%s", strconv.Itoa(listenAddres)), oncall.NewServer())
}
