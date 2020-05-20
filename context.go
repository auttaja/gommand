package gommand

import (
	"context"
	"github.com/andersfylling/disgord"
)

// Context defines the information which might be required to run the command.
type Context struct {
	Prefix           string                 `json:"prefix"`
	Message          *disgord.Message       `json:"message"`
	BotUser          *disgord.User          `json:"botUser"`
	Router           *Router                `json:"-"`
	Session          disgord.Session        `json:"session"`
	Command          CommandInterface       `json:"command"`
	RawArgs          string                 `json:"rawArgs"`
	Args             []interface{}          `json:"args"`
	MiddlewareParams map[string]interface{} `json:"middlewareParams"`
}

// Replay is used to replay a command.
func (c *Context) Replay() error {
	c.Args = []interface{}{}
	return runCommand(c, &StringIterator{Text: c.RawArgs}, c.Command)
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

// WaitForMessage allows you to wait for a message.
func (c *Context) WaitForMessage(CheckFunc func(s disgord.Session, msg *disgord.Message) bool) *disgord.Message {
	c.Router.msgWaitingQueueLock.Lock()
	x := make(chan *disgord.Message)
	c.Router.msgWaitingQueue = append(c.Router.msgWaitingQueue, &msgQueueItem{
		function:  CheckFunc,
		goroutine: x,
	})
	c.Router.msgWaitingQueueLock.Unlock()
	return <-x
}

// DisplayEmbedMenu is used to allow you to easily display a embed menu.
func (c *Context) DisplayEmbedMenu(m *EmbedMenu) error {
	msg, err := c.Reply("Loading...")
	if err != nil {
		return err
	}
	return m.Display(c.Message.ChannelID, msg.ID, c.Session)
}
