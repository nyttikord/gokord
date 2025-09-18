// Package user contains everything related to Discord User, guild.Guild Member and Integration.
//
// Use userapi.Requester to interact with this.
// You can get this with gokord.Session.
package user

import (
	"time"
)

// ExpireBehavior of an Integration.
// https://discord.com/developers/docs/resources/guild#integration-object-integration-expire-behaviors
type ExpireBehavior int

const (
	ExpireBehaviorRemoveRole ExpireBehavior = 0
	ExpireBehaviorKick       ExpireBehavior = 1
)

// IntegrationAccount is integration account information sent while fetching the Connection.
type IntegrationAccount struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Integration stores integration information.
type Integration struct {
	ID                string             `json:"id"`
	Name              string             `json:"name"`
	Type              string             `json:"type"`
	Enabled           bool               `json:"enabled"`
	Syncing           bool               `json:"syncing"`
	RoleID            string             `json:"role_id"`
	EnableEmoticons   bool               `json:"enable_emoticons"`
	ExpireBehavior    ExpireBehavior     `json:"expire_behavior"`
	ExpireGracePeriod int                `json:"expire_grace_period"`
	User              *User              `json:"user"`
	Account           IntegrationAccount `json:"account"`
	SyncedAt          time.Time          `json:"synced_at"`
}

// Connection is a connection returned from the UserConnections endpoint.
type Connection struct {
	ID           string         `json:"id"`
	Name         string         `json:"name"`
	Type         string         `json:"type"`
	Revoked      bool           `json:"revoked"`
	Integrations []*Integration `json:"integrations"`
}
