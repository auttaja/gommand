package gommand

// CommandNotFound is the error which is thrown when a command is not found.
type CommandNotFound struct {
	err string
}

// Error is used to give the error description.
func (c *CommandNotFound) Error() string {
	return c.err
}

// CommandBlank is the error which is thrown when the command is blank.
type CommandBlank struct {
	err string
}

// Error is used to give the error description.
func (c *CommandBlank) Error() string {
	return c.err
}

// IncorrectPermissions is the error which is thrown when the user does not have enough permissions.
type IncorrectPermissions struct {
	err string
}

// Error is used to give the error description.
func (c *IncorrectPermissions) Error() string {
	return c.err
}

// InvalidArgCount is the error when the arg count is not correct.
type InvalidArgCount struct {
	err string
}

// Error is used to give the error description.
func (c *InvalidArgCount) Error() string {
	return c.err
}

// InvalidTransformation is the error argument parsers should use when they can't transform.
type InvalidTransformation struct {
	Description string
}

// Error is used to give the error description.
func (c *InvalidTransformation) Error() string {
	return c.Description
}

// PanicError is used when a string is returned from a panic. If it isn't a string, the error will just be pushed into the handler.
type PanicError struct {
	msg string
}

// Error is used to give the error description.
func (c *PanicError) Error() string {
	return c.msg
}
