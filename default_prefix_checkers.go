package gommand

// StaticPrefix is used for simple static prefixes.
func StaticPrefix(Prefix string) func(_ *Context, r *StringIterator) bool {
	l := len(Prefix)
	return func(_ *Context, r *StringIterator) bool {
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
		return true
	}
}

// MentionPrefix is used to handle a mention which is being used as a prefix.
func MentionPrefix(ctx *Context, r *StringIterator) bool {
	// Get the bot ID.
	BotID := ctx.BotUser.ID.String()

	// Set the stage of parsing this is on.
	stage := uint8(0)

	// Log the ID to compare.
	CmpID := ""

	// Loop through chars until we are sure it's a mention.
	start := true
	for {
		c, err := r.GetChar()
		if err != nil {
			// Ok this string is too small to be the prefix.
			return false
		}
		if c == ' ' {
			if !start {
				// This isn't the start. Isn't a prefix.
				return false
			}
		} else {
			// Set start to false.
			start = false

			// Check the stage.
			if stage == 0 {
				// We expect a '<' char here.
				if c == '<' {
					// This is ok! move to stage 1 (at symbol).
					stage = 1
				} else {
					// This is invalid for this stage.
					return false
				}
			} else if stage == 1 {
				if c == '@' {
					// This is increasingly looking like a mention. Move to stage 2 (possible number/explanation mark).
					stage = 2
				} else {
					// This isn't a mention.
					return false
				}
			} else if stage == 2 {
				if c == '!' {
					// Ok, we should be ok to move to stage 3 without any ID logging.
					stage = 3
				} else {
					// Is this within 0-9?
					if c > 46 && 58 > c {
						// It is. Add to the ID and make it stage 3.
						stage = 3
						CmpID += string(c)
					} else {
						// Not a mention.
						return false
					}
				}
			} else if stage == 3 {
				if c == '>' {
					// This is the end. We should compare it to the ID.
					if CmpID == BotID {
						// Break here.
						break
					} else {
						// Oof.
						return false
					}
				} else if c > 46 && 58 > c {
					// Append to the ID.
					CmpID += string(c)
				} else {
					// Return false.
					return false
				}
			}
		}
	}

	// Remove any whitespace.
	for {
		c, err := r.GetChar()
		if err != nil {
			break
		}
		if c != ' ' {
			r.Pos -= 1
			break
		}
	}

	// Return true.
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
