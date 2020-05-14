package main

import (
	"context"
	"github.com/andersfylling/disgord"
	"github.com/jakemakesstuff/gommand"
	"os"
)

// Create the command router.
var router = gommand.NewRouter(&gommand.RouterConfig{
	// The prefix function should be set here or it will be blank.
	// We are using % and mention prefixes for this example.
	PrefixCheck: gommand.MultiplePrefixCheckers(gommand.StaticPrefix("%"), gommand.MentionPrefix),
})

func init() {
	// A simple command to respond with pong.
	router.SetCommand(&gommand.Command{
		Name:        "ping",
		Description: "Responds with pong.",
		Function: func(ctx *gommand.Context) error {
			_, _ = ctx.Reply("Pong!")
			return nil
		},
	})

	// A simple command to tag the user specified.
	router.SetCommand(&gommand.Command{
		Name:        "tag",
		Description: "Tags the user specified.",
		ArgTransformers: []gommand.ArgTransformer{
			{
				Function: gommand.UserTransformer,
			},
		},
		Function: func(ctx *gommand.Context) error {
			_, _ = ctx.Reply(ctx.Args[0].(*disgord.User).Mention())
			return nil
		},
	})

	// Echos one argument with an additional optional argument.
	router.SetCommand(&gommand.Command{
		Name:        "echo",
		Description: "Echos arguments.",
		ArgTransformers: []gommand.ArgTransformer{
			{
				Function: gommand.StringTransformer,
			},
			{
				Function: gommand.StringTransformer,
				Optional: true,
			},
		},
		Function: func(ctx *gommand.Context) error {
			_, _ = ctx.Reply(ctx.Args[0].(string))
			s, ok := ctx.Args[1].(string)
			if ok {
				_, _ = ctx.Reply("Optional arg: " + s)
			}
			return nil
		},
	})

	// Adds all of the numbers specified and then send the string at the end.
	router.SetCommand(&gommand.Command{
		Name:        "addandecho",
		Description: "Adds numbers and echos the last argument.",
		ArgTransformers: []gommand.ArgTransformer{
			{
				Function: gommand.IntTransformer,
				Greedy:   true,
			},
			{
				Function: gommand.StringTransformer,
			},
		},
		Function: func(ctx *gommand.Context) error {
			nums := ctx.Args[0].([]interface{})
			total := 0
			for _, v := range nums {
				total += v.(int)
			}
			_, _ = ctx.Reply(ctx.Args[1], ": ", total)
			return nil
		},
	})

	// Waits for a message from the user.
	router.SetCommand(&gommand.Command{
		Name:                 "echowait",
		Description:          "Wait for a message then echo it.",
		Function: func(ctx *gommand.Context) error {
			_, _ = ctx.Reply("say the message")
			resp := ctx.WaitForMessage(func(_ disgord.Session, msg *disgord.Message) bool {
				return msg.Author.ID == ctx.Message.Author.ID && msg.ChannelID == ctx.Message.ChannelID
			})
			_, _ = ctx.Reply(resp.Content)
			return nil
		},
	})

	// Handles command errors where possible. If not, just passes it through to the default handler to log to console.
	// Wanted to use Sentry? You could make a handler for this by capturing and returning false. Don't forget it's in the order if the handlers.
	router.AddErrorHandler(func(ctx *gommand.Context, err error) bool {
		// Check all the different types of errors.
		switch err.(type) {
		case *gommand.CommandNotFound, *gommand.CommandBlank:
			// We will ignore.
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
}

func main() {
	client := disgord.New(disgord.Config{
		BotToken: os.Getenv("TOKEN"),
		Logger:   disgord.DefaultLogger(false),
	})
	router.Hook(client)
	_ = client.StayConnectedUntilInterrupted(context.Background())
}
