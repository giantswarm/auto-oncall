package opsgenie

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"text/template"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
)

const (
	// OpsGenieEscalationEndpoint describes an API endpoint for escalations management.
	OpsGenieEscalationEndpoint = "https://api.opsgenie.com/v2/escalations"
	// OpsGenieRoutingRulesEndpoint describes an API endpoint for specific team routing rules management.
	OpsGenieRoutingRulesEndpoint = "https://api.opsgenie.com/v2/teams/4f3b702f-7ee2-4869-bf74-c8042d3fef10/routing-rules"
	// OpsGenieAlertsEndpoint describes an API endpoint for alerts management.
	OpsGenieAlertsEndpoint = "https://api.opsgenie.com/v2/alerts"
)

// Config represents the configuration used to create a new SSH service.
type Config struct {
	// Dependencies.
	Logger micrologger.Logger

	// Internals.
	AuthToken string
}

// New creates a new configured OpsGenie service.
func New(config Config) (*OpsGenie, error) {
	// Dependencies.
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Logger must not be empty", config)
	}

	newOpsGenie := &OpsGenie{
		// Dependencies.
		logger:     config.Logger,
		httpClient: &http.Client{},

		// Internals.
		authToken: config.AuthToken,
	}

	return newOpsGenie, nil
}

type OpsGenie struct {
	// Dependencies.
	logger     micrologger.Logger
	httpClient *http.Client

	// Internals.
	authToken string
}

func (o *OpsGenie) CreateEscalation(routingRule *RoutingRule) error {
	var err error

	// construct escalation request
	var escalationReq bytes.Buffer
	{
		t := template.Must(template.New("").Parse(escalationTmpl))
		err = t.Execute(&escalationReq, routingRule)
		if err != nil {
			return microerror.Maskf(invalidTemplateError, "escalation template")
		}
	}

	req, err := newRequest("POST", OpsGenieEscalationEndpoint, bytes.NewReader(escalationReq.Bytes()), o.authToken)
	resp, err := o.httpClient.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return microerror.Maskf(executionFailedError, fmt.Sprintf("create escalation '%s' for user '%s'", routingRule.Name, routingRule.User))
	}
	// expected status codes for create request are
	// * 201 - new escalation was created
	// * 409 - escalation with such name already exists (already in desired state)
	if resp.StatusCode != 201 && resp.StatusCode != 409 {
		return microerror.Maskf(unexpectedResponseCodeError, fmt.Sprintf("expected 201 or 409, got %d", resp.StatusCode))
	}
	return nil
}

func (o *OpsGenie) CreateRoutingRule(routingRule *RoutingRule) error {
	var err error

	// construct routing rule request request
	var routingRuleReq bytes.Buffer
	{
		t := template.Must(template.New("").Funcs(fnplus).Parse(routingRuleTmpl))
		err = t.Execute(&routingRuleReq, routingRule)
		if err != nil {
			return microerror.Maskf(invalidTemplateError, "routing rule template")
		}
	}

	teamRoutingRules, err := o.GetRoutingRules()
	if err != nil {
		return microerror.Mask(err)
	}
	for _, rule := range teamRoutingRules.Rules {
		if rule.Name == routingRule.Name {
			return microerror.Maskf(routingRuleDuplicationError, fmt.Sprintf("routing rule with name '%s' already exists", routingRule.Name))
		}
	}

	req, err := newRequest("POST", OpsGenieRoutingRulesEndpoint, bytes.NewReader(routingRuleReq.Bytes()), o.authToken)
	resp, err := o.httpClient.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return microerror.Maskf(executionFailedError, fmt.Sprintf("create routing rule '%s' for user '%s'", routingRule.Name, routingRule.User))
	}
	if resp.StatusCode != 201 {
		return microerror.Maskf(unexpectedResponseCodeError, fmt.Sprintf("expected 201, got %d", resp.StatusCode))
	}

	return nil
}

