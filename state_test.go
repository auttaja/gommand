package gommand

import "testing"

// TestState is used to test the state functions.
func TestState(t *testing.T) {
	r := NewRouter(&RouterConfig{
		PrefixCheck: StaticPrefix("%"),
	})
	r.SetCommand(&Command{
		Name: "state",
		Function: func(ctx *Context) error {
			if ctx.State.GetValue() != 0 {
				t.Log("state should have 0 value, current value", ctx.State.GetValue())
				t.FailNow()
			}

			if ctx.State.AddOne() != 1 {
				t.Log("state should have value of 1, current value", ctx.State.GetValue())
				t.FailNow()
			}

			_ = ctx.State.Reset()
			if ctx.State.GetValue() != 0 {
				t.Log("state should have reset to 0, current value", ctx.State.GetValue())
			}
			return nil
		},
	})
	r.CommandProcessor(nil, 0, mockMessage("%state"), true)
}
