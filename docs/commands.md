# Commands
Commands are the fundemental part of how gommand works.

When you have a working router, you should be ready to add commands. We do this with the `SetCommand` function. To use this function, we simply pass in a `Command` object or `CommandInterface` which we want to use:
```go
router.SetCommand(&gommand.Command{
    ...
})
```

## `Command`
The command **MUST** have the `Name` (the name of the command) and `Function` (the function that will be called with the [context](/context) and can return an error, any returned errors will be given to the error handlers) attributes set. The other attributes are optional:

- `Aliases`: Any aliases which a command has.
- `Description`: The description which is used in help commands.
- `Usage`: The usage information for a command.
- `PermissionValidators`: An array of [permission validators](/permission-validators) which only applies to this specific command.
- `ArgTransformers`: This is an array of the `gommand.ArgTransformer` type. Each object in this array contains the following attributes:
    - `Function`: The function which is used to transform the argument which must be set. The function simply takes the [context](/context) and argument as a string and returns an interface and error (if the error is nil - this parsed properly). The following transformers are supported by gommand right now:
        - `gommand.StringTransformer`: Transforms the argument into a string.
        - `gommand.IntTransformer`: Transforms the argument into a integer.
        - `gommand.UIntTransformer`: Transforms the argument into a unsigned integer.
        - `gommand.UserTransformer`: Transforms the argument into a user.
        - `gommand.MemberTransformer`: Transforms the argument into a member.
        - `gommand.ChannelTransformer`: Transforms the argument into a channel.
        - `gommand.MessageURLTransformer`: Transforms a message URL into a message.
        - `gommand.BooleanTransformer`: Transforms the argument into a boolean.
        - `gommand.RoleTransformer`: Transforms the argument into a role.
        - `gommand.DurationTransformer`: Transforms the argument into a duration.
        - `gommand.AnyTransformer(...transformer)`: This takes multiple transformers and tries to find one which works.
    - `Optional`: If this is true and the argument does not exist, it will be set to nil. Note that due to what this does, it has to be either at the end of the argument list or followed by other optional arguments (if you don't combine with Remainder).
    - `Remainder`: If this is true, it will just try and parse the raw remainder of the arguments. If the string is blank it will error with not enough arguments unless optional is set. Note that due to what this does, it has to be at the end of the array.
    - `Greedy`: If this is true, the parser will keep trying to parse the users arguments until it hits the end of their message or a parse fails. When this happens, it will go to the next parser in the array. Note that if the first argument fails, this means that it was not set and an error will be put into the error handler unless it was set as optional. The greedy argument will be of the type `[]interface{}` (unless `Optional` is set and it was not specified).
- `Middleware`: An array of [middleware](/middleware) which only applies to this specific command.
- `Category`: Allows you to set a [category](/categories) for your command.
- `CommandAttributes`: A generic interface which you can use for whatever you want.

## `CommandInterface`

What if you want to create commands as structs or you want more flexibility in the process though? We've thought of you, don't worry! By default, gommand uses the `CommandInterface` interface for commands. This means that your command does not have to be of the `Command` type, it can instead just support the following:

- `GetName() string`: Gets the name of the command.
- `GetAliases() []string`: Gets the aliases of the command. This can't be nil.
- `GetDescription() string`: Gets the description of the command if it's set.
- `GetUsage() string`: Gets the usage of the command if it's set.
- `GetCategory() CategoryInterface`: Gets the category of the command if it's set.
- `GetPermissionValidators() []PermissionValidator`: Gets the permission validators of the command. This can't be nil.
- `GetArgTransformers() []ArgTransformer`: Gets the argument transformers of the command.
- `GetMiddleware() []Middleware`: Get the middleware of the command. This cannot be nil.
- `Init()`: Called to initialise the interface.
- `CommandFunction(ctx *Context) error`: The main function for the command.

What if you just want to use a struct for a command? You won't want to write all of that everytime. Therefore, the `CommandBasics` struct was created. This contains a lot of the attributes in the command struct minus the command function/initialisation, allowing you to simply do this:

```go
type ping struct {
	gommand.CommandBasics
}

func (p *ping) Init() {
	p.Name = "ping"
	p.Description = "Responds with pong."
}

func (ping) CommandFunction(ctx *gommand.Context) error {
	_, _ = ctx.Reply("Pong!")
	return nil
}
```
