package gommand

import "io"

// Get the mention if it exists.
func getMention(r io.ReadSeeker, char uint8, role bool) *string {
	// Defines the parsing stage.
	stage := uint8(0)

	// Defines the ID we will compare to.
	CmpID := ""

	// Defines if this was a mention.
	mention := false

	// Loop through chars until we are sure it's a mention.
	start := true
	for {
		ob := make([]byte, 1)
		_, err := r.Read(ob)
		if err != nil {
			// Check the stage. If it can't be a ID, return nil.
			if stage == 3 {
				return &CmpID
			}
			return nil
		}
		if ob[0] == ' ' {
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
				if ob[0] == '<' {
					// This is ok! move to stage 1 (type symbol).
					stage = 1
					mention = true
				} else if ob[0] > 46 && 58 > ob[0] {
					// We should move to stage 3.
					CmpID += string(ob[0])
					stage = 3
				} else {
					// This is invalid for this stage.
					return nil
				}
			} else if stage == 1 {
				if ob[0] == char {
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
				if ob[0] == x {
					// Ok, we should be ok to move to stage 3 without any ID logging.
					stage = 3
				} else {
					// Is this within 0-9?
					if ob[0] > 46 && 58 > ob[0] {
						// It is. Add to the ID and make it stage 3.
						stage = 3
						CmpID += string(ob[0])
					} else {
						// Not a mention.
						return nil
					}
				}
			} else if stage == 3 {
				if ob[0] == '>' {
					// This is the end. We should return here.
					return &CmpID
				} else if ob[0] > 46 && 58 > ob[0] {
					// Append to the ID.
					CmpID += string(ob[0])
				} else {
					// Return nil if this is a mention. If it is, rewind one and return the ID.
					if mention {
						return nil
					}
					_, _ = r.Seek(-1, io.SeekCurrent)
					return &CmpID
				}
			}
		}
	}
}
