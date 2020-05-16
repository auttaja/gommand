package gommand

// Get the mention if it exists.
func getMention(r *StringIterator, char uint8, role bool) *string {
	// Defines the parsing stage.
	stage := uint8(0)

	// Defines the ID we will compare to.
	CmpID := ""

	// Defines if this was a mention.
	mention := false

	// Loop through chars until we are sure it's a mention.
	start := true
	for {
		c, err := r.GetChar()
		if err != nil {
			// Check the stage. If it can't be a ID, return nil.
			if stage == 3 {
				return &CmpID
			}
			return nil
		}
		if c == ' ' {
			if !start {
				// This isn't the start. Is it a mention?
				if mention {
					// Isn't a prefix.
					return nil
				}

				// Return the ID.
				return &CmpID
			}
		} else {
			// Set start to false.
			start = false

			// Check the stage.
			if stage == 0 {
				// We expect a '<' char here.
				if c == '<' {
					// This is ok! move to stage 1 (type symbol).
					stage = 1
					mention = true
				} else if c > 46 && 58 > c {
					// We should move to stage 3.
					CmpID += string(c)
					stage = 3
				} else {
					// This is invalid for this stage.
					return nil
				}
			} else if stage == 1 {
				if c == char {
					// This is increasingly looking like a mention. Move to stage 2 (possible number/explanation mark).
					stage = 2
				} else {
					// This isn't a mention.
					return nil
				}
			} else if stage == 2 {
				x := uint8('!')
				if role {
					x = '&'
				}
				if c == x {
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
						return nil
					}
				}
			} else if stage == 3 {
				if c == '>' {
					// This is the end. We should return here.
					return &CmpID
				} else if c > 46 && 58 > c {
					// Append to the ID.
					CmpID += string(c)
				} else {
					// Return nil if this is a mention. If it is, rewind one and return the ID.
					if mention {
						return nil
					}
					r.Rewind(1)
					return &CmpID
				}
			}
		}
	}
}
