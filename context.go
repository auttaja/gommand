//go:generate go run generate_wait_for_event.go
package gommand

import (
	"context"
	"github.com/andersfylling/disgord"
	"strings"
)

// Context defines the information which might be required to run the command.
type Context struct {
	ShardID          uint                   `json:"shardId"`
	Prefix           string                 `json:"prefix"`
	Message          *disgord.Message       `json:"message"`
	BotUser          *disgord.User          `json:"botUser"`
	Router           *Router                `json:"-"`
	Session          disgord.Session        `json:"session"`
	Command          CommandInterface       `json:"command"`
	RawArgs          string                 `json:"rawArgs"`
	Args             []interface{}          `json:"args"`
	WaitManager      *WaitManager           `json:"-"`
	MiddlewareParams map[string]interface{} `json:"middlewareParams"`
}

// Replay is used to replay a command.
func (c *Context) Replay() error {
	c.Args = []interface{}{}
	return runCommand(c, strings.NewReader(c.RawArgs), c.Command)
}

// BotMember is used to get the bot as a member of the server this was within.
func (c *Context) BotMember() (*disgord.Member, error) {
	return c.Session.GetMember(context.TODO(), c.Message.GuildID, c.BotUser.ID)
}

// Guild is used to get the guild if the bot needs it.
func (c *Context) Guild() (*disgord.Guild, error) {
	return c.Session.GetGuild(context.TODO(), c.Message.Member.GuildID)
}

// Channel is used to get the channel if the bot needs it.
func (c *Context) Channel() (*disgord.Channel, error) {
	return c.Session.GetChannel(context.TODO(), c.Message.ChannelID)
}

// Reply is used to quickly reply to a command with a message.
func (c *Context) Reply(data ...interface{}) (*disgord.Message, error) {
	return c.Session.SendMsg(context.TODO(), c.Message.ChannelID, data...)
}

// PermissionVerifiedReply is used to reply to a command with a message.
// This is slower than the standard reply command since it checks it has permission first, but it also reduces the risk of a Cloudflare ban from the API.
func (c *Context) PermissionVerifiedReply(data ...interface{}) (*disgord.Message, error) {
	m, err := c.BotMember()
	if err != nil {
		return nil, err
	}
	channel, err := c.Channel()
	if err != nil {
		return nil, err
	}
	perms, err := channel.GetPermissions(context.TODO(), c.Session, m)
	if err != nil {
		return nil, err
	}
	required := disgord.PermissionSendMessages
	for _, v := range data {
		switch v.(type) {
		case disgord.Embed, *disgord.Embed:
			required |= disgord.PermissionEmbedLinks
		case disgord.CreateMessageFileParams, *disgord.CreateMessageFileParams:
			required |= disgord.PermissionAttachFiles
		case disgord.Message, *disgord.Message:
			embedlen := 0
			x, ok := v.(disgord.Message)
			if ok {
				embedlen = len(x.Embeds)
			} else {
				x := v.(*disgord.Message)
				embedlen = len(x.Embeds)
			}
			if embedlen >= 1 {
				required |= disgord.PermissionEmbedLinks
			}
		}
	}
	if (perms&required) != required && (perms&disgord.PermissionAdministrator) != disgord.PermissionAdministrator {
		g, err := c.Guild()
		if err != nil {
			return nil, err
		}
		if g.OwnerID != m.UserID {
			return nil, &disgord.ErrRest{
				Code: 403,
				Msg:  "Permissions check failed.",
			}
		}
	}
	return c.Reply(data...)
}

// WaitForMessage allows you to wait for a message. This function uses WaitForMessageCreate internally and is mainly here for simplicity and compatibility. You should NOT block during the check function.
func (c *Context) WaitForMessage(ctx context.Context, CheckFunc func(s disgord.Session, evt *disgord.Message) bool) *disgord.Message {
	x := c.WaitManager.WaitForMessageCreate(ctx, func(s disgord.Session, evt *disgord.MessageCreate) bool {
		return CheckFunc(s, evt.Message)
	})
	if x == nil {
		return nil
	}
	return x.Message
}

// DisplayEmbedMenu is used to allow you to easily display a embed menu.
func (c *Context) DisplayEmbedMenu(m *EmbedMenu) error {
	msg, err := c.Reply("Loading...")
	if err != nil {
		return err
	}
	return m.Display(c.Message.ChannelID, msg.ID, c.Session)
}
