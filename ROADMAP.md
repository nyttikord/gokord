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

**BREAKING CHANGES EVERYWHERE!**

Reorganizes the source code.
Currently, it's a mess.

This release is unstable and it is not recommended for production use.
Only the REST API was heavily refactored.
The Websocket API (including events and voice) was not touched.

Huge files will be split into smaller ones.
Subpackages will be created to be more maintainable and to be easier to use.

## `0.32.0`

**BREAKING CHANGES IN WEBSOCKET API AND IN STATE!**
Including events and voice.

Refactor the Websocket API (including events and voice).
Rewrite how the Session works.
Rewrite how the State is managed.

This release follows the changes of `0.31.0`.
It does not add new features, but continue the cleaning of the source code.

The goal of this is to be more stable than the `0.31.0`.
We don't know yet if this release could be used in a production environment.

## `0.33.0`

Imports some features from [`anhgelus/gokord`](https://github.com/anhgelus/gokord) to provide useful structs to manage:
- slash commands
- interaction responses
- components

Provides a new easier way to create a bot.

Rewrite interaction package to use contexts.

## `1.0.0`

Release if everything is fine?
