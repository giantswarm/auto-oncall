package oncall

type Oncall struct {
	Repositories  string `yaml:"repositories"`
	OpsgenieToken string `yaml:"opsgenieToken"`
	Users         string `yaml:"users"`
	WebhookSecret string `yaml:"webhookSecret"`
}
