# Getting started

To start off you can check out existing Pull Requests and Issues to get a gasp of what problems we're currently solving and what features you can implement.

## Issues

Our issues are mostly used for bugs, however we welcome refactoring and conceptual issues.

Any other conversation would belong and would be moved into “Discussions”.

## Discussions

We use discussions for ideas, polls, announcements and help questions.

Don't hesitate to ask, we always would try to help.

## Pull Requests

If you want to help us by improving existing or adding new features, you create what's called a Pull Request (aka PR). It allows us to review your code, suggest changes and merge it.

Please, [test your code](#testing-your-code) before submitting a PR!

Here are some tips on how to make a good first PR:
- When creating a PR, please consider a distinctive name and description for it, so the maintainers can understand what your PR changes / adds / removes.
- It's always a good idea to link documentation when implementing a new feature / endpoint
- If you're resolving an issue, don't forget to [link it](https://docs.github.com/en/issues/tracking-your-work-with-issues/linking-a-pull-request-to-an-issue) in the description.
- Enable the checkbox to allow maintainers to edit your PR and make commits in the PR branch when necessary.
- We may ask for changes, usually through suggestions or pull request comments. You can apply suggestions right in the UI. Any other change needs to be done manually.
- Don't forget to mark PR comments resolved when you're done applying the changes.
- Be patient and don't close and reopen your PR when no one responds, sometimes it might be held for a while. There might be a lot of reasons: release preparation, the feature is not significant, maintainers are busy, etc.

When your changes are still incomplete (i.e. in Work In Progress state), you can still create a PR, but consider making it a draft. 
To make a draft PR, you can change the type of PR by clicking to a triangle next to the “Create Pull Request” button.

Once you're done, you can mark it as “Ready for review”, and we'll get right on it.

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

# Code style

To standardize and make things less messy we have a certain code style, that is persistent throughout the codebase.

## Organization

Structures and functions under the same endpoint (e.g., `/guild` or `/channel`) are in the same package named with the
endpoint's name.
For example, `Role` is in the `guild` package because we must call `/guild/roles` to get a role.

If the package has specific REST method, it has a subpackage called `endpointapi` (e.g., `guildapi` or `channelapi`)
containing every these REST method.
In this package, the `Requester` struct implements `discord.Requester` and it is used to send the requests to the Discord
API.
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

If an endpoint/set of endpoints have mostly same parameters, it's a good idea to use a single `Param` structure for them.
Here's an example: 
> Endpoint: `GuildMemberEdit`
>
> `Param` structure: `GuildMemberParams` 

If an endpoint/set of endpoints have differentiating parameters, `Param` structure can be named after the endpoint's verb.
Here's an example:
> Endpoint: `ChannelMessageSendComplex`
>
> `Param` structure: `MessageSend`

> Endpoint: `ChannelMessageEditComplex`
>
> `Param` structure: `MessageEdit` 

### Events

When naming an event, we follow gateway's internal naming (which often matches with the official event name in the docs).
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

