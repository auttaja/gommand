package gommand

import (
	"testing"
)

// TestGreedy is used to test the greedy argument parser.
func TestGreedy(t *testing.T) {
	r := NewRouter(&RouterConfig{
		PrefixCheck: StaticPrefix("%"),
	})
	r.SetCommand(&Command{
		Name: "add",
		ArgTransformers: []ArgTransformer{
			{
				Function: IntTransformer,
				Greedy:   true,
			},
			{
				Function: StringTransformer,
			},
		},
		Function: func(ctx *Context) error {
			num := 0
			for _, v := range ctx.Args[0].([]interface{}) {
				num += v.(int)
			}
			if num != 4 {
				t.Log("Added number is", num)
				t.FailNow()
				return nil
			}
			last := ctx.Args[1].(string)
			if last != "hello" {
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
	r.CommandProcessor(nil, mockMessage("%add 1 1 1  1 hello"), true)
}
