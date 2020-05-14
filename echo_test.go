package gommand

import (
	"github.com/andersfylling/disgord"
	"strings"
	"testing"
)

// TestEcho is used to test a basic echo command.
func TestEcho(t *testing.T) {
	patchMember = false
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
	r.msgCmdProcessor(nil, &disgord.MessageCreate{
		Message: &disgord.Message{
			Lockable:        disgord.Lockable{},
			Author:          &disgord.User{Bot: false},
			Timestamp:       disgord.Time{},
			EditedTimestamp: disgord.Time{},
			Content:         "%echo hello",
			Type:            disgord.MessageTypeDefault,
			GuildID:         1,
			Activity:        disgord.MessageActivity{},
			Application:     disgord.MessageApplication{},
		},
	})
	r.msgCmdProcessor(nil, &disgord.MessageCreate{
		Message: &disgord.Message{
			Lockable:        disgord.Lockable{},
			Author:          &disgord.User{Bot: false},
			Timestamp:       disgord.Time{},
			EditedTimestamp: disgord.Time{},
			Content:         "%echo \"hello\"",
			Type:            disgord.MessageTypeDefault,
			GuildID:         1,
			Activity:        disgord.MessageActivity{},
			Application:     disgord.MessageApplication{},
		},
	})
}

// BenchmarkEcho is used to benchmark echo.
func BenchmarkEcho(b *testing.B) {
	patchMember = false
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
			strings.ToLower(arg)

			return nil
		},
	})
	r.msgCmdProcessor(nil, &disgord.MessageCreate{
		Message: &disgord.Message{
			Lockable:        disgord.Lockable{},
			Author:          &disgord.User{Bot: false},
			Timestamp:       disgord.Time{},
			EditedTimestamp: disgord.Time{},
			Content:         "%echo \"hello\"",
			Type:            disgord.MessageTypeDefault,
			GuildID:         1,
			Activity:        disgord.MessageActivity{},
			Application:     disgord.MessageApplication{},
		},
	})
}
