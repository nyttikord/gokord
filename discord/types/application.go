// Package types contains every type constant documented in the Discord API.
//
// For example, Channel represents the type of channel.Channel and ChannelGuildText is the type of classical channel
// guild text.
package types

// Integration specifies where the application can be installed and its available interaction contexts.
type Integration uint

const (
	// IntegrationGuildInstall indicates that the app is installable to guilds.
	IntegrationGuildInstall Integration = 0
	// IntegrationUserInstall indicates that the app is installable to users.
	IntegrationUserInstall Integration = 1
)

// RoleConnectionMetadata represents the type of application role connection metadata.
type RoleConnectionMetadata int

// Application role connection metadata types.
const (
	RoleConnectionMetadataIntegerLessThanOrEqual     RoleConnectionMetadata = 1
	RoleConnectionMetadataIntegerGreaterThanOrEqual  RoleConnectionMetadata = 2
	RoleConnectionMetadataIntegerEqual               RoleConnectionMetadata = 3
	RoleConnectionMetadataIntegerNotEqual            RoleConnectionMetadata = 4
	RoleConnectionMetadataDatetimeLessThanOrEqual    RoleConnectionMetadata = 5
	RoleConnectionMetadataDatetimeGreaterThanOrEqual RoleConnectionMetadata = 6
	RoleConnectionMetadataBooleanEqual               RoleConnectionMetadata = 7
	RoleConnectionMetadataBooleanNotEqual            RoleConnectionMetadata = 8
)
