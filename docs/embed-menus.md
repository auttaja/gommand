# Embed Menus

Embed menus allow your bot to have interactive components within it. This means that you do not need to make a user run a command for every action. Instead, you can simply build a menu to do what you require. To start you will want to make a new embed menu. You can do this with the `NewEmbedMenu` function, setting the arguments to the embed which you want to be the front page and the command context:
```go
menu := gommand.NewEmbedMenu(&disgord.Embed{
    Title:       "Want to see doge?",
    Description: "Click the option below.",
}, ctx)
```

Embed menus have the following functions attached to them:

- `AddBackButton()`: Adds the back button to the embed menu. **Don't use this on the first embed menu, this is meant for child menus.**
- `AddExitButton()`: Adds an exit button to the menu, causing the menu message to delete. This can be used on any menu.
- `AddParentMenu(Menu *EmbedMenu)`: Sets the parent of the menu.
- `Display(ChannelID, MessageID disgord.Snowflake, client disgord.Session) error`: Manually displays the embed. Note that using the `DisplayEmbedMenu` function below should be prefered since it is much easier.
- `NewChildMenu(options *ChildMenuOptions) *EmbedMenu`: Create a child menu with the options specified. The following options can be set in `ChildMenuOptions`:
	- `Embed *disgord.Embed`: Defines the child embed.
	- `Button *MenuButton`: Defines the [menu button](#menu-button).
	- `BeforeAction func()`: Defines the action to run before the menu is displayed. This can be nil.
	- `AfterAction func()`: Defines the action to run after the menu is displayed. This can be nil.

Embed menus also have the `Reactions` attribute attached to them. This contains the `Add` function which can take a `MenuReaction` object. This contains the following:
- `Button *MenuButton`: Defines the [menu button](#menu-button).
- `Function func(ChannelID, MessageID disgord.Snowflake, menu *EmbedMenu, client disgord.Session)`: Defines the function that will be called when the button is clicked.


The command context has defines two key functions for easily displaying the embed: 
- `DisplayEmbedMenu(m *EmbedMenu)`: Displays the embed and returns any errors while creating it.
```go
err := ctx.DisplayEmbedMenu(menu)
```
- `DisplayEmbedMenuWithLifetime(m *EmbedMenu, lifetime *EmbedLifetimeOptions)`: Displays the embed and allows the usage of [embed lifetimes](#lifetime-options), returning any errors while creating it.
```go
err := ctx.DisplayEmbedMenuWithLifetime(menu, &gommand.EmbedLifetimeOptions{
	MaximumLifetime: time.Minute * 2,
})
```

## Menu button
`MenuButton` objects are used to describe reactions that can be clicked to do actions on the menu. This contains the following attributes:
- `Emoji`: The unicode emoji that will be used.
- `Name`: The name of the menu button.
- `Description`: The description of the menu button.

If the button doesn't include either a `Name` or a `Description`, the help field will be omitted entirely.

## Lifetime Options
`EmbedLifetimeOptions` objects are used to control the maximum duration for which an embed menu will be active / exist. This contains the following attributes:
- `MaximumLifetime`: The maximum duration the embed should exist for, regardless of activity - optional.
- `InactiveLifetime`: The maximum duration after the most recent reaction to the menu that the menu should exist for - optional.
After the duration of either of the above has passed, the menu will be deleted.
- `BeforeDelete`: The function called when the menu is scheduled to be deleted, but just before the message itself is deleted.
- `AfterDelete`: The function called after the menu message is deleted, ran regardless of any errors when deleting the message.