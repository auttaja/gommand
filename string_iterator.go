package gommand

import "errors"

// StringIterator is used to iterate through a string but keep track of the position.
// This also allows for the ability to rewind a string. This is NOT designed to be thread safe.
type StringIterator struct {
	Text string `json:"text"`
	Pos  uint   `json:"pos"`
	len  *uint  `json:"-"`
}

func (s *StringIterator) ensureLen() {
	if s.len == nil {
		l := uint(len(s.Text))
		s.len = &l
	}
}

// Rewind is used to rewind a string iterator N number of chars.
func (s *StringIterator) Rewind(N uint) {
	s.ensureLen()
	s.Pos -= N
}

// GetChar is used to get a character from the string.
func (s *StringIterator) GetChar() (char uint8, err error) {
	s.ensureLen()
	if s.Pos == *s.len {
		err = errors.New("string is fully iterated")
		return
	}
	char = s.Text[s.Pos]
	s.Pos++
	return
}

// GetRemainder is used to get the remainder of a string.
// FillIterator defines if the iterators count should be affected by this.
func (s *StringIterator) GetRemainder(FillIterator bool) (remainder string, err error) {
	s.ensureLen()
	if s.Pos == *s.len {
		err = errors.New("string is fully iterated")
		return
	}
	remainder = s.Text[s.Pos:]
	if FillIterator {
		s.Pos = *s.len
	}
	return
}
