# Welcome to Gommand!

Welcome to Gommand! Gommand has a lot of built in functionality to allow for a higher level commands experience. Want to write your first bot with Gommand? [Start here](/writing-your-first-bot):

- **Custom prefix support:** Gommand allows you to easily set a custom prefix. Gommand has various helper functions for this such as static prefix and mention support.

- **Command alias support:** Out ot the box, gommand features support for command aliases. This allows for the ability to easily write commands that have several wordings. This is very useful for bots where you have command names which are changing during migration.

- **Custom commands support:** Gommand allows you to easily set a handler for custom commands if this is something which you require for your bot.

- **Custom argument support:** Out of the box, gommand features many different argument converters for functionality such as integers and users. Should you need it, gommand also allows for the creation of custom converters through a very simple function.

- **Argument formatting:** Gommand supports required, optional, remainder and greedy processors for commands along with many different argument transformers.

- **Middleware support:** Need to run something before a bunch of commands which requires repeating yourself? No problem! Gommand supports middleware on both a global (through the router object) and local (through the command object) scale. Simply add the function and it will be executed. There is also a map within the context called `MiddlewareParams` in which middleware can very easily store data for the commands to run.

- **Permission validators support:** If you need to validate permissions, you can simply use permission validators. If false is returned, a `IncorrectPermissions` will be sent to the error handlers with the string specified.

- **Ease of use:** We pride this library on ease of use. The idea is to build a higher level wrapper around handling commands using disgord, and whilst doing this we add in helper functions. Additionally, we patch the member object into messages, meaning that you do not need to get this, solving potential code repetition.

- **Advanced error handling:** Gommand features an advanced error handling system. The error is passed through all inserted error handlers in the order which the handlers were inserted. If a handler is able to solve an issue, it will simply return true and gommand will stop iterating through the list. If no handlers return true, it will be passed through to the disgord logger which you set.

- **Battle tested:** The Gommand library is heavily unit tested. Feel free to submit a pull request if you feel there is something important which we are not testing, we will accept it. Find a bug? Feel free to file an issue and we will fix it.
