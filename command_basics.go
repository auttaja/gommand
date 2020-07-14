package gommand

type commandBasics struct {
	Name                 string                `json:"name"`
	Aliases              []string              `json:"aliases"`
	Description          string                `json:"description"`
	Usage                string                `json:"usage"`
	Category             CategoryInterface     `json:"category"`
	Cooldown             Cooldown              `json:"cooldown"`
	PermissionValidators []PermissionValidator `json:"-"`
	ArgTransformers      []ArgTransformer      `json:"-"`
	Middleware           []Middleware          `json:"-"`
	parent               *Command
}

// CommandBasics is the basic command structure minus Init and CommandFunction.
// The objective is that you can inherit this with your structs if you wish to make your own.
type CommandBasics = commandBasics

// GetName is used to get the name.
func (obj *commandBasics) GetName() string {
	if obj.parent == nil {
		return obj.Name
	} else {
		return obj.parent.Name
	}
}

// GetAliases is used to get the aliases.
func (obj *commandBasics) GetAliases() []string {
	var Aliases []string
	if obj.parent == nil {
		Aliases = obj.Aliases
	} else {
		Aliases = obj.parent.Aliases
	}
	if Aliases == nil {
		return []string{}
	}
	return Aliases
}

// GetDescription is used to get the description.
func (obj *commandBasics) GetDescription() string {
	if obj.parent == nil {
		return obj.Description
	} else {
		return obj.parent.Description
	}
}

// GetUsage is used to get the usage.
func (obj *commandBasics) GetUsage() string {
	if obj.parent == nil {
		return obj.Usage
	} else {
		return obj.parent.Usage
	}
}

// GetCategory is used to get the category.
func (obj *commandBasics) GetCategory() CategoryInterface {
	if obj.parent == nil {
		return obj.Category
	} else {
		return obj.parent.Category
	}
}

// GetPermissionValidators is used to get the permission validators.
func (obj *commandBasics) GetPermissionValidators() []PermissionValidator {
	if obj.parent == nil {
		return obj.PermissionValidators
	} else {
		return obj.parent.PermissionValidators
	}
}

// GetArgTransformers is used to get the arg transformers.
func (obj *commandBasics) GetArgTransformers() []ArgTransformer {
	if obj.parent == nil {
		return obj.ArgTransformers
	} else {
		return obj.parent.ArgTransformers
	}
}

// GetCooldown is used to get the cooldown.
func (obj *commandBasics) GetCooldown() Cooldown {
	if obj.parent == nil {
		return obj.Cooldown
	} else {
		return obj.parent.Cooldown
	}
}

// GetMiddleware is used to get the middleware.
func (obj *commandBasics) GetMiddleware() []Middleware {
	if obj.parent == nil {
		return obj.Middleware
	} else {
		return obj.parent.Middleware
	}
}
