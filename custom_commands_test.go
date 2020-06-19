package gommand

import (
	"io"
	"io/ioutil"
	"testing"
)

// TestCustomCommands is used to test custom commands.
func TestCustomCommands(t *testing.T) {
	r := NewRouter(&RouterConfig{
		PrefixCheck: StaticPrefix("%"),
	})
	exists := false
	errored := false
	r.CustomCommandsHandler = func(ctx *Context, cmdname string, r io.ReadSeeker) (bool, error) {
		if !exists {
			return false, nil
		}
		if cmdname != "test" {
			t.Fatal("invalid command name:", cmdname)
		}
		b, _ := ioutil.ReadAll(r)
		if string(b) != "123" {
			t.Fatal("Invalid arguments:", b)
		}
		return true, nil
	}
	r.AddErrorHandler(func(_ *Context, err error) bool {
		if exists {
			t.Fatal(err.Error())
		}
		errored = true
		return true
	})
	r.CommandProcessor(nil, 0, mockMessage("%nonexistent"), true)
	if !errored {
		t.Fatal("didn't error")
	}
	exists = true
	r.CommandProcessor(nil, 0, mockMessage("%test 123"), true)
}
