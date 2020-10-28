package gommand

import (
	"github.com/andersfylling/disgord"
	"sync"
)

var (
	guildStates = map[disgord.Snowflake]*State{}
	stateLock   = sync.Mutex{}
)

// State represents the state for a guild.
type State struct {
	x    int
	lock *sync.RWMutex
}

// AddOne is used to add one to the state and return the result.
func (s *State) AddOne() int {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.x++
	return s.x
}

// GetValue is used to return the current value of the state.
func (s *State) GetValue() int {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.x
}

// Reset is used to reset the state and return the previous value.
func (s *State) Reset() int {
	s.lock.Lock()
	defer s.lock.Unlock()
	x := s.x
	s.x = 0
	return x
}

func setupState(ctx *Context) {
	stateLock.Lock()
	defer stateLock.Unlock()
	x, ok := guildStates[ctx.Message.GuildID]
	if !ok {
		x = &State{lock: &sync.RWMutex{}}
		guildStates[ctx.Message.GuildID] = x
	}
	ctx.State = x
}
