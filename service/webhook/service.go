package webhook

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/giantswarm/opsctl/service/opsgenie"
)

const (
	commitEndpoint        = "https://api.github.com/repos/%s/commits/%s"
	testEnvironmentPrefix = "g"
	routingRuleTTL        = time.Hour * time.Duration(1)
	routingRuleType       = "match-all-conditions"
)

type Config struct {
	Logger micrologger.Logger

	OpsgenieToken string
	Users         map[string]string
	WebhookSecret string
}

type Service struct {
	logger        micrologger.Logger
	opsgenieToken string
	users         map[string]string
	webhookSecret []byte
}

func New(c Config) (*Service, error) {
	if c.OpsgenieToken == "" {
		return nil, microerror.Maskf(invalidConfigError, "Opsgenie token must not be empty")
	}
	if c.WebhookSecret == "" {
		return nil, microerror.Maskf(invalidConfigError, "Github organization webhook secret must not be empty")
	}

	service := &Service{
		logger:        c.Logger,
		opsgenieToken: c.OpsgenieToken,
		users:         c.Users,
		webhookSecret: []byte(c.WebhookSecret),
	}

	return service, nil
}

// Process performs processing of the webhook.
func (s *Service) Process(h Hook) {
	if !strings.HasPrefix(h.DeploymentEvent.Environment, testEnvironmentPrefix) {
		s.logger.Log("level", "debug", "message", "received deployment event", "repository", h.DeploymentEvent.Repository.Name, "ref", h.DeploymentEvent.Ref, "environment", h.DeploymentEvent.Environment)

		err := s.createRoutingRule(h.DeploymentEvent)
		if err != nil {
			s.logger.Log("level", "error", "message", err.Error())
		}
	}
}

func (s *Service) createRoutingRule(event DeploymentEvent) error {
	var err error

	// get commit from refference
	resp, err := http.Get(fmt.Sprintf(commitEndpoint, event.Repository.FullName, event.Ref))
	if err != nil {
		return microerror.Mask(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	commit := Commit{}
	err = json.Unmarshal(body, &commit)
	if err != nil {
		return microerror.Mask(err)
	}

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
		opsgenie.Rule{
			Value: event.Environment,
		},
	}

	ttl := time.Now().Add(routingRuleTTL).UTC().Unix()

	user, ok := s.users[commit.Author.Login]
	if !ok {
		return microerror.Maskf(userNotFoundError, commit.Author.Login)
	}
	routingRule := &opsgenie.RoutingRule{
		Name:       fmt.Sprintf("auto-%s-%s-%s-%s", event.Repository.Name, commit.SHA[:5], commit.Author.Login, strconv.FormatInt(ttl, 10)),
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
