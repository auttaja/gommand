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


The command context has the function `DisplayEmbedMenu` which can be used to easily display the embed. The error defines any issues when it comes to showing the embed:
```go
err := ctx.DisplayEmbedMenu(menu)
```

## Menu button
`MenuButton` objects are used to describe reactions that can be clicked to do actions on the menu. This contains the following attributes:
- `Emoji`: The unicode emoji that will be used.
- `Name`: The name of the menu button.
- `Description`: The description of the menu button.
