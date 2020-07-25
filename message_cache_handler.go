package gommand

import (
	"context"
	"github.com/andersfylling/disgord"
)

// The default number of messages that will be cached by the storage adapter.
var defaultMessageCount = 1000

// MessageCacheStorageAdapter is the interface which is used for message cache storage adapters.
type MessageCacheStorageAdapter interface {
	// Called when the router is created.
	Init()

	// Related to message caching.
	GetAndDelete(ChannelID, MessageID disgord.Snowflake) *disgord.Message
	Delete(ChannelID, MessageID disgord.Snowflake)
	DeleteChannelsMessages(ChannelID disgord.Snowflake)
	Set(ChannelID, MessageID disgord.Snowflake, Message *disgord.Message, Limit uint)
	Update(ChannelID, MessageID disgord.Snowflake, Message *disgord.Message) (old *disgord.Message)

	// Handles guild removal. The behaviour of this changes depending on if GuildChannelRelationshipManagement is implemented.
	// If it is, this will just be used to remove all guild/channel relationships but not messages from the cache (that'll be done by running DeleteChannelsMessages with each channel ID).
	// If it isn't, it will remove all of the guilds messages from the cache.
	RemoveGuild(GuildID disgord.Snowflake)
}

// GuildChannelRelationshipManagement are an optional set of functions which a struct implementing MessageCacheStorageAdapter can use to manage channel/guild ID relationships.
type GuildChannelRelationshipManagement interface {
	GetAllChannelIDs(GuildID disgord.Snowflake) []disgord.Snowflake
	AddChannelID(GuildID, ChannelID disgord.Snowflake)
	RemoveChannelID(GuildID, ChannelID disgord.Snowflake)
}

// MessageCacheHandler is used to handle dispatching events for deleted/edited messages.
// It does this by using the storage adapter to log messages, then the message is deleted from the database at the message limit or when the deleted message handler is called.
type MessageCacheHandler struct {
	MessageCacheStorageAdapter MessageCacheStorageAdapter                       `json:"-"`
	DeletedCallback     func(s disgord.Session, msg *disgord.Message)           `json:"-"`
	UpdatedCallback     func(s disgord.Session, before, after *disgord.Message)	`json:"-"`

	// Limit defines the amount of messages.
	// -1 = unlimited (not suggested if it's in-memory since it'll lead to memory leaks), 0 = default, >0 = user set maximum
	Limit int `json:"limit"`

	// IgnoreBots is whether or not messages from bots should be excluded from the message cache.
	IgnoreBots bool `json:"ignoreBots"`
}

// Removes the guild from the cache.
func (d *MessageCacheHandler) guildDelete(_ disgord.Session, evt *disgord.GuildDelete) {
	if evt.UnavailableGuild.Unavailable {
		// We shouldn't purge the guilds messages. The guild is simply down.
		return
	}
	go func() {
		r, ok := d.MessageCacheStorageAdapter.(GuildChannelRelationshipManagement)
		var ids []disgord.Snowflake
		if ok {
			ids = r.GetAllChannelIDs(evt.UnavailableGuild.ID)
		}
		d.MessageCacheStorageAdapter.RemoveGuild(evt.UnavailableGuild.ID)
		if ok {
			for _, v := range ids {
				d.MessageCacheStorageAdapter.DeleteChannelsMessages(v)
			}
		}
	}()
}

// Removes a channel from the cache.
func (d *MessageCacheHandler) channelDelete(_ disgord.Session, evt *disgord.ChannelDelete) {
	go func() {
		gid := evt.Channel.GuildID
		cid := evt.Channel.ID
		r, ok := d.MessageCacheStorageAdapter.(GuildChannelRelationshipManagement)
		if ok {
			r.RemoveChannelID(gid, cid)
		}
		d.MessageCacheStorageAdapter.DeleteChannelsMessages(cid)
	}()
}

// Adds the guild to the cache.
func (d *MessageCacheHandler) guildCreate(_ disgord.Session, evt *disgord.GuildCreate) {
	go func() {
		gid := evt.Guild.ID
		r, ok := d.MessageCacheStorageAdapter.(GuildChannelRelationshipManagement)
		if ok {
			for _, v := range evt.Guild.Channels {
				r.AddChannelID(gid, v.ID)
			}
		}
	}()
}

// Defines the message deletion handler.
func (d *MessageCacheHandler) messageDelete(s disgord.Session, evt *disgord.MessageDelete) {
	go func() {
		msg := d.MessageCacheStorageAdapter.GetAndDelete(evt.ChannelID, evt.MessageID)
		if msg != nil {
			member, err := s.GetMember(context.TODO(), msg.GuildID, msg.Author.ID)
			if err != nil {
				return
			}
			member.GuildID = evt.GuildID
			msg.Member = member
			msg.Author = member.User
			d.DeletedMessageCallback(s, msg)
		}
	}()
}

// Defines the message creation handler.
func (d *MessageCacheHandler) messageCreate(_ disgord.Session, evt *disgord.MessageCreate) {
	if d.IgnoreBots && evt.Message.Author.Bot {
		return
	}
	Limit := d.Limit
	if Limit == 0 {
		Limit = defaultMessageCount
	} else if 0 > Limit {
		Limit = 0
	}
	go d.MessageCacheStorageAdapter.Set(evt.Message.ChannelID, evt.Message.ID, evt.Message, uint(Limit))
}

// Defines the message update handler.
func (d *MessageCacheHandler) messageUpdate(s disgord.Session, evt *disgord.MessageUpdate) {
	if d.IgnoreBots && evt.Message.Author.Bot {
		return
	}
	go func() {
		before := d.MessageCacheStorageAdapter.Update(evt.Message.ChannelID, evt.Message.ID, evt.Message)
		if before != nil {
			d.UpdatedMessageCallback(s, before, evt.Message)
		}
	}()
}
