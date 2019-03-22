package oncall

type Oncall struct {
	GithubToken   string `yaml:"githubToken"`
	OpsgenieToken string `yaml:"opsgenieToken"`
	Users         string `yaml:"users"`
	WebhookSecret string `yaml:"webhookSecret"`
}
