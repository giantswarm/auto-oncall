package oncall

type Oncall struct {
	GithubToek    string `yaml:"githubToken"`
	OpsgenieToken string `yaml:"opsgenieToken"`
	Users         string `yaml:"users"`
	WebhookSecret string `yaml:"webhookSecret"`
}
