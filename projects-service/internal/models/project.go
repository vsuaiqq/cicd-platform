package models

import "time"

type Project struct {
	ID                  string
	UserID              string
	Name                string
	RepoURL             string
	DefaultBranch       string
	PrivateKeyEncrypted []byte
	PublicKey           string
	WebhookSecret       string
	Status              string
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

type EnvVar struct {
	Key   string
	Value string
}

type SecretMeta struct {
	Key       string
	UpdatedAt time.Time
}

type Secret struct {
	Key   string
	Value string
}

type ProjectMember struct {
	ID          string
	ProjectID   string
	UserID      string
	Email       string
	DisplayName string
	Role        string
	InvitedBy   string
	CreatedAt   time.Time
}
