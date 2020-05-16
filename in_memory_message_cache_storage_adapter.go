package gommand

import (
	"container/list"
	"github.com/andersfylling/disgord"
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
	cache map[string]*map[string]*cachedMessage
	list  *list.List
	len   uint

	guildLock  *sync.RWMutex
	channelMap map[string][]string
}

// Init is used to initialise the in-memory message cache.
func (c *InMemoryMessageCacheStorageAdapter) Init() {
	c.lock = &sync.RWMutex{}
	c.cache = map[string]*map[string]*cachedMessage{}
	c.list = list.New()

	c.guildLock = &sync.RWMutex{}
	c.channelMap = map[string][]string{}
}

// GetAllChannelIDs is used to get all of the channel ID's.
func (c *InMemoryMessageCacheStorageAdapter) GetAllChannelIDs(GuildID string) []string {
	c.guildLock.RLock()
	channels, _ := c.channelMap[GuildID]
	if channels == nil {
		channels = []string{}
	}
	c.guildLock.RUnlock()
	return channels
}

// AddChannelID is used to add a channel ID to the guild.
func (c *InMemoryMessageCacheStorageAdapter) AddChannelID(GuildID, ChannelID string) {
	c.guildLock.Lock()
	channels, _ := c.channelMap[GuildID]
	if channels == nil {
		channels = []string{}
	}
	c.channelMap[GuildID] = append(channels, ChannelID)
	c.guildLock.Unlock()
}

// RemoveChannelID is used to remove a channel ID to the guild.
func (c *InMemoryMessageCacheStorageAdapter) RemoveChannelID(GuildID, ChannelID string) {
	c.guildLock.Lock()
	channels, _ := c.channelMap[GuildID]
	if channels == nil {
		channels = []string{}
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
func (c *InMemoryMessageCacheStorageAdapter) RemoveGuild(GuildID string) {
	c.guildLock.Lock()
	delete(c.channelMap, GuildID)
	c.guildLock.Unlock()
}

// GetAndDelete is used to get and delete from the cache where this is possible.
func (c *InMemoryMessageCacheStorageAdapter) GetAndDelete(ChannelID, MessageID string) *disgord.Message {
	// Read lock the cache.
	c.lock.RLock()

	// Get the channel cache from the base.
	msgs, _ := c.cache[ChannelID]
	if msgs == nil {
		// Nope not in this cache, return nil.
		c.lock.RUnlock()
		return nil
	}

	// Try and get the compressed message from the cache.
	msg, _ := (*msgs)[MessageID]
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
func (c *InMemoryMessageCacheStorageAdapter) Delete(ChannelID, MessageID string) {
	// Write lock the cache.
	c.lock.Lock()

	// Get the channel cache from the base.
	msgs, _ := c.cache[ChannelID]
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
func (c *InMemoryMessageCacheStorageAdapter) DeleteChannelsMessages(ChannelID string) {
	// Write lock the cache.
	c.lock.Lock()

	// Get the channel cache and remove all messages in it.
	msgs, _ := c.cache[ChannelID]
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
func (c *InMemoryMessageCacheStorageAdapter) Set(ChannelID, MessageID string, Message *disgord.Message, Limit uint) {
	// Write lock the cache.
	c.lock.Lock()

	// Check if we are over the limit.
	if c.len == Limit && Limit != 0 {
		f := c.list.Front()
		if f != nil {
			c.lock.Unlock()
			s := strings.Split(f.Value.(string), "-")
			c.Delete(s[0], s[1])
			c.list.Remove(f)
			c.lock.Lock()
		}
	}

	// Set the message.
	msgs, _ := c.cache[ChannelID]
	if msgs == nil {
		m := map[string]*cachedMessage{}
		msgs = &m
		c.cache[ChannelID] = msgs
	}
	(*msgs)[MessageID] = &cachedMessage{
		msg: Message,
		el:  c.list.PushBack(ChannelID + "-" + MessageID),
	}
	c.len++

	// Write unlock the cache.
	c.lock.Unlock()
}
