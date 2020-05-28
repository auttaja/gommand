# Categories
In Go, categories use the `gommand.CategoryInterface` interface to ensure that they can be modular. The interface has the following functions which must be set:

- `GetName() string`: Gets the name of the category.
- `GetDescription() string`: Gets the description of the category.
- `GetPermissionValidators() []gommand.PermissionValidator`: Gets the array of permission validators. This array cannot be nil.
- `GetMiddleware() []gommand.Middleware`: Gets the array of permission validators. This array cannot be nil.

For ease of use, gommand has the `Category` struct that implements all of these for you. The following attributes can be set in this:

- `Name`: The name of the category.
- `Description`: The description of the category.
- `PermissionValidators`: An array of [permission validators](./permission-validators.md) which will be used on each item in the category. This can be nil.
- `Middleware`: An array of [middleware](./middleware.md) which will be used on each item in the category. This can be nil.

The default help command will automatically take advantage of categories when it is displaying commands. Note that you might want to change the category of the default help command. This is simple to do:
```go
router.GetCommand("help").(*gommand.Command).Category = <category>
```

Note that to allow for the easy categorisation of commands and prevent repetition, a pointer should be created somewhere in your codebase (using `var` or before your commands) to a category which multiple commands use and they should all just pass through the same pointer.
