package gommand

import "testing"

// TestSubcommand is used to test that subcommands work properly.
func TestSubcommand(t *testing.T) {
	r := NewRouter(&RouterConfig{
		PrefixCheck: StaticPrefix("%"),
	})
	r.SetCommand(&CommandGroup{
		Name: "a",
		NoCommandSpecified: &Command{
			Function: func(ctx *Context) error {
				return nil
			},
		},
	})
	shouldFail := false
	failed := false
	r.AddErrorHandler(func(_ *Context, err error) bool {
		if shouldFail {
			failed = true
		} else {
			t.Fatal(err)
		}
		return true
	})
	r.CommandProcessor(nil, 0, mockMessage("%a"), true)
	r.SetCommand(&CommandGroup{
		Name: "b",
		subcommands: map[string]CommandInterface{
			"test": &Command{Function: func(ctx *Context) error {
				return nil
			}},
			"arg_expected": &Command{ArgTransformers: []ArgTransformer{{Function: StringTransformer}}, Function: func(ctx *Context) error {
				_ = ctx.Args[0].(string)
				return nil
			}},
		},
	})
	shouldFail = true
	r.CommandProcessor(nil, 0, mockMessage("%b"), true)
	if !failed {
		t.Fatal("command did not fail")
	}
	shouldFail = false
	failed = false
	r.CommandProcessor(nil, 0, mockMessage("%b test"), true)
	shouldFail = true
	r.CommandProcessor(nil, 0, mockMessage("%b arg_expected"), true)
	shouldFail = false
	failed = false
	r.CommandProcessor(nil, 0, mockMessage("%b arg_expected test"), true)
}
