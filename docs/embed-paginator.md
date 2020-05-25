# Embed Paginator

Gommand contains embed menus, but what if you simply want to create a page system by passing through an array of embeds? Gommand supports the ability to do this with a simple high level API. To do this, we can use the `EmbedsPaginator` function within Gommand. The paginator takes the command [context](/context) as the first argument, and the array of embeds as the second argument. In the event that this is unable to send, the error result will be set.
```go
err := gommand.EmbedsPaginator(ctx, []*disgord.Embed{
    {
        Title: "Hi",
    },
    {
        Title: "World",
    },
})
```
