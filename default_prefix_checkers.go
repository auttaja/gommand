package gommand

// StaticPrefix is used for simple static prefixes.
func StaticPrefix(Prefix string) func(_ *Context, r *StringIterator) bool {
	l := len(Prefix)
	return func(ctx *Context, r *StringIterator) bool {
		i := 0
		for i != l {
			b, err := r.GetChar()
			if err != nil {
				return false
			}
			if b != Prefix[i] {
				return false
			}
			i++
		}
		ctx.Prefix = Prefix
		return true
	}
}

// MentionPrefix is used to handle a mention which is being used as a prefix.
func MentionPrefix(ctx *Context, r *StringIterator) bool {
	// Get the bot ID.
	BotID := ctx.BotUser.ID.String()

	// The ID to compare.
	CmpID := getMention(r, '@')

	// Is it nil or not the ID?
	if CmpID == nil || *CmpID != BotID {
		return false
	}

	// Remove any whitespace.
	for {
		c, err := r.GetChar()
		if err != nil {
			break
		}
		if c != ' ' {
			r.Pos--
			break
		}
	}

	// Return true and set the prefix.
	ctx.Prefix = "<@" + BotID + "> "
	return true
}

// MultiplePrefixCheckers is used to handle multiple prefix checkers.
func MultiplePrefixCheckers(Handlers ...PrefixCheck) func(ctx *Context, r *StringIterator) bool {
	return func(ctx *Context, r *StringIterator) bool {
		pos := r.Pos
		for _, v := range Handlers {
			if v(ctx, r) {
				return true
			}
			r.Pos = pos
		}
		return false
	}
}
