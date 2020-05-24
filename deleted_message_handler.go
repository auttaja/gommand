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
	GetAndDelete(ChannelID, MessageID string) *disgord.Message
	Delete(ChannelID, MessageID string)
	DeleteChannelsMessages(ChannelID string)
	Set(ChannelID, MessageID string, Message *disgord.Message, Limit uint)

	// Related to channel and guild ID relationship caching.
	// Channel ID's are NOT confirmed to be unique and will be repeated on bot reboot as per the Discord API.
	// You should manage this in your adapter.
	GetAllChannelIDs(GuildID string) []string
	AddChannelID(GuildID, ChannelID string)
	RemoveChannelID(GuildID, ChannelID string)
	RemoveGuild(GuildID string)
}

// DeletedMessageHandler is used to handle dispatching events for deleted messages.
// It does this by using the storage adapter to log messages, then the message is deleted from the database at the message limit or when the deleted message handler is called.
type DeletedMessageHandler struct {
	MessageCacheStorageAdapter MessageCacheStorageAdapter                    `json:"-"`
	Callback                   func(s disgord.Session, msg *disgord.Message) `json:"-"`

	// Limit defines the amount of messages.
	// -1 = unlimited (not suggested if it's in-memory since it'll lead to memory leaks), 0 = default, >0 = user set maximum
	Limit int `json:"limit"`
}

// Removes the guild from the cache.
func (d *DeletedMessageHandler) guildDelete(_ disgord.Session, evt *disgord.GuildDelete) {
	if evt.UnavailableGuild.Unavailable {
		// We shouldn't purge the guilds messages. The guild is simply down.
		return
	}
	go func() {
		ids := d.MessageCacheStorageAdapter.GetAllChannelIDs(evt.UnavailableGuild.ID.String())
		d.MessageCacheStorageAdapter.RemoveGuild(evt.UnavailableGuild.ID.String())
		for _, v := range ids {
			d.MessageCacheStorageAdapter.DeleteChannelsMessages(v)
		}
	}()
}

// Removes a channel from the cache.
func (d *DeletedMessageHandler) channelDelete(_ disgord.Session, evt *disgord.ChannelDelete) {
	go func() {
		gid := evt.Channel.GuildID.String()
		cid := evt.Channel.ID.String()
		d.MessageCacheStorageAdapter.RemoveChannelID(gid, cid)
		d.MessageCacheStorageAdapter.DeleteChannelsMessages(cid)
	}()
}

// Adds the guild to the cache.
func (d *DeletedMessageHandler) guildCreate(_ disgord.Session, evt *disgord.GuildCreate) {
	go func() {
		gid := evt.Guild.ID.String()
		for _, v := range evt.Guild.Channels {
			d.MessageCacheStorageAdapter.AddChannelID(gid, v.ID.String())
		}
	}()
}

// Defines the message deletion handler.
func (d *DeletedMessageHandler) messageDelete(s disgord.Session, evt *disgord.MessageDelete) {
	go func() {
		msg := d.MessageCacheStorageAdapter.GetAndDelete(evt.ChannelID.String(), evt.MessageID.String())
		if msg != nil {
			member, err := s.GetMember(context.TODO(), msg.GuildID, msg.Author.ID)
			if err != nil {
				return
			}
			member.GuildID = evt.GuildID
			msg.Member = member
			msg.Author = member.User
			d.Callback(s, msg)
		}
	}()
}

// Defines the message creation handler.
func (d *DeletedMessageHandler) messageCreate(_ disgord.Session, evt *disgord.MessageCreate) {
	Limit := d.Limit
	if Limit == 0 {
		Limit = defaultMessageCount
	} else if 0 > Limit {
		Limit = 0
	}
	go d.MessageCacheStorageAdapter.Set(evt.Message.ChannelID.String(), evt.Message.ID.String(), evt.Message, uint(Limit))
}
