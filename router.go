package gommand

import (
	"github.com/andersfylling/disgord"
	"github.com/auttaja/fastparse"
	"io"
	"strings"
	"sync"
)

// PrefixCheck is the type for a function to check the prefix. true here means the prefix is there and was read.
// Note that prefix functions cannot see the member object in the message since that is patched in AFTER this is ran.
type PrefixCheck = func(ctx *Context, r io.ReadSeeker) bool

// ErrorHandler is a function which is used to handle errors.
// true here means that the error was handled by this handler. If this is false, the processor will go to the next handler.
type ErrorHandler = func(ctx *Context, err error) bool

// CustomCommandsHandler is the handler which is used for custom commands.
// An error being returned by the custom commands handler will return in that being passed through to the error handler instead.
// true here represents this being a custom command. This means the Router will not go through the errors handler unless an error is set.
type CustomCommandsHandler = func(ctx *Context, cmdname string, r io.ReadSeeker) (bool, error)

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
	MessageCacheHandler  *MessageCacheHandler
	PrefixCheck          PrefixCheck
	ErrorHandlers        []ErrorHandler
	PermissionValidators []PermissionValidator
	Middleware           []Middleware
	Cooldown             Cooldown

	// The number if message pads which will be created in memory to allow for quicker parsing.
	// Please set this to -1 if you do not want any, 0 will default to 100.
	MessagePads int
}

// Router defines the command router which is being used.
// Please call NewRouter to initialise this rather than creating a new struct.
type Router struct {
	PrefixCheck           PrefixCheck           `json:"-"`
	CustomCommandsHandler CustomCommandsHandler `json:"-"`
	cmds                  map[string]CommandInterface
	botUsers              map[uint]*disgord.User
	cmdLock               *sync.RWMutex
	errorHandlers         []ErrorHandler
	permissionValidators  []PermissionValidator
	middleware            []Middleware
	parserManager         *fastparse.ParserManager
	MessageCacheHandler   *MessageCacheHandler
	Cooldown              Cooldown
}

// NewRouter creates a new command Router.
func NewRouter(Config *RouterConfig) *Router {
	if Config.PrefixCheck == nil {
		// No prefix.
		Config.PrefixCheck = func(ctx *Context, r io.ReadSeeker) bool {
			return true
		}
	}

	// Set the Router.
	if Config.MessagePads == 0 {
		Config.MessagePads = 100
	} else if 0 > Config.MessagePads {
		Config.MessagePads = 0
	}
	r := &Router{
		PrefixCheck:          Config.PrefixCheck,
		cmds:                 map[string]CommandInterface{},
		cmdLock:              &sync.RWMutex{},
		errorHandlers:        Config.ErrorHandlers,
		permissionValidators: Config.PermissionValidators,
		middleware:           Config.Middleware,
		MessageCacheHandler:  Config.MessageCacheHandler,
		Cooldown:             Config.Cooldown,
		botUsers:             map[uint]*disgord.User{},
		parserManager:        fastparse.NewParserManager(2000, Config.MessagePads),
	}

	// Set the help command.
	r.SetCommand(defaultHelpCommand())

	// If deleted message handler isn't nil, initialise the storage adapter.
	if r.MessageCacheHandler != nil {
		if r.MessageCacheHandler.MessageCacheStorageAdapter == nil {
			r.MessageCacheHandler.MessageCacheStorageAdapter = &InMemoryMessageCacheStorageAdapter{}
		}
		r.MessageCacheHandler.MessageCacheStorageAdapter.Init()
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
func (r *Router) GetCommand(Name string) CommandInterface {
	r.cmdLock.RLock()
	cmd := r.cmds[strings.ToLower(Name)]
	r.cmdLock.RUnlock()
	return cmd
}

// SetCommand is used to set a command.
func (r *Router) SetCommand(c CommandInterface) {
	c.Init()
	cooldown := c.GetCooldown()
	if cooldown != nil {
		cooldown.Init()
	}
	r.cmdLock.Lock()
	name := strings.ToLower(c.GetName())
	r.cmds[name] = c
	if c.GetAliases() != nil {
		for _, v := range c.GetAliases() {
			v = strings.ToLower(v)
			r.cmds[v] = c
		}
	}
	r.cmdLock.Unlock()
}

// RemoveCommand is used to remove a command from the Router.
func (r *Router) RemoveCommand(c CommandInterface) {
	r.cmdLock.Lock()
	delete(r.cmds, strings.ToLower(c.GetName()))
	if c.GetAliases() != nil {
		for _, v := range c.GetAliases() {
			delete(r.cmds, strings.ToLower(v))
		}
	}
	r.cmdLock.Unlock()
}

// GetAllCommands is used to get all of the commands.
func (r *Router) GetAllCommands() []CommandInterface {
	// Read lock the commands.
	r.cmdLock.RLock()

	// Get the command count.
	count := 0
	for k, v := range r.cmds {
		if k == v.GetName() {
			count++
		}
	}

	// Allocate the commands.
	a := make([]CommandInterface, count)

	// Set the index.
	index := 0

	// Go through each command and add it to the array.
	for k, v := range r.cmds {
		if k == v.GetName() {
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
func (r *Router) GetCommandsOrderedByCategory() map[CategoryInterface][]CommandInterface {
	m := map[CategoryInterface][]CommandInterface{}
	cmds := r.GetAllCommands()
	for _, v := range cmds {
		items, ok := m[v.GetCategory()]
		if !ok {
			items = []CommandInterface{}
		}
		m[v.GetCategory()] = append(items, v)
	}
	return m
}
