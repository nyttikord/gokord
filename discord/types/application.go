// Package types contains every type constant documented in the Discord API.
//
// For example, Channel represents the type of channel.Channel and ChannelGuildText is the type of classical guild text
// channel.
//
// Flags and specific stuff are not included in this package.
// Only global constant used in multiple files are located here.
package types

// IntegrationInstall specifies where the application can be installed and its available interaction contexts.
type IntegrationInstall uint

const (
	// IntegrationInstallGuild indicates that the app is installable to guilds.
	IntegrationInstallGuild IntegrationInstall = 0
	// IntegrationInstallUser indicates that the app is installable to users.
	IntegrationInstallUser IntegrationInstall = 1
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
