# Getting started

Before opening a pull request (PR) or an issue, check the existing ones to avoid creating a duplicate.

If you find a bug or if you have a suggestion to improve the library, you may open an issue.

## Pull Requests

Fork the repo, create a branch, commits and then open a PR.
We encourage you to open an issue first (if yours does not resolve one).
When you are writing the description of your PR, don't forget to
[link it](https://docs.github.com/en/issues/tracking-your-work-with-issues/linking-a-pull-request-to-an-issue).

Please, [test your code](#testing-your-code)!

When you have finished your PR, one of our maintainers will review your work.
If everything is fine, it will be merged (yay :D).
If your work is mostly good, the reviewer will ask you to fix issues.
If your PR has conceptual issues, it will be closed and the reviewer will explain you why.

**Read our the rest of this document to avoid losing your (and our) times.**

We encourage you to watch
[this conference at FOSDEM 2026](https://fosdem.org/2026/schedule/event/L7ERNP-prs-maintainers-will-love/) to understand
how to write good PR.

# Use of AI

The maintainers of gokord do not use LLMs.
We are against this technology for many reasons, but we can't stop you from using these tools.

**When you are interacting with us, do not use an LLM.**
If you are, we will instantly close your issue/your PR.

**If you have used an LLM to help you, you must inform us.**
We will not reject your work for this reason.
If you hide this, we will instantly close your issue/your PR.

We will not accept poorly written code with useless comments.
We will not accept code that does not follow our code style.
We will not accept PR that does not follow our contributing guide.
You have to modify the LLM's output.
You have to test the code that you want to be merged.
If you don't understand something created by an LLM, avoid creating a PR/an issue, thanks you.

**Remember that gokord is made for humans and by humans.**
If your contribution does not follow this vision, avoid contributing.
Every PR is reviewed by at least one human, every issue is triaged by at least one human, every line of code was written
for humans.

# Style

To standardize and make things less messy, we have a certain code style that is persistent throughout the codebase.

## Commits

A commit is an atomic modification.
It cannot be divided into smaller ones.
It must update tests and it must work without additional commits.

We follow this simple schema for their name:
```
kind(scope): description
```
`kind` is the kind of modification:
- `feat` for an addition
- `refactor` for a refactor
- `fix` for fixing an issue
- `style` for the code style
- `docs` for the documentations
- `build` for building tools
- `ci` for CI/CD

`scope` indicates the part touched of your modification.
We commonly use `ws` for websocket, `guild` for `guild` and `guild/guildapi` packages...

`description` is a *short* description of your modification.
If you want to explain more things, include them in other lines (not the first one).

If your history is messy, you must modify it and force push the updated version.
In futur versions of git, you will be able to use the new `git-history(1)` to edit this easily.

```
fix(ws): not panicing if bad setup during connection
```
is fixing an issue related to websocket.
Before this commit, the bot was not panicing if there is a bad setup during the connection.
After this commit, this issue is fixed.

```
feat(logger): option to trim version in caller
```
is adding something to the logger.
Now, the developer can use an option to trim versions.

## Organization

Structures and functions under the same endpoint (e.g., `/guild` or `/channel`) are in the same package named with the
endpoint's name.
For example, `Role` is in the `guild` package because we must call `/guild/roles` to get a role.

If the package has specific REST method, it has a subpackage called `endpointapi` (e.g., `guildapi` or `channelapi`)
containing every these REST method.
In this package, the `Requester` struct implements `discord.Requester` and it is used to send the requests to the
Discord API.
This `Requester` is obtainable with the method `gokord.Session.EndpointAPI()` (e.g., `GuildAPI()` or `ChannelAPI()`).
```go
var s *gokord.Session // this is a valid gokord session
var g *guild.Guild
var err error
g, err = s.GuildAPI().Guild("0123456789") // this request the guild with the ID "0123456789"
if err != nil {
	// an error occurred
}
```

Constant used by Discord are located in the package `discord`. Its subpackage `types` contains types used in multiple
places.

## Naming

### REST methods

When naming a REST method, while it might seem counterintuitive, we specify the entity before the action verb (for GET
endpoints we don't specify one however).
Here's an example:

> Endpoint name: Get Channel Message
>
> Method name: `ChannelMessage`

> Endpoint name: Edit Channel Message
>
> Method name: `ChannelMessageEdit`

### Parameter structures

When making a complex REST endpoint, sometimes you might need to implement a `Param` structure.
This structure contains parameters for certain endpoint/set of endpoints.

If an endpoint/set of endpoints have mostly same parameters, it's a good idea to use a single `Param` structure for
them.
Here's an example: 
> Endpoint: `GuildMemberEdit`
>
> `Param` structure: `GuildMemberParams` 

If an endpoint/set of endpoints have differentiating parameters, `Param` structure can be named after the endpoint's
verb.
Here's an example:
> Endpoint: `ChannelMessageSendComplex`
>
> `Param` structure: `MessageSend`

> Endpoint: `ChannelMessageEditComplex`
>
> `Param` structure: `MessageEdit` 

### Events

When naming an event, we follow gateway's internal naming (which often matches with the official event name in the
docs).
Here's an example:
> Event name: Interaction Create (`INTERACTION_CREATE`)
>
> Structure name: `InteractionCreate`

# Testing your code

Before submitting a PR, you must test your changes.

First, you can simply test the websocket with
```bash
go run ./tools/cmd/session/main.go -token YOUR_TOKEN
```
You can also set the token with the environment variable `DG_TOKEN`.

If everything looks fine, we encourage you to create a simple bot in another directory.
```bash
go work init . # init a new go.work file for this module
go work use path/to/gokord # override the gokord in the go.mod by the one present in this folder
```
