package gommand

import (
	"github.com/andersfylling/disgord"
	"io"
	"strings"
)

// Handles setting the bot user initially.
func (r *Router) readyEvt(_ disgord.Session, evt *disgord.Ready) {
	// Set the bot user.
	r.cmdLock.Lock()
	r.botUsers[evt.ShardID] = evt.User
	r.cmdLock.Unlock()
}

// Handles changing the bot user if this is required.
func (r *Router) userUpdate(_ disgord.Session, evt *disgord.UserUpdate) {
	// Check if it's the same user.
	r.cmdLock.RLock()
	if r.botUsers[evt.ShardID].ID != evt.User.ID {
		r.cmdLock.RUnlock()
		return
	}
	r.cmdLock.RUnlock()

	// Set the bot user.
	r.cmdLock.Lock()
	r.botUsers[evt.ShardID] = evt.User
	r.cmdLock.Unlock()
}

// CommandProcessor is used to do the message command processing.
func (r *Router) CommandProcessor(s disgord.Session, ShardID uint, msg *disgord.Message, prefix bool) {
	// If the message is from a bot or this isn't in a guild, ignore it.
	if msg.Author.Bot || msg.IsDirectMessage() {
		return
	}

	// Read lock the commands.
	r.cmdLock.RLock()

	// Create the context.
	ctx := &Context{
		ShardID:          ShardID,
		Message:          msg,
		BotUser:          r.botUsers[ShardID],
		Router:           r,
		Session:          s,
		Args:             []interface{}{},
		MiddlewareParams: map[string]interface{}{},
	}
	ctx.WaitManager = &WaitManager{ctx: ctx}

	// Create a read seeker of the message content.
	reader := strings.NewReader(msg.Content)

	// Run a prefix check.
	if prefix {
		if !r.PrefixCheck(ctx, reader) {
			// The prefix was not used.
			r.cmdLock.RUnlock()
			return
		}
	}

	// Parts of the member should be patched into the message object here to make it easier to use.
	msg.Member.GuildID = msg.GuildID
	msg.Member.User = msg.Author

	// Iterate the message until the space.
	cmdname := ""
	for {
		b, err := reader.ReadByte()
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
	p := r.parserManager.Parser(reader)
	remainder, _ := p.Remainder()
	p.Done()
	_, _ = reader.Seek(int64(len(remainder)*-1), io.SeekCurrent)
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
	// Launch the handler in a go-routine.
	go r.CommandProcessor(s, evt.ShardID, evt.Message, true)
}

// Hook is used to hook all required events into the disgord client.
func (r *Router) Hook(s disgord.Session) {
	s.On(disgord.EvtReady, r.readyEvt)
	s.On(disgord.EvtUserUpdate, r.userUpdate)
	s.On(disgord.EvtMessageCreate, r.msgCreate)
	if r.MessageCacheHandler != nil {
		s.On(disgord.EvtGuildDelete, r.MessageCacheHandler.guildDelete)
		s.On(disgord.EvtChannelDelete, r.MessageCacheHandler.channelDelete)
		s.On(disgord.EvtGuildCreate, r.MessageCacheHandler.guildCreate)
		s.On(disgord.EvtMessageCreate, r.MessageCacheHandler.messageCreate)
		s.On(disgord.EvtMessageDelete, r.MessageCacheHandler.messageDelete)
		s.On(disgord.EvtMessageUpdate, r.MessageCacheHandler.messageUpdate)
	}
	s.On(disgord.EvtMessageReactionAdd, handleMenuReactionEdit)
	s.On(disgord.EvtMessageDelete, handleEmbedMenuMessageDelete)
}
