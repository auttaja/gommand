package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/andersfylling/disgord"
	"github.com/auttaja/gommand"
	"github.com/sirupsen/logrus"
)

// Create the command router.
var router = gommand.NewRouter(&gommand.RouterConfig{
	// The prefix function should be set here or it will be blank.
	// We are using % and mention prefixes for this example.
	PrefixCheck: gommand.MultiplePrefixCheckers(gommand.StaticPrefix("%"), gommand.MentionPrefix),

	// Prints deleted/edited messages.
	MessageCacheHandler: &gommand.MessageCacheHandler{
		DeletedCallback: func(s disgord.Session, msg *disgord.Message) {
			println(msg.Content)
		},
		UpdatedCallback: func(s disgord.Session, before, after *disgord.Message) {
			println(before.Content, ">", after.Content)
		},
		BulkDeletedCallback: func(s disgord.Session, channelID disgord.Snowflake, messages []*disgord.Message) {
			builder := &strings.Builder{}
			for _, m := range messages {
				builder.WriteString(fmt.Sprintf("[Timestamp: %s] [Author: %s] [ID: %s] [Content: %s]\n",
					m.Timestamp.Format("Mon Jan 2 15:04:05"),
					m.Author.String(),
					m.ID.String(),
					m.Content,
				))
			}

			_, _ = s.SendMsg(channelID, &disgord.CreateMessageParams{
				Content: fmt.Sprintf("Bulk delete in <#%s>. %d deleted messages logged.", channelID.String(), len(messages)),
				Files: []disgord.CreateMessageFileParams{
					{
						Reader:   strings.NewReader(builder.String()),
						FileName: channelID.String() + "_deleted.txt",
					},
				},
			})
		},
	},
})

// An example of building your own struct for a command.
type ping struct {
	gommand.CommandBasics
}

func (p *ping) Init() {
	p.Name = "ping"
	p.Description = "Responds with pong."
}

func (ping) CommandFunction(ctx *gommand.Context) error {
	_, _ = ctx.Reply("Pong!")
	return nil
}

func init() {
	// A simple command to respond with pong.
	router.SetCommand(&ping{})

	// A simple command to tag the user specified.
	router.SetCommand(&gommand.Command{
		Name:        "tag",
		Description: "Tags the user specified.",
		Cooldown: &gommand.UserCooldown{
			MaxRuns:      2,
			UsageExpires: time.Minute,
		},
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
		Name:        "echowait",
		Description: "Wait for a message then echo it.",
		Function: func(ctx *gommand.Context) error {
			_, _ = ctx.Reply("say the message")
			c, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
			defer cancel()
			resp := ctx.WaitForMessage(c, func(_ disgord.Session, msg *disgord.Message) bool {
				return msg.Author.ID == ctx.Message.Author.ID && msg.ChannelID == ctx.Message.ChannelID
			})
			if resp == nil {
				_, _ = ctx.Reply("timed out")
			} else {
				_, _ = ctx.Reply(resp.Content)
			}
			return nil
		},
	})

	// Waits for a reaction from the user.
	router.SetCommand(&gommand.Command{
		Name:        "reactionwait",
		Description: "Wait for a reaction from the user then echo it.",
		Function: func(ctx *gommand.Context) error {
			msg, err := ctx.Reply("react with what you want echoed")
			if err != nil {
				return nil
			}
			c, cancel := context.WithTimeout(context.TODO(), time.Minute)
			defer cancel()
			resp := ctx.WaitManager.WaitForMessageReactionAdd(c, func(_ disgord.Session, evt *disgord.MessageReactionAdd) bool {
				return evt.UserID == ctx.Message.Author.ID && evt.MessageID == msg.ID
			})
			if resp == nil {
				_, _ = ctx.Reply("timed out")
			} else {
				_, _ = ctx.Reply(resp.PartialEmoji)
			}
			return nil
		},
	})

	// Create a basic embed menu.
	router.SetCommand(&gommand.Command{
		Name:        "embedmenu",
		Description: "Displays a embed menu.",
		Function: func(ctx *gommand.Context) error {
			menu := gommand.NewEmbedMenu(&disgord.Embed{
				Title:       "Want to see doge?",
				Description: "Click the option below.",
			}, ctx)

			child := menu.NewChildMenu(&gommand.ChildMenuOptions{
				Embed: &disgord.Embed{
					Image: &disgord.EmbedImage{
						URL: "https://cdn.vox-cdn.com/thumbor/s6HznC4HCYrV3axUS-7wVOPbC2c=/0x0:1020x680/2050x1367/cdn.vox-cdn.com/assets/3785529/DOGE-10.jpg",
					},
				},
				Button: &gommand.MenuButton{
					Emoji:       "🇦",
					Name:        "Show doge",
					Description: "such doge such wow",
				},
				AfterAction: func() {
					println("Doge was here!")
				},
			})
			child.AddBackButton()

			_ = ctx.DisplayEmbedMenu(menu)

			return nil
		},
	})

	// Creates embed pages.
	router.SetCommand(&gommand.Command{
		Name:        "embedpages",
		Description: "Shows embed pages.",
		Function: func(ctx *gommand.Context) error {
			_ = gommand.EmbedsPaginator(ctx, []*disgord.Embed{
				{
					Title: "Hi",
				},
				{
					Title: "World",
				},
				{
					Image: &disgord.EmbedImage{
						URL: "https://cdn.vox-cdn.com/thumbor/s6HznC4HCYrV3axUS-7wVOPbC2c=/0x0:1020x680/2050x1367/cdn.vox-cdn.com/assets/3785529/DOGE-10.jpg",
					},
				},
			}, 0, "")
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
	logrus.SetLevel(logrus.DebugLevel)
	client := disgord.New(disgord.Config{
		BotToken: os.Getenv("TOKEN"),
		Logger:   logrus.New(),
	})
	router.Hook(client)
	_ = client.StayConnectedUntilInterrupted(context.Background())
}
