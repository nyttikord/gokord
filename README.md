# gokord

[![Go Reference](https://pkg.go.dev/badge/github.com/nyttikord/gokord.svg)](https://pkg.go.dev/github.com/nyttikord/gokord)
[![CI](https://github.com/nyttikord/gokord/actions/workflows/ci.yml/badge.svg)](https://github.com/nyttikord/gokord/actions/workflows/ci.yml)

gokord is a [Go](https://go.dev/) package that provides low level bindings to the [Discord](https://discord.com/)
chat client API.
gokord has nearly complete support for all the Discord API endpoints, and websocket interface.
We have decided to remove the really outdated voice API, because our maintainers do not use it.
Feel free to open a PR to add it!

<!--
If you would like to help the Gokord package please use 
[this link](https://discord.com/oauth2/authorize?client_id=173113690092994561&scope=bot)
to add the official Gokord test bot **dgo** to your server. This provides 
indispensable help to this project.
-->

gokord is a hard fork of [DiscordGo](https://github.com/bwmarrin/discordgo) because:
- maintainers are inactives;
- [maintainers do not want to upgrade to newer versions to fix 4 CVE](https://github.com/bwmarrin/discordgo/pull/1528).

Check the [ROADMAP](/ROADMAP.md) for more information.

See below for the main differences between gokord and DiscordGo.

## Getting Started

### Installing

This assumes you already have a working Go environment.

`go get` *will always pull the latest tagged release from the main branch.*

```sh
go get -u github.com/nyttikord/gokord
```

### Usage

Import the package into your project.

```go
import "github.com/nyttikord/gokord"
```

Construct a new Discord client which can be used to access the variety of  Discord API functions and to set callback
functions for Discord events.

```go
dg := gokord.New("Bot " + "authentication token")
// do some config things, like setting the intents or the sharding
err := dg.Open() // this starts the websocket and connect it to the Discord API.
```

See Documentation and Examples below for more detailed information.

## Why should you use gokord instead of DiscordGo?

We have completely refactored the code base to clean it and to update it.

### What you can see

We have:
- recoded the websocket to follow Discord documentations;
- modified how goroutines were managed to avoid data races with Mutex (this actually leads to a faster code :D);
- fixed lot of bugs related to invalid handling of rate limits;
- updated the documentation;
- added missing features.

With gokord, you can also choose where you state is stored: you can stay with maps, or you can go with in-memory
key-value database like Valkey or Redis.

You can use our powerful logger based on `log/slog`!

#### What you will see

We are currently simplifying how data is shared between goroutines using `context` package to build a more seamless
development experience.

We are going to refactor how interactions are handled and we will create helpful structs to achieve the same goal.

Check the [ROADMAP](./ROADMAP.md) to have the details ;D

### What you cannot see

Everything were stored in huge files (`restapi.go`, `websocket.go`, `structs.go`...) located in the package `discordgo`.
We have decided to split these to increase the readability and the ease of maintenance of our code.

We also have:
- changed the websocket API to use a well-maintained one;
- refactored how data are handled internally to avoid data races and to be faster;
- modified how the code is written to fit with what we think that are the best practices.

## Migrating to gokord from DiscordGo

Due to its nature as a hard fork, migrating to gokord is equivalent to rewriting your application.
Take a look at our examples to see the major differences.
You can also check the documentation which is nearly complete.

## Documentation

**NOTICE**: This library and the Discord API are unfinished.
We are following the [Semantic Versioning](https://semver.org/) and we plan to release the `1.0.0` after cleaning the
library.
Next major breaking changes will only be introduced when the major version is increased.

Currently, we are using Discord API v10.
When we will switch to Discord API v11 with its breaking changes, we will increase the major version. 
We can also increase it if we introduce breaking changes in our library.
We *will not* increase it if Discord introduces breaking changes in the current supported API version.

The gokord code is fairly well documented at this point and this is currently the only documentation available.
Go reference (below) presents that information in a nice format.

- [![Go Reference](https://pkg.go.dev/badge/github.com/nyttikord/gokord.svg)](https://pkg.go.dev/github.com/nyttikord/gokord)

## Examples

Below is a list of examples and other projects using Gokord. 

- [gokord examples](https://github.com/nyttikord/gokord/tree/main/examples) — A collection of example programs written with gokord (really outdated)

If you want real world example, you can check [our bots](https://github.com/nyttikord) or 
[Les Copaings Bot](https://git.anhgelus.world/anhgelus/les-copaings-bot).

<!--
## Troubleshooting
For help with common problems please reference the 
[Troubleshooting](https://github.com/bwmarrin/discordgo/wiki/Troubleshooting) 
section of the project wiki.
-->

## Contributing
Contributions are very welcomed, however please follow the below guidelines.

- First open an issue describing the bug or enhancement so it can be discussed.  
- Try to match current naming conventions as closely as possible.  
- This package is intended to be a low level direct mapping of the Discord API, so please avoid adding enhancements
outside of that scope without first discussing it.
- Create a Pull Request with your changes against the main branch.

Check [CONTRIBUTING.md](/CONTRIBUTING.md) for more information.

