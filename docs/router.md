# Router
Creating the router is very simple to do. You can simply create a router object in the part of your project where all of your commands are by calling the `NewRouter` function with a `RouterConfig` object. The configuration object can contain the following attributes:

- `PrefixCheck`: Defines the [prefix checker](./prefix-checkers.md) which will be used.
- `ErrorHandlers`: An array of functions as described in [writing your first bot](./writing-your-first-bot.md) which will run one after another. This can be nil and you can also add one with the `AddErrorHandler` function attached to the router.
- `PermissionValidators`: This is any [permission validators](./permission-validators.md) which you wish to add on a global router scale. This can be nil.
- `Middleware`: This is any [middleware](./middleware.md) which you wish to add on a global router scale. This can be nil.
- `MessageCacheHandler`: See the [deleted message handler](./handling-deleted-messages.md) documentation below.
- `Cooldown`: The cooldown interface for this router. You should keep this as nil if you don't want a router wide cooldown.

From here, we can use the functions attached to the router:

- `AddErrorHandler(Handler ErrorHandler)`: Used to add a error handler as described in [writing your first bot](./writing-your-first-bot.md).
- `CommandProcessor(s disgord.Session, msg *disgord.Message, prefix bool)`: Used to process a command. You will probably never need to use this.
- `GetAllCommands() []CommandInterface`: Get all commands.
- `GetCommand(Name string) CommandInterface`: Get a command by its name.
- `GetCommandsOrderedByCategory() map[CategoryInterface][]CommandInterface`: Get all commands ordered by their category.
- `Hook(s disgord.Session)`: Used to hook to a disgord session.
- `RemoveCommand(c CommandInterface)`: Used to remove a [command](./command.md).
- `SetCommand(c CommandInterface)`: Used to set the [command](./command.md).
