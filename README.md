# gommand

Package gommand provides an easy to use and high performance commands router and processor for Go(lang) using disgord as the base Discord framework. It features the following functionality:
- **Custom prefix support:** Gommand allows you to easily set a custom prefix. Gommand has various helper functions for this such as static prefix and mention support.
- **Command alias support:** Out ot the box, gommand features support for command aliases. This allows for the ability to easily write commands that have several wordings. This is very useful for bots where you have command names which are changing during migration.
- **Custom commands support:** Gommand allows you to easily set a handler for custom commands if this is something which you require for your bot.
- **Custom argument support:** Out of the box, gommand features many different argument converters for functionality such as integers and users. Should you need it, gommand also allows for the creation of custom converters through a very simple function.
- **Argument formatting:** Gommand supports required, optional, remainder and greedy processors for commands.
- **Middleware support:** Need to run something before a bunch of commands which requires repeating yourself? No problem! Gommand supports middleware on both a global (through the router object) and local (through the command object) scale. Simply add the function and it will be executed. There is also a map within the context called `MiddlewareParams` in which middleware can very easily store data for the commands to run.
- **Permission validators support:** If you need to validate permissions, you can simply use permission validators. If false is returned, a `IncorrectPermissions` will be sent to the error handlers with the string specified.
- **Ease of use:** We pride this library on ease of use. The idea is to build a higher level wrapper around handling commands using disgord, and whilst doing this we add in helper functions. Additionally, we patch the member object into messages, meaning that you do not need to get this, solving potential code repetition.
- **Advanced error handling:** Gommand features an advanced error handling system. The error is passed through all inserted error handlers in the order which the handlers were inserted. If a handler is able to solve an issue, it will simply return true and gommand will stop iterating through the list. If no handlers return true, it will be passed through to the disgord logger which you set.
- **Battle tested:** The Gommand library is heavily unit tested. Feel free to submit a pull request if you feel there is something important which we are not testing, we will accept it. Find a bug? Feel free to file an issue and we will fix it.

An example of the bot is accessable in the `example` folder.

## Contributing
Do you have something which you wish to contribute? Feel free to make a pull request. The following simple criteria applies:

- **Is it useful to everyone?:** If this is a domain specific edit, you probably want to keep this as middleware since it will not be accepted into the main project.
- **Does it slow down parsing by >1ms?:** This will likely be denied. We want to keep parsing as high performance as possible.
- **Have you ran `gofmt -w .` and `golint .`?:** We like to stick to using Go standards within this project, therefore if this is not done you may be asked to do this for it to be acccepted.

Have you experienced a bug? If so, please make an issue! We take bugs seriously and this will be a large priority for us.

## Creating your router
Creating the router is very simple to do. You can simply create a router object in the part of your project where all of your commands are by calling the `NewRouter` function with a `RouterConfig` object. The configuration object can contain the following attributes:
- `PrefixCheck`: This is used to set the checker which will be used for prefixes. Gommand contains the following prefix checks which you can use:
    - `gommand.StaticPrefix(<prefix>)`: This will return a function which can be used in the place of the prefix check attribute to specifically look for a static prefix.
    - `gommand.MentionPrefix`: This is used to check if the bot is mentioned.
    - `gommand.MultiplePrefixCheckers(<prefix checker>...)` - This allows you to combine prefix checkers. In the event that a prefix checker returns false, the string iterator will be rewinded back to where it was and the next checker will be called.
    
    In the event that these prefix checkers won't suffice, you can write your own with the function type `func(ctx *gommand.Context, r *gommand.StringIterator) bool`. Note that the context does not contain the member object in the message or the command yet. See [using the string iterator](#using-the-string-iterator) below on how to use the string iterator. If this is nil, it defaults to no prefix.
