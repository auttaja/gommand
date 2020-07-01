# Permission validators
Permission validators allow for a quick method to check if the user has permission to run a command.

## Built-in Permission Validators

Gommand contains built-in permission validators for all [Discord permisions](https://discord.com/developers/docs/topics/permissions#permissions-bitwise-permission-flags).

To use them, simply call the permission ID you want with a bitwise OR of all of the validation types which you require. The following validation types are supported:
- `CheckMembersUserPermissions`: Used to check the members user permissions.
- `CheckMembersChannelPermissions`: Used to check the members channel permissions.
- `CheckBotUserPermissions`: Used to check the bots user permissions.
- `CheckBotChannelPermissions`: Used to check the bots channel permissions.

For example, if you wanted to check if a user was administrator, you would use the permission validator `gommand.ADMINISTRATOR(gommand.CheckMembersUserPermissions)`. If you also wanted to check if the bot was adminstrator, the validator would be `gommand.ADMINISTRATOR(gommand.CheckMembersUserPermissions | gommand.CheckBotUserPermissions)`.

## DIY Permission Validators

If you wish to write your own permission validators, they follow the format `func(ctx *Context) (string, bool)`. If the boolean is true, the user does have permission. If not, the string is used to construct a `IncorrectPermissions` error.
