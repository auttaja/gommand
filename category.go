package gommand

// CategoryInterface is the interface which is used for categories.
// This can be used to write your own category handler if you wish.
type CategoryInterface interface {
	GetName() string
	GetDescription() string
	GetPermissionValidators() []func(ctx *Context) (string, bool) // If I change this to the PermissionValidator type, it triggers #25838 in Go.
	GetMiddleware() []Middleware
	GetCooldown() Cooldown
}

// Category is the generic category struct which uses the category interface.
type Category struct {
	Name                 string                `json:"name"`
	Description          string                `json:"description"`
	Cooldown Cooldown `json:"cooldown"`
	PermissionValidators []PermissionValidator `json:"-"`
	Middleware           []Middleware          `json:"-"`
}

// GetName is used to get the name of the category.
func (c *Category) GetName() string {
	return c.Name
}

// GetDescription is used to get the description of the category.
func (c *Category) GetDescription() string {
	return c.Description
}

// GetPermissionValidators is used to get the permission validators of the category.
func (c *Category) GetPermissionValidators() []PermissionValidator {
	if c.PermissionValidators == nil {
		return []PermissionValidator{}
	}
	return c.PermissionValidators
}

// GetMiddleware is used to get the middleware of the category.
func (c *Category) GetMiddleware() []Middleware {
	if c.Middleware == nil {
		return []Middleware{}
	}
	return c.Middleware
}

// GetCooldown is used to get the cooldown from the category.
func (c *Category) GetCooldown() Cooldown {
	return c.Cooldown
}
