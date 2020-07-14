package gommand

import (
	"io"
	"strings"
)

// ArgTransformer defines a transformer which is to be used on arguments.
type ArgTransformer struct {
	// Greedy defines if the parser should keep going until an argument fails.
	Greedy bool

	// Optional defines if the argument is optional. This can be mixed with greedy or remainder.
	// This has to be either at the end of the argument list or followed by other optional arguments (if you don't combine with Remainder).
	Optional bool

	// Remainder defines if it should just parse the rest of the arguments.
	// Remainders need to be at the end of a command.
	Remainder bool

	// Function is used to transform the argument. The function should error if this is not possible.
	Function func(ctx *Context, Arg string) (interface{}, error)

	// Default is the value the arg will have if it isn't supplied by the user.
	// This will not be used in case of an error from the transformer function, only when no argument is passed by the user.
	// Similarly to optional, this either has to be one of the end arguments (or followed by other default arguments only).
	Default interface{}
}

// CommandInterface is used to define a interface which is used for commands.
type CommandInterface interface {
	// Get the attributes of the command.
	GetName() string
	GetAliases() []string
	GetDescription() string
	GetUsage() string
	GetCooldown() Cooldown
	GetCategory() CategoryInterface
	GetPermissionValidators() []PermissionValidator
	GetArgTransformers() []ArgTransformer
	GetMiddleware() []Middleware

	// Initialisation function and command.
	Init()
	CommandFunction(ctx *Context) error
}

// Command defines a command which can be used within the Router.
type Command struct {
	*CommandBasics       `json:"-"`
	Name                 string                   `json:"name"`
	Aliases              []string                 `json:"aliases"`
	Description          string                   `json:"description"`
	Usage                string                   `json:"usage"`
	Category             CategoryInterface        `json:"category"`
	Cooldown             Cooldown                 `json:"cooldown"`
	CommandAttributes    interface{}              `json:"commandAttributes"`
	PermissionValidators []PermissionValidator    `json:"-"`
	ArgTransformers      []ArgTransformer         `json:"-"`
	Middleware           []Middleware             `json:"-"`
	Function             func(ctx *Context) error `json:"-"`
}

// Init is used to initialise the command.
func (c *Command) Init() {
	c.CommandBasics = &CommandBasics{parent: c}
}

// CommandFunction is used to run the command function.
func (c *Command) CommandFunction(ctx *Context) error {
	return c.Function(ctx)
}

// CommandHasPermission is used to run through the permission validators and check if the user has permission.
// The error will be of the IncorrectPermissions type if they do not have permission.
func CommandHasPermission(ctx *Context, c CommandInterface) error {
	// Run any permission validators on a global scale.
	if ctx.Router.permissionValidators != nil {
		for _, v := range ctx.Router.permissionValidators {
			msg, ok := v(ctx)
			if !ok {
				return &IncorrectPermissions{err: msg}
			}
		}
	}

	// Run any permission validators on a category scale.
	if c.GetCategory() != nil {
		for _, v := range c.GetCategory().GetPermissionValidators() {
			msg, ok := v(ctx)
			if !ok {
				return &IncorrectPermissions{err: msg}
			}
		}
	}

	// Run any permission validators on a local scale.
	if c.GetPermissionValidators() != nil {
		for _, v := range c.GetPermissionValidators() {
			msg, ok := v(ctx)
			if !ok {
				return &IncorrectPermissions{err: msg}
			}
		}
	}

	// Return no errors.
	return nil
}

// HasPermission is a shorthand for running CommandHasPermission on the command.
// This is here to prevent a breaking change.
func (c *Command) HasPermission(ctx *Context) error {
	return CommandHasPermission(ctx, c)
}

