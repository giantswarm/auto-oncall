package webhook

type Commit struct {
	Author Author `json:"author"`
}

type Author struct {
	Login string `json:"login"`
}

type DeploymentEvent struct {
	Deployment Deployment
	Repository Repository `json:"repository"`
}

type Deployment struct {
	Creator     Creator
	Environment string `json:"environment"`
	Ref         string `json:"ref"`
}

type Creator struct {
	Login string `json:"login"`
}

type Repository struct {
	FullName string `json:"full_name"`
	Name     string `json:"name"`
}
