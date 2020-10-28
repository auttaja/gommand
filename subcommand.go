package gommand

import (
	"sort"
	"strings"
)

// CommandGroup is used to have a group of commands which will be executed as sub-commands.
type CommandGroup struct {
	// Name is used to define the category name.
	Name string `json:"name"`

	// Aliases is used to define aliases for this command group.
	Aliases []string `json:"aliases"`

	// Description is used to define a group description.
	Description string `json:"description"`

	// Category is used to define a group category.
	Category CategoryInterface `json:"category"`

	// Cooldown is used to define a group cooldown.
	Cooldown Cooldown `json:"cooldown"`

	// PermissionValidators defines the permission validators for this group.
	PermissionValidators []PermissionValidator `json:"-"`

	// Middleware is used to define the sub-command middleware. Note that this applies to all items in the group.
	Middleware []Middleware `json:"-"`

	// NoCommandSpecified is the command to call when no command is specified.
	NoCommandSpecified CommandInterface

	// Defines the sub-commands.
	subcommands map[string]CommandInterface
}

// GetName is used to get the name.
func (g *CommandGroup) GetName() string {
	return g.Name
}

// GetAliases is used to get the aliases.
func (g *CommandGroup) GetAliases() []string {
	if g.Aliases == nil {
		return []string{}
	}
	return g.Aliases
}

// GetDescription is used to get the description.
func (g *CommandGroup) GetDescription() string {
	return g.Description
}

// GetUsage is used to get the usage.
func (g *CommandGroup) GetUsage() string {
	x := "<"
	keys := make([]string, len(g.subcommands))
	i := 0
	for k := range g.subcommands {
		keys[i] = k
		i++
	}
	sort.Strings(keys)
	for _, cmdname := range keys {
		x += cmdname + "/"
	}
	x = x[:len(x)-1]
	if len(x) == 0 {
		if g.NoCommandSpecified != nil {
			return g.NoCommandSpecified.GetUsage()
		}
		return ""
	}
	return x + ">"
}

// GetCategory is used to get the category.
func (g *CommandGroup) GetCategory() CategoryInterface {
	return g.Category
}

// GetPermissionValidators is used to get the permission validators.
func (g *CommandGroup) GetPermissionValidators() []PermissionValidator {
	return g.PermissionValidators
}

// GetArgTransformers is used to get the arg transformers.
func (g *CommandGroup) GetArgTransformers() []ArgTransformer {
	return []ArgTransformer{
		{
			Function: StringTransformer,
			Optional: true,
		},
		{
			Remainder: true,
			Optional:  true,
			Function:  StringTransformer,
		},
	}
}

// GetCooldown is used to get the cooldown.
func (g *CommandGroup) GetCooldown() Cooldown {
	return g.Cooldown
}

// GetMiddleware is used to get the middleware.
func (g *CommandGroup) GetMiddleware() []Middleware {
	return g.Middleware
}

// CommandFunction is the command function which will be called.
func (g *CommandGroup) CommandFunction(ctx *Context) error {
	cmdname, ok := ctx.Args[0].(string)
	if !ok {
		// Handle the no command specified event if it is set.
		if g.NoCommandSpecified != nil {
			return runCommand(ctx, strings.NewReader(""), g.NoCommandSpecified)
		}

		// Send an error to the router.
		return &CommandBlank{err: "This group expects a command but none was given."}
	}
	args, _ := ctx.Args[1].(string)
	subcommand, ok := g.subcommands[strings.ToLower(cmdname)]
	if ok {
		// Return this command handler.
		return runCommand(ctx, strings.NewReader(args), subcommand)
	}
	return &CommandNotFound{err: "The command specified for the group was not found."}
}

// Init is used to initialise the commands group.
func (g *CommandGroup) Init() {
	if g.subcommands == nil {
		g.subcommands = map[string]CommandInterface{}
	}
	if g.NoCommandSpecified != nil {
		g.NoCommandSpecified.Init()
	}
	uniques := make([]CommandInterface, 0, len(g.subcommands))
	for _, v := range g.subcommands {
		isUnique := true
		for _, x := range uniques {
			if v == x {
				isUnique = false
				break
			}
		}
		if isUnique {
			uniques = append(uniques, v)
			v.Init()
		}
	}
}

// AddCommand is used to add a command to the group.
func (g *CommandGroup) AddCommand(cmd CommandInterface) {
	if g.subcommands == nil {
		g.subcommands = map[string]CommandInterface{}
	}
	cmdname := strings.ToLower(cmd.GetName())
	g.subcommands[cmdname] = cmd
	for _, alias := range cmd.GetAliases() {
		g.subcommands[strings.ToLower(alias)] = cmd
	}
}
