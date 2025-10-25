package types

// Component is type of component.Component.
type Component uint

const (
	ComponentActionsRow            Component = 1
	ComponentButton                Component = 2
	ComponentSelectMenu            Component = 3
	ComponentTextInput             Component = 4
	ComponentUserSelectMenu        Component = 5
	ComponentRoleSelectMenu        Component = 6
	ComponentMentionableSelectMenu Component = 7
	ComponentChannelSelectMenu     Component = 8
	ComponentSection               Component = 9
	ComponentTextDisplay           Component = 10
	ComponentThumbnail             Component = 11
	ComponentMediaGallery          Component = 12
	ComponentFile                  Component = 13
	ComponentSeparator             Component = 14
	ComponentContainer             Component = 17
	ComponentLabel                 Component = 18
	ComponentFileUpload            Component = 19
)

// SelectMenuDefaultValue represents the type of entity selected by default in autopopulated component.SelectMenuDefaultValue.
type SelectMenuDefaultValue string

const (
	SelectMenuDefaultValueUser    SelectMenuDefaultValue = "user"
	SelectMenuDefaultValueRole    SelectMenuDefaultValue = "role"
	SelectMenuDefaultValueChannel SelectMenuDefaultValue = "channel"
)

// SelectMenu represents the type of component.SelectMenu.
type SelectMenu Component

const (
	SelectMenuString      = SelectMenu(ComponentSelectMenu)
	SelectMenuUser        = SelectMenu(ComponentUserSelectMenu)
	SelectMenuRole        = SelectMenu(ComponentRoleSelectMenu)
	SelectMenuMentionable = SelectMenu(ComponentMentionableSelectMenu)
	SelectMenuChannel     = SelectMenu(ComponentChannelSelectMenu)
)
