package events

type GitEvent struct {
	ProjectID  string       `json:"project_id,omitempty"`
	Platform   PlatformType `json:"platform,omitempty"`
	Event      EventType    `json:"event,omitempty"`
	Repository *Repository  `json:"repository,omitempty"`
	Ref        string       `json:"ref,omitempty"`
	Branch     string       `json:"branch,omitempty"`
	CommitSHA  string       `json:"commit_sha,omitempty"`
	Commits    []*Commit    `json:"commits,omitempty"`
	Author     *User        `json:"author,omitempty"`
	Timestamp  int64        `json:"timestamp,omitempty"`
}

type PlatformType string

const (
	GitHub PlatformType = "github"
)

type EventType string

const (
	Push EventType = "push"
)

type Repository struct {
	Name string `json:"name,omitempty"`
	URL  string `json:"url,omitempty"`
}

type Commit struct {
	SHA     string `json:"sha,omitempty"`
	Message string `json:"message,omitempty"`
	Author  *User  `json:"author,omitempty"`
	URL     string `json:"url,omitempty"`
}

type User struct {
	Login string `json:"login,omitempty"`
}
