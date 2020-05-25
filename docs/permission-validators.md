# Permission validators
Permission validators allow for a quick method to check if the user has permission to run a command.

Gommand contains built-in permission validators for all [Discord permisions](https://discord.com/developers/docs/topics/permissions#permissions-bitwise-permission-flags). To use them, simply use the permission ID as the permission validator function. For example, if you want to check if someone is an administrator, simply add `gommand.ADMINISTRATOR` to your `PermissionValidators` array on the command, category or router.

If you wish to write your own permission validators, they follow the format `func(ctx *Context) (string, bool)`. If the boolean is true, the user does have permission. If not, the string is used to construct a `IncorrectPermissions` error.
