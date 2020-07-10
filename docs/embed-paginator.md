# Embed Paginator

Gommand contains embed menus, but what if you simply want to create a page system by passing through an array of embeds? Gommand supports the ability to do this with a simple high level API. To do this, we can use the `EmbedsPaginator` function within Gommand. The paginator takes the command [context](./context.md) as the first argument, the array of embeds as the second argument, the initial page number as the third argument (used to flick through pages in the event that the buttons cannot be created) and the text content that will be sent with the embed in the event that the buttons cannot be created as the fourth argument. In the event that this is unable to send, the error result will be set.
```go
err := gommand.EmbedsPaginator(ctx, []*disgord.Embed{
    {
        Title: "Hi",
    },
    {
        Title: "World",
    },
}, 1, "use command <page number> to flick through pages")
```
