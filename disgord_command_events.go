package gommand

import (
	"context"
	"github.com/andersfylling/disgord"
	"strings"
)

// Defines if the member should be patched in.
// This is true by default, and should only be made false for unit tests.
var patchMember = true

// Handles setting the bot user initially.
func (r *Router) readyEvt(_ disgord.Session, evt *disgord.Ready) {
	// Set the bot user.
	r.BotUser = evt.User
}

// Handles changing the bot user if this is required.
func (r *Router) userUpdate(_ disgord.Session, evt *disgord.UserUpdate) {
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

// CommandProcessor is used to do the message command processing.
func (r *Router) CommandProcessor(s disgord.Session, msg *disgord.Message) {
	// If the message is from a bot or this isn't in a guild, ignore it.
	if msg.Author.Bot || msg.IsDirectMessage() {
		return
	}

	// Read lock the commands.
	r.cmdLock.RLock()

	// Create the context.
	ctx := &Context{
		Message:          msg,
		BotUser:          r.BotUser,
		Router:           r,
		Session:          s,
		Args:             []interface{}{},
		MiddlewareParams: map[string]interface{}{},
	}

	// Create a string iterator of the message content.
	reader := &StringIterator{Text: msg.Content}

	// Run a prefix check.
	if !r.PrefixCheck(ctx, reader) {
		// The prefix was not used.
		r.cmdLock.RUnlock()
		return
	}

	// The member should be patched into the message object here to make it easier.
	if patchMember {
		member, err := s.GetMember(context.TODO(), msg.GuildID, msg.Author.ID)
		if err != nil {
			r.cmdLock.RUnlock()
			r.errorHandler(ctx, err)
			return
		}
		member.GuildID = msg.GuildID
		msg.Member = member
	}

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
		r.cmdLock.RUnlock()
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
	cmd := r.cmds[strings.ToLower(cmdname)]
	ctx.Command = cmd
	if cmd == nil {
		r.cmdLock.RUnlock()
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
	r.cmdLock.RUnlock()
	err := runCommand(ctx, reader, cmd)
	if err != nil {
		r.errorHandler(ctx, err)
	}
}

// Handles processing new messages.
func (r *Router) msgCreate(s disgord.Session, evt *disgord.MessageCreate) {
	// Create a go-routine to handle message waiting.
	go func() {
		// Lock the message queue.
		// This is not a R/W mutex because of potential race conditions if we have to read unlock and then lock.
		r.msgWaitingQueueLock.Lock()

		// Create the index array for the items we want to purge.
		indexes := make([]int, 0, 1)

		// Iterate through the message queue.
		for i, v := range r.msgWaitingQueue {
			if v.function(s, evt.Message) {
				// This event can now be dropped from the queue.
				indexes = append(indexes, i)
				v.goroutine <- evt.Message
			}
		}

		// Remove each index from the array.
		for _, i := range indexes {
			r.msgWaitingQueue[i] = r.msgWaitingQueue[len(r.msgWaitingQueue)-1]
			r.msgWaitingQueue = r.msgWaitingQueue[:len(r.msgWaitingQueue)-1]
		}

		// Unlock the message queue.
		r.msgWaitingQueueLock.Unlock()
	}()

	// Launch the handler in a go-routine.
	go r.CommandProcessor(s, evt.Message)
}

// Hook is used to hook all required events into the disgord client.
func (r *Router) Hook(s disgord.Session) {
	s.On(disgord.EvtReady, r.readyEvt)
	s.On(disgord.EvtUserUpdate, r.userUpdate)
	s.On(disgord.EvtMessageCreate, r.msgCreate)
	if r.DeletedMessageHandler != nil {
		s.On(disgord.EvtGuildDelete, r.DeletedMessageHandler.guildDelete)
		s.On(disgord.EvtChannelDelete, r.DeletedMessageHandler.channelDelete)
		s.On(disgord.EvtGuildCreate, r.DeletedMessageHandler.guildCreate)
		s.On(disgord.EvtMessageCreate, r.DeletedMessageHandler.messageCreate)
		s.On(disgord.EvtMessageDelete, r.DeletedMessageHandler.messageDelete)
	}
	s.On(disgord.EvtMessageReactionAdd, handleMenuReactionEdit)
	s.On(disgord.EvtMessageDelete, handleEmbedMenuMessageDelete)
}
