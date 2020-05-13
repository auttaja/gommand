package gommand

import (
	"context"
	"github.com/andersfylling/disgord"
	"strings"
)

// Handles setting the bot user initially.
func (r *router) readyEvt(_ disgord.Session, evt *disgord.Ready) {
	// Set the bot user.
	r.BotUser = evt.User
}

// Handles changing the bot user if this is required.
func (r *router) userUpdate(_ disgord.Session, evt *disgord.UserUpdate) {
	// Check if it's the same user.
	r.cmdLock.RLock()
	if r.BotUser.ID != evt.User.ID {
		r.cmdLock.RUnlock()
		return
	}
	r.cmdLock.RUnlock()

	// Set the bot user.
	r.cmdLock.Lock()
	r.BotUser = evt.User
	r.cmdLock.Unlock()
}

// Handles processing new messages.
func (r *router) msgCreate(s disgord.Session, evt *disgord.MessageCreate) {
	// If the message is from a bot or this isn't in a guild, ignore it.
	if evt.Message.Author.Bot || evt.Message.IsDirectMessage() {
		return
	}

	// Read lock the commands.
	r.cmdLock.RLock()

	// Defer read unlocking commands.
	defer r.cmdLock.RUnlock()

	// Create the context.
	ctx := &Context{
		Message:          evt.Message,
		BotUser:          r.BotUser,
		Router:           r,
		Session:          s,
		Args:             []interface{}{},
		MiddlewareParams: map[string]interface{}{},
	}

	// Create a string iterator of the message content.
	reader := &StringIterator{Text: evt.Message.Content}

	// Run a prefix check.
	if !r.PrefixCheck(ctx, reader) {
		// The prefix was not used.
		return
	}

	// The member should be patched into the message object here to make it easier.
	member, err := s.GetMember(context.TODO(), evt.Message.GuildID, evt.Message.Author.ID)
	if err != nil {
		r.errorHandler(ctx, err)
		return
	}
	member.GuildID = evt.Message.GuildID
	evt.Message.Member = member

	// Iterate the message until the space.
	cmdname := ""
	for {
		b, err := reader.GetChar()
		if err != nil || b == ' ' {
			break
		}
		cmdname += string(b)
	}
	if cmdname == "" {
		var ok bool
		var err error
		if r.CustomCommandsHandler != nil {
			ok, err = r.CustomCommandsHandler(ctx, "", reader)
		}
		if err == nil {
			if !ok {
				r.errorHandler(ctx, &CommandBlank{err: "The command is blank."})
			}
		} else {
			r.errorHandler(ctx, err)
		}
		return
	}

	// Get the remainder as raw arguments.
	remainder, _ := reader.GetRemainder(false)
	ctx.RawArgs = remainder

	// Get the command if it exists.
	cmd, _ := r.cmds[strings.ToLower(cmdname)]
	ctx.Command = cmd
	if cmd == nil {
		var ok bool
		var err error
		if r.CustomCommandsHandler != nil {
			ok, err = r.CustomCommandsHandler(ctx, cmdname, reader)
		}
		if err == nil {
			if !ok {
				r.errorHandler(ctx, &CommandBlank{err: "The command \"" + cmdname + "\" does not exist."})
			}
		} else {
			r.errorHandler(ctx, err)
		}
		return
	}

	// Run the command handler.
	err = cmd.run(ctx, reader)
	if err != nil {
		r.errorHandler(ctx, err)
	}
}

// Hook is used to hook all required events into the disgord client.
func (r *router) Hook(s disgord.Session) {
	s.On(disgord.EvtReady, r.readyEvt)
	s.On(disgord.EvtUserUpdate, r.userUpdate)
	s.On(disgord.EvtMessageCreate, r.msgCreate)
}
