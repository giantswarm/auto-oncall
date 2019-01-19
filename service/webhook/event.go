package webhook

type Commit struct {
	Author Author `json:"author"`
	SHA    string `json:"sha"`
}

type Author struct {
	Login string `json:"login"`
}

type DeploymentEvent struct {
	Environment string     `json:"environment"`
	Ref         string     `json:"ref"`
	Repository  Repository `json:"repository"`
}

type Repository struct {
	FullName string `json:"full_name"`
	Name     string `json:"name"`
}
