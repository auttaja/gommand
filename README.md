# gommand

A fast command router and processor for Go.

disgord is great on its own if you want a low level Discord API wrapper, but sometimes you will want some higher level functionality. This is where gommand comes in.

gommand is a command router at heart, but is far more capable than that. It features the following benefits over just manually handling the disgord message event:

- **High performance** - The gommand codebase is designed to be as high performance as possible, attempting to reduce memory alocation.
- **Very usable** - gommand features a lot of functionality when it is parsing commands such as a flexible transformation API, greedy/optional commands and the member object is added back into messages.

An example is accessable in the `example` folder.
