package gommand

import (
	"github.com/andersfylling/disgord"
	"strings"
	"sync"
)

// PrefixCheck is the type for a function to check the prefix. true here means the prefix is there and was read.
// Note that prefix functions cannot see the member object in the message since that is patched in AFTER this is ran.
type PrefixCheck = func(ctx *Context, r *StringIterator) bool

// ErrorHandler is a function which is used to handle errors.
// true here means that the error was handled by this handler. If this is false, the processor will go to the next handler.
type ErrorHandler = func(ctx *Context, err error) bool

// CustomCommandsHandler is the handler which is used for custom commands.
// An error being returned by the custom commands handler will return in that being passed through to the error handler instead.
// true here represents this being a custom command. This means the Router will not go through the errors handler unless an error is set.
type CustomCommandsHandler = func(ctx *Context, cmdname string, r *StringIterator) (bool, error)

// PermissionValidator is a function which is used to validate someones permissions who is running a command.
// If the boolean is set to true, it means the user has permissions.
// If it is false, the user does not, and the string will be used as the error with the IncorrectPermissions type.
type PermissionValidator = func(ctx *Context) (string, bool)

// Middleware is used to handle setting data within the context before the command is ran.
// This is useful if there is stuff you want to do, but you don't want to type on every command.
// If this errors, it will just get passed through to the error handler.
type Middleware = func(ctx *Context) error

// RouterConfig defines the config which will be used for the Router.
type RouterConfig struct {
	DeletedMessageHandler *DeletedMessageHandler
	PrefixCheck           PrefixCheck
	ErrorHandlers         []ErrorHandler
	PermissionValidators  []PermissionValidator
	Middleware            []Middleware
}

type msgQueueItem struct {
	function  func(s disgord.Session, msg *disgord.Message) bool
	goroutine chan *disgord.Message
}

// Router defines the command router which is being used.
// Please call NewRouter to initialise this rather than creating a new struct.
type Router struct {
	BotUser               *disgord.User         `json:"-"`
	PrefixCheck           PrefixCheck           `json:"-"`
	CustomCommandsHandler CustomCommandsHandler `json:"-"`
	cmds                  map[string]*Command
	cmdLock               *sync.RWMutex
	errorHandlers         []ErrorHandler
	permissionValidators  []PermissionValidator
	middleware            []Middleware
	msgWaitingQueue       []*msgQueueItem
	msgWaitingQueueLock   *sync.Mutex
	DeletedMessageHandler *DeletedMessageHandler
}

// NewRouter creates a new command Router.
func NewRouter(Config *RouterConfig) *Router {
	if Config.PrefixCheck == nil {
		// No prefix.
		Config.PrefixCheck = func(ctx *Context, r *StringIterator) bool {
			return true
		}
	}

	// Set the Router.
	r := &Router{
		PrefixCheck:           Config.PrefixCheck,
		cmds:                  map[string]*Command{},
		cmdLock:               &sync.RWMutex{},
		errorHandlers:         Config.ErrorHandlers,
		permissionValidators:  Config.PermissionValidators,
		middleware:            Config.Middleware,
		DeletedMessageHandler: Config.DeletedMessageHandler,
		msgWaitingQueue:       []*msgQueueItem{},
		msgWaitingQueueLock:   &sync.Mutex{},
	}

	// If deleted message handler isn't nil, initialise the storage adapter.
	if r.DeletedMessageHandler != nil {
		if r.DeletedMessageHandler.MessageCacheStorageAdapter == nil {
			r.DeletedMessageHandler.MessageCacheStorageAdapter = &InMemoryMessageCacheStorageAdapter{}
		}
		r.DeletedMessageHandler.MessageCacheStorageAdapter.Init()
	}

	// Return the router.
	return r
}

// Dispatches the required error handlers in the event of an error.
func (r *Router) errorHandler(ctx *Context, err error) {
	if r.errorHandlers != nil {
		// There are error handlers, go through these first.
		for _, v := range r.errorHandlers {
			if v(ctx, err) {
				// The error has been handled. Return here.
				return
			}
		}
	}

	// No error handlers, go straight to console logging.
	ctx.Session.Logger().Error(err)
}

// AddErrorHandler is used to add a error handler to the Router.
// An error handler takes the structure of the ErrorHandler type above.
// Error handlers are executed in the order in which they are added to the Router.
func (r *Router) AddErrorHandler(Handler ErrorHandler) {
	r.cmdLock.Lock()
	r.errorHandlers = append(r.errorHandlers, Handler)
	r.cmdLock.Unlock()
}

// GetCommand is used to get a command from the Router if it exists.
// If the command doesn't exist, this will be a nil pointer.
func (r *Router) GetCommand(Name string) *Command {
	r.cmdLock.RLock()
	cmd := r.cmds[strings.ToLower(Name)]
	r.cmdLock.RUnlock()
	return cmd
}

// SetCommand is used to set a command.
func (r *Router) SetCommand(c *Command) {
	r.cmdLock.Lock()
	c.Name = strings.ToLower(c.Name)
	r.cmds[c.Name] = c
	if c.Aliases == nil {
		c.Aliases = []string{}
	}
	for i, v := range c.Aliases {
		v = strings.ToLower(v)
		c.Aliases[i] = v
		r.cmds[v] = c
	}
	r.cmdLock.Unlock()
}

// RemoveCommand is used to remove a command from the Router.
func (r *Router) RemoveCommand(c *Command) {
	r.cmdLock.Lock()
	delete(r.cmds, strings.ToLower(c.Name))
	if c.Aliases != nil {
		for _, v := range c.Aliases {
			delete(r.cmds, strings.ToLower(v))
		}
	}
	r.cmdLock.Unlock()
}

// GetAllCommands is used to get all of the commands.
func (r *Router) GetAllCommands() []*Command {
	// Read lock the commands.
	r.cmdLock.RLock()

	// Get the command count.
	count := 0
	for k, v := range r.cmds {
		if k == v.Name {
			count++
		}
	}

	// Allocate the commands.
	a := make([]*Command, count)

	// Set the index.
	index := 0

	// Go through each command and add it to the array.
	for k, v := range r.cmds {
		if k == v.Name {
			a[index] = v
			index++
		}
	}

	// Read unlock the commands.
	r.cmdLock.RUnlock()

	// Return the commands.
	return a
}

// GetCommandsOrderedByCategory is used to order commands by the categories.
func (r *Router) GetCommandsOrderedByCategory() map[CategoryInterface][]*Command {
	m := map[CategoryInterface][]*Command{}
	cmds := r.GetAllCommands()
	for _, v := range cmds {
		items, ok := m[v.Category]
		if !ok {
			items = []*Command{}
		}
		m[v.Category] = append(items, v)
	}
	return m
}