- `ErrorHandlers`: An array of functions of the [ErrorHandler](#error-handling) type which will run one after another. This can be nil and you can also add one with the `AddErrorHandler` function attached to the router.
- `PermissionValidators`: This is any [permission validators](#permission-validators) which you wish to add on a global router scale. This can be nil.
- `Middleware`: This is any [middleware](#middleware) which you wish to add on a global router scale. This can be nil.

The router is then trivial to make:
```go
var router = gommand.NewRouter(&gommand.RouterConfig{
    ...
})
```

## Adding a command
When you have a working router, you should be ready to add commands. We do this with the `SetCommand` function. To use this function, we simply pass in a `Command` object:
```go
router.SetCommand(&gommand.Command{
    ...
})
```
The command **MUST** have the `Name` (the name of the command) and `Function` (the function that will be called with the [context](#context) and can return an error, any returned errors will be given to the error handlers) attributes set. The other attributes are optional:
- `Aliases`: Any aliases which a command has.
- `Description`: The description which is used in help commands.
- `Usage`: The usage information for a command.
- `PermissionValidators`: An array of [permission validators](#permission-validators) which only applies to this specific command.
- `ArgTransformers`: This is an array of the `gommand.ArgTransformer` type. Each object in this array contains the following attributes:
    - `Function`: The function which is used to transform the argument which must be set. The function simply takes the [context](#context) and argument as a string and returns an interface and error (if the error is nil - this parsed properly). The following transformers are supported by gommand right now:
        - `gommand.StringTransformer`: Transforms the argument into a string.
        - `gommand.IntTransformer`: Transforms the argument into a integer.
        - `gommand.UIntTransformer`: Transforms the argument into a unsigned integer.
        - `gommand.UserTransformer`: Transforms the argument into a user.
        - `gommand.MemberTransformer`: Transforms the argument into a member.
        - `gommand.ChannelTransformer`: Transforms the argument into a channel.
    - `Optional`: If this is true and the argument does not exist, it will be sert to nil. Note that due to what this does, it has to be at the end of the array.
    - `Remainder`: If this is true, it will just try and parse the raw remainder of the arguments. If the string is blank it will error with not enough arguments unless optional is set. Note that due to what this does, it has to be at the end of the array.
    - `Greedy`: If this is true, the parser will keep trying to parse the users arguments until it hits the end of their message or a parse fails. When this happens, it will go to the next parser in the array. Note that if the first argument fails, this means that it was not set and an error will be put into the error handler unless it was set as optional. The greedy argument will be of the type `[]interface{}` (unless `Optional` is set and it was not specified).
- `Middleware`: An array of [middleware](#middleware) which only applies to this specific command.

## Context
The context is a core part of the gommand functionality. The context contains several crucial bits of information:
- `Message`: The base message which the command is relating to. Unless otherwise specified, the member object will be patched into this message.
- `BotUser`: The `*disgord.User` object which is repersenting the bot. Do **NOT** edit this since it is shared across command calls.
- `Router`: The base router.
- `Session`: The `*disgord.Session` which was used to emit this event.
- `Command`: The actual command which was called.
- `RawArgs`: A string of the raw arguments.
- `Args`: The transformed arguments.
- `MidddlewareParams`: The params set by [middleware](#middleware).

It also contains several helper functions:
- `Replay() error`: Allows you to replay a command.
- `BotMember() (*disgord.Member, error)`: Get the bot as a member of the guild which the command is being ran in.
- `Channel() (*disgord.Channel, error)`: Get the channel which this is being ran in.
- `Reply(data ...interface{}) (*disgord.Message, error)`: A shorter way to quickly reply to a message.
- `WaitForMessage(CheckFunc func(s disgord.Session, msg *disgord.Message) bool) *disgord.Message`: Waits for a message based on the check function you gave.

## Hooking the router to your disgord session
In the initialisation of your disgord session, you will want to hook the gommand handler with the `Hook` function:
```go
// Your client config can be how you please.
client := disgord.New(disgord.Config{
    BotToken: os.Getenv("TOKEN"),
    Logger:   disgord.DefaultLogger(false),
})

// Hook the router.
router.Hook(client)

// ANY OTHER INITIALISATION OF DISGORD EVENTS HERE

// Connect to Discord.
err := client.StayConnectedUntilInterrupted(context.Background())
if err != nil {
    panic(err)
}
```

## Error Handling
In gommand, every negative action is treated as an error. It is then your job to handle these errors. If the error is not handled within the router, it is then just simply passed off to the logger which was configured in disgord. Error handlers take the context and the error. From this they return a boolean. If the boolean is true, it means the error was caught by the handler. If not it simply goes to the next handler in the array. Gommand has several errors which can pass through of its own:
- `*gommand.CommandNotFound`: The command was not found within the router.
- `*gommand.CommandBlank`: The command name was blank.
- `*gommand.IncorrectPermissions`: The permissions this user has are incorrect for the command.
- `*gommand.InvalidArgCount`: The argument count is not correct.
- `*gommand.InvalidTransformation`: Passed through from a transformer when it cannot transform properly.

Using this, we can build a simple error handler that ignores command not found events and logs errors to the chat for the others, although you may wish to implement this differently:
```go
func basicErrorHandler(ctx *gommand.Context, err error) bool {
    // Check all the different types of errors.
    switch err.(type) {
    case *gommand.CommandNotFound, *gommand.CommandBlank:
        // We will ignore.
        return true
    case *gommand.InvalidTransformation:
        _, _ = ctx.Reply("Invalid argument:", err.Error())
        return true
    case *gommand.IncorrectPermissions:
        _, _ = ctx.Reply("Invalid permissions:", err.Error())
        return true
    case *gommand.InvalidArgCount:
        _, _ = ctx.Reply("Invalid argument count.")
        return true
    }

    // This was not handled here.
    return false
}
```
This can then be added to the `ErrorHandlers` array or passed to `AddErrorHandler`. Note that they execute in the order they were added.

## Permission Validators
TODO

## Middleware
TODO

## Using The String Iterator
If you are handling parts of the parsing which are very early in the process as is the case with prefixes and custom commands,0 and you are writing your own code to implement them, you will need to handle the `gommand.StringIterator` type. The objective of this is to try and prevent multiple iterations of the string, which can be computationally expensive, where this is possible. The iterator implements the following:
- `GetRemainder(FillIterator bool) (string, error)`: This will get the remainder of the iterator. If it's already at the end, the error will be set. `FillIterator` defines if it should fill the iterator when it is done or if it should leave it where it is.
- `GetChar() (uint8, error)`: Used to get a character from the iterator. If it's already at the end, the error will be set.
- `Rewind(N uint)`: Used to rewind by N number of chars. Useful if you only iterated a few times to check something.
