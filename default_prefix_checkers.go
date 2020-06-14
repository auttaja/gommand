package gommand

import "io"

import "strings"

// StaticPrefix is used for simple static prefixes.
func StaticPrefix(Prefix string) func(_ *Context, r io.ReadSeeker) bool {
	bytes := []byte(Prefix)
	l := len(bytes)
	return func(ctx *Context, r io.ReadSeeker) bool {
		i := 0
		ob := make([]byte, 1)
		for i != l {
			_, err := r.Read(ob)
			if err != nil {
				return false
			}
			if ob[0] != bytes[i] {
				return false
			}
			i++
		}
		ctx.Prefix = Prefix
		return true
	}
}

// MentionPrefix is used to handle a mention which is being used as a prefix.
func MentionPrefix(ctx *Context, r io.ReadSeeker) bool {
	// Get the bot ID.
	BotID := ctx.BotUser.ID.String()

	// The ID to compare.
	CmpID := getMention(r, '@', false)

	// Is it nil or not the ID?
	if CmpID == nil || *CmpID != BotID {
		return false
	}

	// Remove any whitespace.
	for {
		ob := make([]byte, 1)
		_, err := r.Read(ob)
		if err != nil {
			break
		}
		if ob[0] != ' ' {
			_, _ = r.Seek(-1, io.SeekCurrent)
			break
		}
	}

	// Return true and set the prefix.
	ctx.Prefix = "<@" + BotID + "> "
	return true
}

// MultiplePrefixCheckers is used to handle multiple prefix checkers.
func MultiplePrefixCheckers(Handlers ...PrefixCheck) func(ctx *Context, r io.ReadSeeker) bool {
	return func(ctx *Context, r io.ReadSeeker) bool {
		sr := r.(*strings.Reader)
		s := sr.Size()
		for _, v := range Handlers {
			if v(ctx, r) {
				return true
			}
			read := s - int64(sr.Len())
			_, _ = r.Seek(read*-1, io.SeekCurrent)
		}
		return false
	}
}
