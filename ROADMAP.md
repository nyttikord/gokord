# ROADMAP

Goals of gokord:
- updating DiscordGo
- stay up to date with Discord API
- provides a cleaner library to interact with Discord API

## `0.33.0`

**BREAKING CHANGES EVERYWHERE!**
We have decided to replace `gorilla/websocket` by `coder/websocket` to support contexts and to have a well maintained
library.
Now, you must use contexts to open and close the bot.
Events gives the current context in the handler, so you will have to modify every handlers to follow the new signature. 

Rewrite state to be able to work with a custom implementation of storage.

The library is more stable, thanks to contexts.
This release can be used in production, yay :D
(and it seems to be more stable than DiscordGo.)

## `0.34.0`

**BREAKING CHANGES EVERYWHERE**
Introduce contexts everywhere.

New higher-level interaction package using contexts:
- easier to declare
- easier to handle
- easier to respond (like the one in [`anhgelus/gokord`](https://github.com/anhgelus/gokord))
- everything is managed via contexts to catch timeout errors and to provide a cleaner syntax

Rewrite HTTP API to create dynamic requests and to leave OOP.

Create a higher-level component creator to limit errors and to provide a cleaner syntax.

Use AVL as default state storage to reduce memory usage.

## `0.35.0`

New higher-level slash commands package using contexts (like the one in 
[`anhgelus/gokord`](https://github.com/anhgelus/gokord)):
- easier to declare
- easier to handle
- automatic deploy

Provides a new easier way to create a bot.

## `1.0.0`

Release if everything is fine?