// Used to run the command.
func runCommand(ctx *Context, reader io.ReadSeeker, c CommandInterface) (err error) {
	// Handle recovering from exceptions.
	defer func() {
		if r := recover(); r != nil {
			switch v := r.(type) {
			case string:
				ctx.Router.errorHandler(ctx, &PanicError{msg: v})
			case error:
				ctx.Router.errorHandler(ctx, v)
			}
		}
	}()

	// Run any permission validators.
	err = CommandHasPermission(ctx, c)
	if err != nil {
		return
	}

	// Get the category.
	cat := c.GetCategory()

	// Check if the command is on cooldown.
	cmdCooldown := c.GetCooldown()
	if cmdCooldown != nil {
		msg, ok := cmdCooldown.Check(ctx)
		if !ok {
			return &CommandOnCooldown{Message: msg}
		}
	}
	routerCooldown := ctx.Router.Cooldown
	var catCooldown Cooldown
	if cat != nil {
		catCooldown = cat.GetCooldown()
	}
	if catCooldown != nil && cmdCooldown != catCooldown && routerCooldown != catCooldown {
		msg, ok := catCooldown.Check(ctx)
		if !ok {
			return &CommandOnCooldown{Message: msg}
		}
	}
	if routerCooldown != nil && cmdCooldown != routerCooldown {
		msg, ok := routerCooldown.Check(ctx)
		if !ok {
			return &CommandOnCooldown{Message: msg}
		}
	}

	// Run any middleware.
	if ctx.Router.middleware != nil {
		for _, v := range ctx.Router.middleware {
			err = v(ctx)
			if err != nil {
				return
			}
		}
	}
	if cat != nil {
		for _, v := range cat.GetMiddleware() {
			err = v(ctx)
			if err != nil {
				return
			}
		}
	}
	if c.GetMiddleware() != nil {
		for _, v := range c.GetMiddleware() {
			err = v(ctx)
			if err != nil {
				return
			}
		}
	}

	// Transform all arguments if this is possible. If not, error.
	if c.GetArgTransformers() != nil {
		// Slice the arguments.
		ArgCount := 0
		for _, v := range c.GetArgTransformers() {
			ArgCount++
			if v.Remainder {
				break
			}
		}

		// The array containing all transformed arguments.
		Args := make([]interface{}, ArgCount)

		// Defines the argument parser.
		parser := ctx.Router.parserManager.Parser(reader)

		// This is where we transform each argument.
		for i, v := range c.GetArgTransformers() {
			if v.Remainder {
				// Get the remainder.
				remainder, _ := parser.Remainder()
				remainder = strings.Trim(remainder, " ")
				if remainder == "" {
					if v.Default != nil {
						Args[i] = v.Default
						continue
					}
					// Is this an optional argument?
					if !v.Optional {
						parser.Done()
						return &IncorrectPermissions{err: "Remainder cannot be optional."}
					}
				} else {
					x, err := v.Function(ctx, remainder)
					if err != nil {
						parser.Done()
						return err
					}
					Args[i] = x
				}
				break
			} else if v.Greedy {
				// Keep going until there's an error.
				FirstArg := true
				ArgsTransformed := make([]interface{}, 0, 1)
				for {
					Argument := parser.GetNextArg()
					if Argument == nil {
						if FirstArg {
							// This is the first argument.
							// This is important because we are expecting a result if this is not optional.
							if v.Default != nil {
								// The user didn't provide an argument here, but there was a default argument for this converter - use that instead.
								ArgsTransformed = append(ArgsTransformed, v.Default)
								break
							} else if v.Optional {
								// This is optional! We can break the loop here.
								break
							} else {
								// This isn't optional and no default was provided - throw an error.
								err = &InvalidArgCount{err: "Expected an argument for the greedy converter."}
								parser.Done()
								return
							}
						} else {
							break
						}
					} else {
						// Attempt to parse this argument.
						res, err := v.Function(ctx, Argument.Text)
						if err != nil {
							if FirstArg {
								parser.Done()
								return err
							}

							_ = Argument.Rewind()
							break
						}
						ArgsTransformed = append(ArgsTransformed, res)
					}
					FirstArg = false
				}
				if len(ArgsTransformed) != 0 {
					Args[i] = ArgsTransformed
				}
			} else {
				// Try and get one argument.
				Argument := parser.GetNextArg()
				if Argument != nil {
					x, err := v.Function(ctx, Argument.Text)
					if err != nil {
						parser.Done()
						return err
					}
					Args[i] = x
					continue
				}
				if v.Default != nil {
					Args[i] = v.Default
					continue
				}
				if v.Optional {
					break
				} else {
					parser.Done()
					return &InvalidArgCount{err: "A required argument is missing."}
				}
			}
		}

		// Mark the parser as done.
		parser.Done()

		// Set the arguments.
		ctx.Args = Args
	}

	// Run the command and return.
	err = c.CommandFunction(ctx)
	return
}
