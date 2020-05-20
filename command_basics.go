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
}

// GetName is used to get the name.
func (c *CommandBasics) GetName() string {
	return c.Name
}

// GetAliases is used to get the aliases.
func (c *CommandBasics) GetAliases() []string {
	if c.Aliases == nil {
		return []string{}
	}
	return c.Aliases
}

// GetDescription is used to get the description.
func (c *CommandBasics) GetDescription() string {
	return c.Description
}

// GetUsage is used to get the usage.
func (c *CommandBasics) GetUsage() string {
	return c.Usage
}

// GetCategory is used to get the category.
func (c *CommandBasics) GetCategory() CategoryInterface {
	return c.Category
}

// GetPermissionValidators is used to get the permission validators.
func (c *CommandBasics) GetPermissionValidators() []PermissionValidator {
	return c.PermissionValidators
}

// GetArgTransformers is used to get the arg transformers.
func (c *CommandBasics) GetArgTransformers() []ArgTransformer {
	return c.ArgTransformers
}

// GetMiddleware is used to get the middleware.
func (c *CommandBasics) GetMiddleware() []Middleware {
	return c.Middleware
}
