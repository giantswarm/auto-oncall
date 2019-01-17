package endpoint

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/giantswarm/opsctl/service/opsgenie"

	"github.com/giantswarm/auto-oncall/service/githubhook"
)

const (
	masterRef       = "refs/heads/master"
	listenAddres    = 8000
	webhookEndpoint = "/webhook"
	routingRuleTTL  = time.Hour * time.Duration(1)
	routingRuleType = "match-all-conditions"
)

type healthCheckResponse struct {
	Status string `json:"status"`
}

type errorResponse struct {
	Error string `json:"error"`
}

type Config struct {
	Logger        micrologger.Logger
	OpsgenieToken string
	Repositories  []string          `yaml:"repositories"`
	Users         map[string]string `yaml:"users"`
	WebhookSecret string
}

type Oncall struct {
	logger        micrologger.Logger
	opsgenieToken string
	repositories  []string
	users         map[string]string
	webhookSecret []byte
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func writeJSONResponse(w http.ResponseWriter, status int, data []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", strconv.Itoa(len(data)))
	w.WriteHeader(status)
	w.Write(data)
}

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	var response []byte

	response, _ = json.Marshal(healthCheckResponse{Status: "webhook handler"})
	writeJSONResponse(w, http.StatusOK, response)
}

func (o *Oncall) webhookHandler(w http.ResponseWriter, r *http.Request) {

	var response []byte

	h, err := githubhook.Parse(o.webhookSecret, r)
	if err != nil {
		o.logger.Log("level", "error", "message", err.Error())
		response, _ = json.Marshal(errorResponse{Error: "invalid request"})
		writeJSONResponse(w, http.StatusBadRequest, response)
	} else {
		response, _ = json.Marshal(healthCheckResponse{Status: "webhook request received"})
		writeJSONResponse(w, http.StatusOK, response)

		o.logger.Log("level", "debug", "message", fmt.Sprintf("received push event into repository %#q", h.Event.Repository.Name), "user", h.Event.Pusher.Name, "ref", h.Event.Ref)

		if h.Event.Ref == masterRef && stringInSlice(h.Event.Repository.Name, o.repositories) {
			o.logger.Log("level", "debug", "repository", h.Event.Repository.Name, "message", "push event into master branch received", "user", h.Event.Pusher.Name)

			err := o.createRoutingRule(h.Event)
			if err != nil {
				o.logger.Log("level", "error", "message", err.Error())
			}
		}
	}

}

func (o *Oncall) createRoutingRule(event githubhook.Event) error {
	var err error

	var opsGenieService *opsgenie.OpsGenie
	{
		serviceConfig := opsgenie.Config{
			Logger:    o.logger,
			AuthToken: o.opsgenieToken,
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

	user, ok := o.users[event.Pusher.Name]
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
	o.logger.Log("level", "debug", "message", fmt.Sprintf("escalation %#q for user %#q has been created", routingRule.Name, routingRule.User))

	err = opsGenieService.CreateRoutingRule(routingRule)
	if err != nil {
		return microerror.Mask(err)
	}
	o.logger.Log("level", "debug", "message", fmt.Sprintf("routing rule %#q for user %#q has been created", routingRule.Name, routingRule.User))

	return nil
}

func New(c Config) (Oncall, error) {
	if c.OpsgenieToken == "" {
		return Oncall{}, microerror.Maskf(invalidConfigError, "OPSGENIE_TOKEN environment variable token must not be empty")
	}
	if c.WebhookSecret == "" {
		return Oncall{}, microerror.Maskf(invalidConfigError, "GITHUB_WEBHOOK_SECRET environment variable must not be empty")
	}

	oncall := Oncall{
		logger:        c.Logger,
		opsgenieToken: c.OpsgenieToken,
		repositories:  c.Repositories,
		users:         c.Users,
		webhookSecret: []byte(c.WebhookSecret),
	}

	return oncall, nil
}

func (o *Oncall) NewServer() *http.ServeMux {
	s := http.NewServeMux()
	s.HandleFunc("/", defaultHandler)
	s.HandleFunc(webhookEndpoint, o.webhookHandler)

	return s
}
