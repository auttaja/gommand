# Prefix Checkers
Gommand contains the following prefix checkers:

- `gommand.StaticPrefix(<prefix>)`: This will return a function which can be used in the place of the prefix check attribute to specifically look for a static prefix.
- `gommand.MentionPrefix`: This is used to check if the bot is mentioned.
- `gommand.MultiplePrefixCheckers(<prefix checker>...)` - This allows you to combine prefix checkers. In the event that a prefix checker returns false, the read seeker will be rewound back to where it was and the next checker will be called.

In the event that these prefix checkers won't suffice, you can write your own with the function type `func(ctx *gommand.Context, r io.ReadSeeker) bool` where `true` represents if the prefix is used. If the prefix was used, you should also set `ctx.Prefix` to your prefix. Note that the context does not contain the member object in the message or the command yet. If this is nil, it defaults to no prefix.
