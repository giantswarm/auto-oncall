package webhook

type Event struct {
	HeadCommit HeadCommit `json:"head_commit"`
	Pusher     Pusher     `json:"pusher"`
	Ref        string     `json:"ref"`
	Repository Repository `json:"repository"`
}

type HeadCommit struct {
	ID string `json:"id"`
}

type Repository struct {
	Name string `json:"name"`
}

type Pusher struct {
	Name string `json:"name"`
}
