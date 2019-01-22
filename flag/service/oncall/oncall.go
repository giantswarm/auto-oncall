package oncall

type Oncall struct {
	OpsgenieToken string `yaml:"opsgenieToken"`
	Users         string `yaml:"users"`
	WebhookSecret string `yaml:"webhookSecret"`
}
