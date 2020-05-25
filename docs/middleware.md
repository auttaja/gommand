# Middleware
Middleware allows you to write powerful extensions on a per-command or per-router basis. Middleware is seperate from permission validators to allow the application to tell if the user has permission without re-executing all of the middleware which has been set. Middleware follows the format `func(ctx *Context) error`, with any errors being passed to the error handler. If you wish to get an argument from a middleware function to another function or command during execution, you can use the `MiddlewareParams` map within the context.

Middleware can be added in an array to the command, within categories, or within the global router.
