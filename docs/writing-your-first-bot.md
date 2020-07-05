# Writing your first bot

Creating your first bot with Gommand is a breeze:

## Creating the router
Firstly, we will create the new [router](./router.md) for all disgord/gommand events. In this example, we are using several of the built-in [prefix checkers](./prefix-checkers.md) (`MultiplePrefixCheckers` to use multiple prefix checkers, `StaticPrefix("%")` to use the `%` prefix, and `MentionPrefix` to allow for the mentioning of the bot):
```go
var router = gommand.NewRouter(&gommand.RouterConfig{
	// The prefix function should be set here or it will be blank.
	// We are using % and mention prefixes for this example.
	PrefixCheck: gommand.MultiplePrefixCheckers(gommand.StaticPrefix("%"), gommand.MentionPrefix),
})
```

The router will also create a basic help command which you can either use or delete.

## Setting the command
To set a command, we will want to call the `SetCommand` function on the router. To set the command, we will create a new instance of the [`Command`](./commands.md#Command) struct and then set it to the router:
```go
func init() {
	router.SetCommand(&gommand.Command{
		Name:        "ping",
		Description: "Responds with pong.",
		Function: func(ctx *gommand.Context) error {
			_, _ = ctx.Reply("Pong!")
			return nil
		},
	})
}
```

## Starting the bot
Firstly, we will want to setup the fatal error handler for disgord's event logger. This is where errors will be passed through to in the event that none of the error handlers are able to handle it.

In this example, we will use [logrus](https://github.com/sirupsen/logrus) since it is plug and play with disgord, but you are free to use your own handler here. In the beginning of the main function, we will configure this:
```go
func main() {
	logrus.SetLevel(logrus.DebugLevel)
    ...
}
```
From here, we will create the disgord session with the logger and grabbing the `TOKEN` environment variable:
```go
	s := disgord.New(disgord.Config{
		BotToken: os.Getenv("TOKEN"),
		Logger:   logrus.New(),
	})
```
We will then need to hook the router to the client:
```go
    router.Hook(s)
```
Any further events being added to the session can be added here. From here, we can simply add the code to start the bot:
```go
    err := s.StayConnectedUntilInterrupted(context.Background())
    if err != nil {
        panic(err)
    }
```

Congratulations, you have created a bot which responds with pong to the ping command:

![cmd](https://i.imgur.com/NXKXQFT.png)

Congratulations, you wrote your first bot! Now it is time to extend this bot with some new functionality.

## Creating the error handler
Right now, the console will be spammed with errors for situations such as when a command isn't found. This is because Gommand doesn't have an error handler, it will automatically pass errors through to the disgord logger (which right now just logs to console). We can change this by adding a error handler in the initialisation of this bot with the `AddErrorHandler` function:

```go
    router.AddErrorHandler(func(ctx *gommand.Context, err error) bool {
		// Check all the different types of errors.
		switch err.(type) {
		case *gommand.CommandNotFound, *gommand.CommandBlank:
			// We will ignore. The command was not found in the router or the command was blank.
			return true
		case *gommand.InvalidTransformation:
			_, _ = ctx.Reply("Invalid argument:", err.Error())
			return true
		case *gommand.IncorrectPermissions:
			_, _ = ctx.Reply("Invalid permissions:", err.Error())
			return true
		case *gommand.InvalidArgCount:
			_, _ = ctx.Reply("Invalid argument count.")
			return true
		}

		// This was not handled here.
		return false
	})
```

Gommand contains the following built-in errors:

- `*gommand.CommandNotFound`: The command was not found within the router.
- `*gommand.CommandBlank`: The command name was blank.
- `*gommand.IncorrectPermissions`: The permissions this user has are incorrect for the command.
- `*gommand.InvalidArgCount`: The argument count is not correct.
- `*gommand.InvalidTransformation`: Passed through from a transformer when it cannot transform properly.
- `*gommand.PanicError`: This is used when a string is returned from a panic. If it isn't a string, the error will just be pushed into the handler.

The boolean in this function represents whether the error should be parsed through to the next error handler. If true is returned, it is handled within the function. If false is returned, it will be passed through to the next error handler, or to disgord's default logger if there are no error handlers after it.

## Waiting for a message
The [`Context`](./context.md) struct which is provided when a command is ran contains some extremely powerful features. One example of this would be waiting for a message. To wait for a message, we can use the `WaitForMessage` function to wait for a message based on the condition specified. When the condition is met, the response will be the message:
```go
    // Echos the message. You can change the context to have a timeout.
    resp := ctx.WaitForMessage(context.TODO(), func(_ disgord.Session, msg *disgord.Message) bool {
        return msg.Author.ID == ctx.Message.Author.ID && msg.ChannelID == ctx.Message.ChannelID
    })
    _, _ = ctx.Reply(resp.Content)
```

## Aliases
But what if you want the ping command to respond to `%p`? The [`Command`](./commands.md#Command) struct contains a `Aliases` field among many other things which you can use to set a simple string array to handle exactly this:
```go
    Aliases: []string{"p"},
```

## What to check out next
Now you have learned the basic structure of Gommand, you may want to check out the following:

- [Categories](./categories.md)
- [Handling deleted messages](./handling-deleted-messages.md)
- [Permission validators](./permission-validators.md)
- [Middleware](./middleware.md)
- [Embed paginator](./embed-paginator.md)
- [Embed menus](./embed-menus.md)
