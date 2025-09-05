package types

// Integration dictates where application can be installed and its available interaction contexts.
type Integration uint

const (
	// IntegrationGuildInstall indicates that app is installable to guilds.
	IntegrationGuildInstall Integration = 0
	// IntegrationUserInstall indicates that app is installable to users.
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
