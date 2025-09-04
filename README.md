# Gokord

Gokord is a hard fork of [DiscordGo](https://github.com/bwmarrin/discordgo) because:
- maintainers are inactives
- the code base does not follow Go recommendations (and it has many spelling mistakes in its docs)
- [maintainers do not want to upgrade Go and libraries used to a newer version fixing 4 CVE](https://github.com/bwmarrin/discordgo/pull/1528)

Check the [ROADMAP](/ROADMAP.md) for more information.

[![Go Reference](https://pkg.go.dev/badge/github.com/nyttikord/gokord.svg)](https://pkg.go.dev/github.com/nyttikord/gokord) [![CI](https://github.com/nyttikord/gokord/actions/workflows/ci.yml/badge.svg)](https://github.com/nyttikord/gokord/actions/workflows/ci.yml)

Gokord is a [Go](https://golang.org/) package that provides low level 
bindings to the [Discord](https://discord.com/) chat client API. Gokord 
has nearly complete support for all the Discord API endpoints, websocket
interface, and voice interface.

<!--
If you would like to help the Gokord package please use 
[this link](https://discord.com/oauth2/authorize?client_id=173113690092994561&scope=bot)
to add the official Gokord test bot **dgo** to your server. This provides 
indispensable help to this project.
-->

## Getting Started

### Installing

This assumes you already have a working Go environment, if not, please see
[this page](https://golang.org/doc/install) first.

`go get` *will always pull the latest tagged release from the master branch.*

```sh
go get -u github.com/nyttikord/gokord
```

### Usage

Import the package into your project.

```go
import "github.com/nyttikord/gokord"
```

Construct a new Discord client which can be used to access the variety of 
Discord API functions and to set callback functions for Discord events.

```go
discord, err := gokord.New("Bot " + "authentication token")
```

See Documentation and Examples below for more detailed information.


## Documentation

**NOTICE**: This library and the Discord API are unfinished.
Because of that there may be major changes to the library in the future.

The Gokord code is fairly well documented at this point and is currently
the only documentation available. Go reference (below) presents that information in a nice format.

- [![Go Reference](https://pkg.go.dev/badge/github.com/nyttikord/gokord.svg)](https://pkg.go.dev/github.com/nyttikord/gokord)


## Examples

Below is a list of examples and other projects using Gokord. 

- [Gokord examples](https://github.com/nyttikord/gokord/tree/main/examples) — A collection of example programs written with Gokord (really outdated)

<!--
## Troubleshooting
For help with common problems please reference the 
[Troubleshooting](https://github.com/bwmarrin/discordgo/wiki/Troubleshooting) 
section of the project wiki.
-->


## Contributing
Contributions are very welcomed, however please follow the below guidelines.

- First open an issue describing the bug or enhancement so it can be
discussed.  
- Try to match current naming conventions as closely as possible.  
- This package is intended to be a low level direct mapping of the Discord API, 
so please avoid adding enhancements outside of that scope without first 
discussing it.
- Create a Pull Request with your changes against the main branch.


## List of Discord APIs

See [this chart](https://abal.moe/Discord/Libraries.html) for a feature 
comparison and list of other Discord API libraries.

## Special Thanks

[Chris Rhodes](https://github.com/iopred) - For the DiscordGo logo and tons of PRs.
