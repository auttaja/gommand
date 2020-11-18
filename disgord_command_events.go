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
	if r.GetState != nil {
		if err := r.GetState(ctx); err != nil {
			r.errorHandler(ctx, err)
			return
		}
	}

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
		r.errorHandler(ctx, &CommandBlank{err: "The command is blank."})
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
			parser := ctx.Router.parserManager.Parser(reader)
			ok, err = r.CustomCommandsHandler(ctx, cmdname, parser)
			parser.Done()
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
	gateway := s.Gateway()
	gateway.Ready(r.readyEvt)
	gateway.UserUpdate(r.userUpdate)
	gateway.MessageCreate(r.msgCreate)

	if r.MessageCacheHandler != nil {
		gateway.GuildCreate(r.MessageCacheHandler.guildCreate)
		gateway.GuildDelete(r.MessageCacheHandler.guildDelete)
		gateway.ChannelDelete(r.MessageCacheHandler.channelDelete)
		gateway.MessageCreate(r.MessageCacheHandler.messageCreate)
		gateway.MessageUpdate(r.MessageCacheHandler.messageUpdate)
		gateway.MessageDelete(r.MessageCacheHandler.messageDelete)
		gateway.MessageDeleteBulk(r.MessageCacheHandler.bulkDeleteHandler)
	}

	gateway.MessageReactionAdd(handleMenuReactionEdit)
	gateway.MessageDelete(handleEmbedMenuMessageDelete)
}
