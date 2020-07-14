package gommand

import (
	"github.com/andersfylling/disgord"
	"github.com/hako/durafmt"
	"sync"
	"time"
)

// Cooldown is used to define the interface which is used to handle command cooldowns.
type Cooldown interface {
	// Init is ran on the insertion of the cooldown bucket into a command/category/attribute which has been added to the router.
	// Note that your struct should implement logic to make sure that Init only modifies the struct once since a cooldown might be shared across objects!
	Init()

	// Check should add one X (where X is what you're measuring) to the cooldown bucket and return true if the command should run.
	Check(ctx *Context) (message string, ok bool)

	// Clear should clear any command related ratelimits.
	Clear()
}

// This is used so that any internal pointers for cooldowns will not point to a newly initialised map causing bugs.
type cooldownInternals struct {
	// Internal map of snowflake > usage.
	coolingDown map[disgord.Snowflake]uint

	// Internal lock used to modify/access the map.
	coolingDownLock *sync.Mutex
}

// Used to expire a usage.
func (i *cooldownInternals) expire(id disgord.Snowflake) {
	i.coolingDownLock.Lock()
	defer i.coolingDownLock.Unlock()
	usages, ok := i.coolingDown[id]
	if !ok {
		return
	}
	usages--
	if usages == 0 {
		delete(i.coolingDown, id)
	} else {
		i.coolingDown[id] = usages
	}
}

// Used to check a usage.
func (i *cooldownInternals) check(id disgord.Snowflake, max uint, expires time.Duration) (message string, shouldRun bool) {
	// Lock the mutex until we're done.
	i.coolingDownLock.Lock()
	defer i.coolingDownLock.Unlock()

	// Check how many usages by the guild.
	usages := i.coolingDown[id]

	// Is the usages equal to max runs? If so, return false.
	if usages == max {
		durationFmt := durafmt.Parse(expires).String()
		return "This command has a " + durationFmt + " cooldown.", false
	}

	// Add 1 to usages and set it.
	usages++
	i.coolingDown[id] = usages

	// Expire this usage.
	time.AfterFunc(expires, func() {
		i.expire(id)
	})

	// Return true here.
	return "", true
}

// GuildCooldown implements the Cooldown interface and is used to handle guild level ratelimits.
type GuildCooldown struct {
	// The internals used for cooldowns.
	internals *cooldownInternals

	// MaxRuns is the maximum amount of times a command can be ran in a guild until each usage expires.
	MaxRuns uint

	// UsageExpires is used to define how long until a usage expires (and is therefore not counted in the cooldown).
	UsageExpires time.Duration
}

// Init is used to initialise the guild cooldowns.
func (g *GuildCooldown) Init() {
	if g.internals != nil {
		// This has been initialised already.
		return
	}
	g.internals = &cooldownInternals{
		coolingDown:     map[disgord.Snowflake]uint{},
		coolingDownLock: &sync.Mutex{},
	}
}

// Check is used to check if the command should run and add 1 to the guild count.
func (g *GuildCooldown) Check(ctx *Context) (string, bool) {
	return g.internals.check(ctx.Message.GuildID, g.MaxRuns, g.UsageExpires)
}

// Clear is used to clear all cooldowns.
func (g *GuildCooldown) Clear() {
	oldLock := g.internals.coolingDownLock
	oldLock.Lock()
	g.internals.coolingDown = map[disgord.Snowflake]uint{}
	g.internals = &cooldownInternals{
		coolingDown:     map[disgord.Snowflake]uint{},
		coolingDownLock: &sync.Mutex{},
	}
	oldLock.Unlock()
}

// UserCooldown implements the Cooldown interface and is used to handle user level ratelimits.
type UserCooldown struct {
	// The internals used for cooldowns.
	internals *cooldownInternals

	// MaxRuns is the maximum amount of times a command can be ran by a user until each usage expires.
	MaxRuns uint

	// UsageExpires is used to define how long until a usage expires (and is therefore not counted in the cooldown).
	UsageExpires time.Duration
}

// Init is used to initialise the user cooldowns.
func (u *UserCooldown) Init() {
	if u.internals != nil {
		// This has been initialised already.
		return
	}
	u.internals = &cooldownInternals{
		coolingDown:     map[disgord.Snowflake]uint{},
		coolingDownLock: &sync.Mutex{},
	}
}

// Check is used to check if the command should run and add 1 to the user count.
func (u *UserCooldown) Check(ctx *Context) (string, bool) {
	return u.internals.check(ctx.Message.Author.ID, u.MaxRuns, u.UsageExpires)
}

// Clear is used to clear all cooldowns.
func (u *UserCooldown) Clear() {
	oldLock := u.internals.coolingDownLock
	oldLock.Lock()
	u.internals.coolingDown = map[disgord.Snowflake]uint{}
	u.internals = &cooldownInternals{
		coolingDown:     map[disgord.Snowflake]uint{},
		coolingDownLock: &sync.Mutex{},
	}
	oldLock.Unlock()
}

// ChannelCooldown implements the Cooldown interface and is used to handle channel level ratelimits.
type ChannelCooldown struct {
	// The internals used for cooldowns.
	internals *cooldownInternals

	// MaxRuns is the maximum amount of times a command can be ran by a channel until each usage expires.
	MaxRuns uint

	// UsageExpires is used to define how long until a usage expires (and is therefore not counted in the cooldown).
	UsageExpires time.Duration
}

// Init is used to initialise the channel cooldowns.
func (c *ChannelCooldown) Init() {
	if c.internals != nil {
		// This has been initialised already.
		return
	}
	c.internals = &cooldownInternals{
		coolingDown:     map[disgord.Snowflake]uint{},
		coolingDownLock: &sync.Mutex{},
	}
}

// Check is used to check if the command should run and add 1 to the channel count.
func (c *ChannelCooldown) Check(ctx *Context) (string, bool) {
	return c.internals.check(ctx.Message.ChannelID, c.MaxRuns, c.UsageExpires)
}

// Clear is used to clear all cooldowns.
func (c *ChannelCooldown) Clear() {
	oldLock := c.internals.coolingDownLock
	oldLock.Lock()
	c.internals.coolingDown = map[disgord.Snowflake]uint{}
	c.internals = &cooldownInternals{
		coolingDown:     map[disgord.Snowflake]uint{},
		coolingDownLock: &sync.Mutex{},
	}
	oldLock.Unlock()
}

// Handles multiple cooldowns.
type multiCooldownHandler struct {
	cooldowns []Cooldown
}

// Init is used to call Init on all cooldown handlers.
func (m *multiCooldownHandler) Init() {
	for _, v := range m.cooldowns {
		v.Init()
	}
}

// Clear is used to call Clear on all cooldown handlers.
func (m *multiCooldownHandler) Clear() {
	for _, v := range m.cooldowns {
		v.Clear()
	}
}

// Check is used to call Check on all cooldown handlers. If one returns false, we return the result of it.
func (m *multiCooldownHandler) Check(ctx *Context) (msg string, ok bool) {
	for _, v := range m.cooldowns {
		msg, ok = v.Check(ctx)
		if !ok {
			return
		}
	}
	return
}

// MultipleCooldowns is used to chain multiple cooldowns together.
func MultipleCooldowns(cooldowns ...Cooldown) Cooldown {
	return &multiCooldownHandler{cooldowns: cooldowns}
}
