package webhook

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

func New(c Config) (Service, error) {
	if c.ConfigFilePath == "" {
		return Oncall{}, microerror.Maskf(invalidConfigError, "ConfigFile path must not be empty")
	}
	if c.OpsgenieToken == "" {
		return Oncall{}, microerror.Maskf(invalidConfigError, "OPSGENIE_TOKEN environment variable token must not be empty")
	}
	if c.WebhookSecret == "" {
		return Oncall{}, microerror.Maskf(invalidConfigError, "GITHUB_WEBHOOK_SECRET environment variable must not be empty")
	}

	yamlFile, err := ioutil.ReadFile(c.ConfigFilePath)
	if err != nil {
		return Service{}, microerror.Mask(err)
	}

	configFile := ConfigFile{}
	err = yaml.Unmarshal(yamlFile, &configFile)
	if err != nil {
		return Service{}, microerror.Mask(err)
	}

	service := Service{
		logger:        c.Logger,
		opsgenieToken: c.OpsgenieToken,
		repositories:  configFile.Repositories,
		users:         configFile.Users,
		webhookSecret: []byte(c.WebhookSecret),
	}

	return service, nil
}

// NewHook returns a Hook from an incoming HTTP Request.
func (s *Service) NewHook(req *http.Request) (hook Hook, err error) {
	if !strings.EqualFold(req.Method, "POST") {
		return Hook{}, microerror.Maskf(executionFailedError, fmt.Sprintf("%#q requests are not supported", req.Method))
	}

	if hook.Signature = req.Header.Get("x-hub-signature"); len(hook.Signature) == 0 {
		return Hook{}, microerror.Maskf(executionFailedError, "no signature found")
	}

	if hook.EventName = req.Header.Get("x-github-event"); len(hook.EventName) == 0 {
		return Hook{}, microerror.Maskf(executionFailedError, "no event found")
	}

	if hook.ID = req.Header.Get("x-github-delivery"); len(hook.ID) == 0 {
		return Hook{}, microerror.Maskf(executionFailedError, "no event id found")
	}

	if !signedBy(hook, s.webhookSecret) {
		return Hook{}, microerror.Maskf(executionFailedError, "invalid signature found")
	}

	hook.Payload, err = ioutil.ReadAll(req.Body)
	if err != nil {
		return Hook{}, microerror.Mask(err)
	}

	err = json.Unmarshal(hook.Payload, &hook.Event)
	if err != nil {
		return Hook{}, microerror.Mask(err)
	}

	return hook, nil
}
