package gommand

// StaticPrefix is used for simple static prefixes.
func StaticPrefix(Prefix string) func(_ *Context, r *StringIterator) bool {
	PrefixBytes := []byte(Prefix)
	l := len(PrefixBytes)
	return func(_ *Context, r *StringIterator) bool {
		i := 0
		for i != l {
			b, err := r.GetChar()
			if err != nil {
				return false
			}
			if b != PrefixBytes[i] {
				return false
			}
			i++
		}
		return true
	}
}
