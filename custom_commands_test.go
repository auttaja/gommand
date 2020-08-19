package gommand

import (
	"github.com/auttaja/fastparse"
	"testing"
)

// TestCustomCommands is used to test custom commands.
func TestCustomCommands(t *testing.T) {
	r := NewRouter(&RouterConfig{
		PrefixCheck: StaticPrefix("%"),
	})
	exists := false
	errored := false
	r.CustomCommandsHandler = func(ctx *Context, cmdname string, parser *fastparse.Parser) (bool, error) {
		if !exists {
			return false, nil
		}
		if cmdname != "test" {
			t.Fatal("invalid command name:", cmdname)
		}
		arg := parser.GetNextArg()
		if arg.Text != "123" {
			t.Fatal("Invalid arguments:", arg.Text)
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
