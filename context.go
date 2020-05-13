package gommand

import (
	"context"
	"github.com/andersfylling/disgord"
)

// Context defines the information which might be required to run the command.
type Context struct {
	Message          *disgord.Message       `json:"message"`
	BotUser          *disgord.User          `json:"botUser"`
	Router           *Router                `json:"-"`
	Session          disgord.Session        `json:"session"`
	Command          *Command               `json:"command"`
	RawArgs          string                 `json:"rawArgs"`
	Args             []interface{}          `json:"args"`
	MiddlewareParams map[string]interface{} `json:"middlewareParams"`
}

// Replay is used to replay a command.
func (c *Context) Replay() error {
	c.Args = []interface{}{}
	return c.Command.run(c, &StringIterator{Text: c.RawArgs})
}

// BotMember is used to get the bot as a member of the server this was within.
func (c *Context) BotMember() (*disgord.Member, error) {
	return c.Session.GetMember(context.TODO(), c.Message.GuildID, c.BotUser.ID)
}

// Channel is used to get the channel if the bot needs it.
func (c *Context) Channel() (*disgord.Channel, error) {
	return c.Session.GetChannel(context.TODO(), c.Message.ChannelID)
}

// Reply is used to quickly reply to a command with a message.
func (c *Context) Reply(data ...interface{}) (*disgord.Message, error) {
	return c.Session.SendMsg(context.TODO(), c.Message.ChannelID, data...)
}
