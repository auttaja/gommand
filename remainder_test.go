package gommand

import (
	"github.com/andersfylling/disgord"
	"testing"
)

// TestRemainder is used to test the remainder argument parser.
func TestRemainder(t *testing.T) {
	patchMember = false
	r := NewRouter(&RouterConfig{
		PrefixCheck: StaticPrefix("%"),
	})
	r.SetCommand(&Command{
		Name: "remainder",
		ArgTransformers: []ArgTransformer{
			{
				Function: StringTransformer,
				Remainder: true,
			},
		},
		Function: func(ctx *Context) error {
			last := ctx.Args[0].(string)
			if last != "\"hello\"" {
				t.Log("String is", last)
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
	r.msgCmdProcessor(nil, &disgord.MessageCreate{
		Message: mockMessage("%remainder   \"hello\""),
	})
}
