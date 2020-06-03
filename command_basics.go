package gommand

// CommandBasics is the basic command structure minus Init and CommandFunction.
// The objective is that you can inherit this with your structs if you wish to make your own.
type CommandBasics struct {
	Name                 string                `json:"name"`
	Aliases              []string              `json:"aliases"`
	Description          string                `json:"description"`
	Usage                string                `json:"usage"`
	Category             CategoryInterface     `json:"category"`
	PermissionValidators []PermissionValidator `json:"-"`
	ArgTransformers      []ArgTransformer      `json:"-"`
	Middleware           []Middleware          `json:"-"`
	parent               *Command
}

// GetName is used to get the name.
func (obj *CommandBasics) GetName() string {
	if obj.parent == nil {
		return obj.Name
	} else {
		return obj.parent.Name
	}
}

// GetAliases is used to get the aliases.
func (obj *CommandBasics) GetAliases() []string {
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
func (obj *CommandBasics) GetDescription() string {
	if obj.parent == nil {
		return obj.Description
	} else {
		return obj.parent.Description
	}
}

// GetUsage is used to get the usage.
func (obj *CommandBasics) GetUsage() string {
	if obj.parent == nil {
		return obj.Usage
	} else {
		return obj.parent.Usage
	}
}

// GetCategory is used to get the category.
func (obj *CommandBasics) GetCategory() CategoryInterface {
	if obj.parent == nil {
		return obj.Category
	} else {
		return obj.parent.Category
	}
}

// GetPermissionValidators is used to get the permission validators.
func (obj *CommandBasics) GetPermissionValidators() []PermissionValidator {
	if obj.parent == nil {
		return obj.PermissionValidators
	} else {
		return obj.parent.PermissionValidators
	}
}

// GetArgTransformers is used to get the arg transformers.
func (obj *CommandBasics) GetArgTransformers() []ArgTransformer {
	if obj.parent == nil {
		return obj.ArgTransformers
	} else {
		return obj.parent.ArgTransformers
	}
}

// GetMiddleware is used to get the middleware.
func (obj *CommandBasics) GetMiddleware() []Middleware {
	if obj.parent == nil {
		return obj.Middleware
	} else {
		return obj.parent.Middleware
	}
}
