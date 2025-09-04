# ROADMAP

Goals of Gokord:
- updating discordgo
- stay up to date with Discord API
- provides a cleaner library to interact with Discord API

## `0.30.0`

Upgrades Go version and libraries to fix *4 CVE*! 
See [bwmarrin/discordgo#1528](https://github.com/bwmarrin/discordgo/pull/1528) for more information.

Supports primary guild (tag), colors for roles, label component and buyable user perks.

Upgrades to API v10 (discordgo is still using API v9).

## `0.31.0`

Reorganizes the source code.
Currently, it's a mess.

Huge files will be split into smaller ones.
Subpackages will be created to be more maintainable and to be easier to use.

## `0.32.0`

Imports some features from [`anhgelus/gokord`](https://github.com/anhgelus/gokord) to provide useful structs to manage:
- slash commands
- interaction responses
- components

## `0.33.0`

Provides a new easier way to create a bot.

## `1.0.0`

Release if everything is fine?
