# Context
The context is a core part of the gommand functionality. The context contains several crucial bits of information:

- `Message`: The base message which the command is relating to. Unless otherwise specified, the member object will be patched into this message.
- `BotUser`: The `*disgord.User` object which is repersenting the bot. Do **NOT** edit this since it is shared across command calls.
- `Router`: The base router.
- `Session`: The `*disgord.Session` which was used to emit this event.
- `Command`: The actual command which was called.
- `RawArgs`: A string of the raw arguments.
- `Args`: The transformed arguments.
- `Prefix`: Defines the prefix which was used.
- `MidddlewareParams`: The params set by [middleware](./middleware.md).

It also contains several helper functions:

- `Replay() error`: Allows you to replay a command.
- `BotMember() (*disgord.Member, error)`: Get the bot as a member of the guild which the command is being ran in.
- `Channel() (*disgord.Channel, error)`: Get the channel which this is being ran in.
- `Reply(data ...interface{}) (*disgord.Message, error)`: A shorter way to quickly reply to a message.
- `WaitForMessage(CheckFunc func(s disgord.Session, msg *disgord.Message) bool) *disgord.Message`: Waits for a message based on the check function you gave.
- `DisplayEmbedMenu(m *EmbedMenu) error`: Used to display an [embed menu](./embed-menus.md).
