package gommand

import (
	"strings"
	"testing"
)

// TestSingleArgument is used to test a basic echo command.
func TestSingleArgument(t *testing.T) {
	r := NewRouter(&RouterConfig{
		PrefixCheck: StaticPrefix("%"),
	})
	r.SetCommand(&Command{
		Name: "echo",
		ArgTransformers: []ArgTransformer{
			{
				Function: StringTransformer,
			},
		},
		Function: func(ctx *Context) error {
			arg := ctx.Args[0].(string)
			if arg != "hello" {
				t.Log(arg, "was returned")
				t.FailNow()
			}
			return nil
		},
	})
	r.AddErrorHandler(func(_ *Context, err error) bool {
		t.Log(err)
		t.FailNow()
		return true
	})
	r.CommandProcessor(nil, 0, mockMessage("%echo hello"), true)
	r.CommandProcessor(nil, 0, mockMessage("%echo \"hello\""), true)
}

// BenchmarkSingleArgument is used to benchmark a basic echo command.
func BenchmarkSingleArgument(b *testing.B) {
	r := NewRouter(&RouterConfig{
		PrefixCheck: StaticPrefix("%"),
	})
	r.SetCommand(&Command{
		Name: "echo",
		ArgTransformers: []ArgTransformer{
			{
				Function: StringTransformer,
			},
		},
		Function: func(ctx *Context) error {
			arg := ctx.Args[0].(string)

			// Do something here to cause a bit of activity in the function.
			//nolint:staticcheck
			strings.ToLower(arg)

			return nil
		},
	})
	r.CommandProcessor(nil, 0, mockMessage("%echo \"hello\""), true)
}
