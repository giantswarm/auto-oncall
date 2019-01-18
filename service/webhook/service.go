package webhook

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"strconv"
	"time"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/giantswarm/opsctl/service/opsgenie"
)

const (
	masterRef       = "refs/heads/master"
	routingRuleTTL  = time.Hour * time.Duration(1)
	routingRuleType = "match-all-conditions"
)

type Config struct {
	Logger micrologger.Logger

	ConfigFilePath string
	OpsgenieToken  string
	WebhookSecret  string
}

type ConfigFile struct {
	Repositories []string          `yaml:"repositories"`
	Users        map[string]string `yaml:"users"`
}

type Service struct {
	logger        micrologger.Logger
	opsgenieToken string
	repositories  []string
	users         map[string]string
	webhookSecret []byte
}

func New(c Config) (*Service, error) {
	if c.ConfigFilePath == "" {
		return nil, microerror.Maskf(invalidConfigError, "ConfigFile path must not be empty")
	}
	if c.OpsgenieToken == "" {
		return nil, microerror.Maskf(invalidConfigError, "OPSGENIE_TOKEN environment variable token must not be empty")
	}
	if c.WebhookSecret == "" {
		return nil, microerror.Maskf(invalidConfigError, "GITHUB_WEBHOOK_SECRET environment variable must not be empty")
	}

	yamlFile, err := ioutil.ReadFile(c.ConfigFilePath)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	configFile := ConfigFile{}
	err = yaml.Unmarshal(yamlFile, &configFile)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	service := &Service{
		logger:        c.Logger,
		opsgenieToken: c.OpsgenieToken,
		repositories:  configFile.Repositories,
		users:         configFile.Users,
		webhookSecret: []byte(c.WebhookSecret),
	}

	return service, nil
}

// Process performs processing of the webhook.
func (s *Service) Process(h Hook) {
	s.logger.Log("level", "debug", "message", fmt.Sprintf("received push event into repository %#q", h.Event.Repository.Name), "user", h.Event.Pusher.Name, "ref", h.Event.Ref)

	if h.Event.Ref == masterRef && stringInSlice(h.Event.Repository.Name, s.repositories) {
		s.logger.Log("level", "debug", "repository", h.Event.Repository.Name, "message", "push event into master branch received", "user", h.Event.Pusher.Name)

		err := s.createRoutingRule(h.Event)
		if err != nil {
			s.logger.Log("level", "error", "message", err.Error())
		}
	}
}

func (s *Service) createRoutingRule(event Event) error {
	var err error

	var opsGenieService *opsgenie.OpsGenie
	{
		serviceConfig := opsgenie.Config{
			Logger:    s.logger,
			AuthToken: s.opsgenieToken,
		}
		opsGenieService, err = opsgenie.New(serviceConfig)
		if err != nil {
			return microerror.Mask(err)
		}
	}

	conditions := []opsgenie.Rule{
		opsgenie.Rule{
			Value: event.Repository.Name,
		},
	}

	ttl := time.Now().Add(routingRuleTTL).UTC().Unix()

	user, ok := s.users[event.Pusher.Name]
	if !ok {
		return microerror.Maskf(userNotFoundError, event.Pusher.Name)
	}
	routingRule := &opsgenie.RoutingRule{
		Name:       fmt.Sprintf("auto-%s-%s-%s-%s", event.Repository.Name, event.HeadCommit.ID[:5], event.Pusher.Name, strconv.FormatInt(ttl, 10)),
		Conditions: conditions,
		Type:       routingRuleType,
		User:       user,
	}

	err = opsGenieService.CreateEscalation(routingRule)
	if err != nil {
		return microerror.Mask(err)
	}
	s.logger.Log("level", "debug", "message", fmt.Sprintf("escalation %#q for user %#q has been created", routingRule.Name, routingRule.User))

	err = opsGenieService.CreateRoutingRule(routingRule)
	if err != nil {
		return microerror.Mask(err)
	}
	s.logger.Log("level", "debug", "message", fmt.Sprintf("routing rule %#q for user %#q has been created", routingRule.Name, routingRule.User))

	return nil
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
