package gommand

import (
	"container/list"
	"github.com/andersfylling/disgord"
	"github.com/andersfylling/snowflake/v4"
	"strings"
	"sync"
)

type cachedMessage struct {
	msg *disgord.Message
	el  *list.Element
}

// InMemoryMessageCacheStorageAdapter is used to hold cached messages in RAM. This is extremely fast, but will lead to increased RAM usage.
type InMemoryMessageCacheStorageAdapter struct {
	lock  *sync.RWMutex
	cache map[snowflake.Snowflake]*map[snowflake.Snowflake]*cachedMessage
	list  *list.List
	len   uint

	guildLock  *sync.RWMutex
	channelMap map[snowflake.Snowflake][]snowflake.Snowflake
}

// Init is used to initialise the in-memory message cache.
func (c *InMemoryMessageCacheStorageAdapter) Init() {
	c.lock = &sync.RWMutex{}
	c.cache = map[snowflake.Snowflake]*map[snowflake.Snowflake]*cachedMessage{}
	c.list = list.New()

	c.guildLock = &sync.RWMutex{}
	c.channelMap = map[snowflake.Snowflake][]snowflake.Snowflake{}
}

// GetAllChannelIDs is used to get all of the channel ID's.
func (c *InMemoryMessageCacheStorageAdapter) GetAllChannelIDs(GuildID snowflake.Snowflake) []snowflake.Snowflake {
	c.guildLock.RLock()
	channels := c.channelMap[GuildID]
	if channels == nil {
		channels = []snowflake.Snowflake{}
	}
	c.guildLock.RUnlock()
	return channels
}

// AddChannelID is used to add a channel ID to the guild.
func (c *InMemoryMessageCacheStorageAdapter) AddChannelID(GuildID, ChannelID snowflake.Snowflake) {
	c.guildLock.Lock()
	channels := c.channelMap[GuildID]
	if channels == nil {
		channels = []snowflake.Snowflake{}
	}
	c.channelMap[GuildID] = append(channels, ChannelID)
	c.guildLock.Unlock()
}

// RemoveChannelID is used to remove a channel ID to the guild.
func (c *InMemoryMessageCacheStorageAdapter) RemoveChannelID(GuildID, ChannelID snowflake.Snowflake) {
	c.guildLock.Lock()
	channels := c.channelMap[GuildID]
	if channels == nil {
		channels = []snowflake.Snowflake{}
	}
	for i, v := range channels {
		if v != ChannelID {
			continue
		}
		channels[i] = channels[len(channels)-1]
		channels = channels[:len(channels)-1]
		break
	}
	c.channelMap[GuildID] = channels
	c.guildLock.Unlock()
}

// RemoveGuild is used to remove a guild from the cache.
func (c *InMemoryMessageCacheStorageAdapter) RemoveGuild(GuildID snowflake.Snowflake) {
	c.guildLock.Lock()
	delete(c.channelMap, GuildID)
	c.guildLock.Unlock()
}

// GetAndDelete is used to get and delete from the cache where this is possible.
func (c *InMemoryMessageCacheStorageAdapter) GetAndDelete(ChannelID, MessageID snowflake.Snowflake) *disgord.Message {
	// Read lock the cache.
	c.lock.RLock()

	// Get the channel cache from the base.
	msgs := c.cache[ChannelID]
	if msgs == nil {
		// Nope not in this cache, return nil.
		c.lock.RUnlock()
		return nil
	}

	// Try and get the compressed message from the cache.
	msg := (*msgs)[MessageID]
	c.lock.RUnlock()
	if msg == nil {
		// Nothing to delete.
		return nil
	}

	// Delete the message from the cache.
	c.Delete(ChannelID, MessageID)

	// Return the message.
	return msg.msg
}

// Delete is used to delete a specific message from the cache.
func (c *InMemoryMessageCacheStorageAdapter) Delete(ChannelID, MessageID snowflake.Snowflake) {
	// Write lock the cache.
	c.lock.Lock()

	// Get the channel cache from the base.
	msgs := c.cache[ChannelID]
	if msgs == nil {
		// Nope not in this cache, return.
		c.lock.Unlock()
		return
	}

	// Check the message exists in the cache.
	msg, ok := (*msgs)[MessageID]
	if !ok {
		// Not here, return.
		c.lock.Unlock()
		return
	}

	// Delete from the cache.
	delete(*msgs, MessageID)

	// Check the length of messages.
	if len(*msgs) == 0 {
		// Remove the channel cache from the base cache.
		delete(c.cache, ChannelID)
	}

	// Remove the message from the list.
	c.list.Remove(msg.el)

	// Remove 1 from the length.
	c.len--

	// Write unlock the cache.
	c.lock.Unlock()
}

// DeleteChannelsMessages is used to delete a channels messages from a cache.
func (c *InMemoryMessageCacheStorageAdapter) DeleteChannelsMessages(ChannelID snowflake.Snowflake) {
	// Write lock the cache.
	c.lock.Lock()

	// Get the channel cache and remove all messages in it.
	msgs := c.cache[ChannelID]
	for _, v := range *msgs {
		c.list.Remove(v.el)
		c.len--
	}

	// Delete the channel from the cache.
	delete(c.cache, ChannelID)

	// Write unlock the cache.
	c.lock.Unlock()
}

// Set is used to set a item in the cache.
func (c *InMemoryMessageCacheStorageAdapter) Set(ChannelID, MessageID snowflake.Snowflake, Message *disgord.Message, Limit uint) {
	// Write lock the cache.
	c.lock.Lock()

	// Check if we are over the limit.
	if c.len == Limit && Limit != 0 {
		f := c.list.Front()
		if f != nil {
			c.lock.Unlock()
			s := strings.Split(f.Value.(string), "-")
			c.Delete(snowflake.ParseSnowflakeString(s[0]), snowflake.ParseSnowflakeString(s[1]))
			c.list.Remove(f)
			c.lock.Lock()
		}
	}

	// Set the message.
	msgs := c.cache[ChannelID]
	if msgs == nil {
		m := map[snowflake.Snowflake]*cachedMessage{}
		msgs = &m
		c.cache[ChannelID] = msgs
	}
	(*msgs)[MessageID] = &cachedMessage{
		msg: Message,
		el:  c.list.PushBack(ChannelID.String() + "-" + MessageID.String()),
	}
	c.len++

	// Write unlock the cache.
	c.lock.Unlock()
}