func (o *OpsGenie) CreateAlert(alert *Alert) error {
	var err error

	// construct alert request
	var alertReq bytes.Buffer
	{
		t := template.Must(template.New("").Parse(alertTmpl))
		err = t.Execute(&alertReq, alert)
		if err != nil {
			return microerror.Maskf(invalidTemplateError, "alert template")
		}
	}

	req, err := newRequest("POST", OpsGenieAlertsEndpoint, bytes.NewReader(alertReq.Bytes()), o.authToken)
	resp, err := o.httpClient.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return microerror.Maskf(executionFailedError, fmt.Sprintf("alerting team '%s'", alert.Team))
	}
	if resp.StatusCode != 202 {
		return microerror.Maskf(unexpectedResponseCodeError, fmt.Sprintf("expected 202, got %d", resp.StatusCode))
	}

	return nil
}

func (o *OpsGenie) DeleteEscalation(name string) error {
	var err error

	deleteEndpoint := fmt.Sprintf("%s/%s?identifierType=name", OpsGenieEscalationEndpoint, name)
	req, err := newRequest("DELETE", deleteEndpoint, nil, o.authToken)
	resp, err := o.httpClient.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return microerror.Maskf(executionFailedError, "delete escalation")
	}
	if resp.StatusCode != 200 {
		return microerror.Maskf(unexpectedResponseCodeError, fmt.Sprintf("expected 200, got %d", resp.StatusCode))
	}

	return nil
}

func (o *OpsGenie) DeleteRoutingRule(id string) error {
	var err error

	deleteEndpoint := fmt.Sprintf("%s/%s", OpsGenieRoutingRulesEndpoint, id)
	req, err := newRequest("DELETE", deleteEndpoint, nil, o.authToken)
	resp, err := o.httpClient.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return microerror.Maskf(executionFailedError, "delete routing rule")
	}
	if resp.StatusCode != 200 {
		return microerror.Maskf(unexpectedResponseCodeError, fmt.Sprintf("expected 200, got %d", resp.StatusCode))
	}

	return nil
}

func (o *OpsGenie) GetRoutingRule(id string) (*TeamRoutingRule, error) {
	var err error

	getRuleEndpoint := fmt.Sprintf("%s/%s", OpsGenieRoutingRulesEndpoint, id)
	req, err := newRequest("GET", getRuleEndpoint, nil, o.authToken)
	resp, err := o.httpClient.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return nil, microerror.Maskf(executionFailedError, "get routing rule")
	}
	if resp.StatusCode != 200 {
		return nil, microerror.Maskf(unexpectedResponseCodeError, fmt.Sprintf("expected 200, got %d", resp.StatusCode))
	}

	var teamRoutingRuleData *TeamRoutingRuleData
	{
		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)
		err = json.Unmarshal(buf.Bytes(), &teamRoutingRuleData)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}
	return teamRoutingRuleData.Rule, nil
}

func (o *OpsGenie) GetRoutingRules() (*TeamRoutingRules, error) {
	// Fetch a list of all routing rules except default.
	req, err := newRequest("GET", OpsGenieRoutingRulesEndpoint, nil, o.authToken)
	resp, err := o.httpClient.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return nil, microerror.Maskf(executionFailedError, "list routing rules")
	}
	if resp.StatusCode != 200 {
		return nil, microerror.Maskf(unexpectedResponseCodeError, fmt.Sprintf("expected 200, got %d", resp.StatusCode))
	}

	var teamRoutingRules *TeamRoutingRules
	{
		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)
		err = json.Unmarshal(buf.Bytes(), &teamRoutingRules)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}
	return teamRoutingRules, nil
}

func newRequest(op string, endpoint string, payload io.Reader, token string) (*http.Request, error) {
	req, err := http.NewRequest(op, endpoint, payload)
	if err != nil {
		return nil, microerror.Mask(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("GenieKey %s", token))
	return req, nil
}
