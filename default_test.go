package gommand

import (
	"fmt"
	"reflect"
	"testing"
)

func TestDefault(t *testing.T) {
	tables := []struct {
		rawArgs  string
		expected []interface{}
	}{
		{"42 \"test string\"", []interface{}{[]interface{}{42}, "test string", "string 2"}},
		{"5 \"test string\"", []interface{}{[]interface{}{5}, "test string", "string 2"}},
		{"1 \"test string\"", []interface{}{[]interface{}{1}, "test string", "string 2"}},
		{"1 2 3 4 hello", []interface{}{[]interface{}{1, 2, 3, 4}, "hello", "string 2"}},
		{"1", []interface{}{[]interface{}{1}, "test string", "string 2"}},
		{"1 hello \"replace string 2\"", []interface{}{[]interface{}{1}, "hello", "replace string 2"}},
	}

	test := 0

	r := NewRouter(&RouterConfig{
		PrefixCheck: StaticPrefix("%"),
	})

	r.SetCommand(&Command{
		Name: "default",
		ArgTransformers: []ArgTransformer{
			{
				Function: IntTransformer,
				Greedy:   true,
			},
			{
				Function: StringTransformer,
				Default:  "test string",
			},
			{
				Function: StringTransformer,
				Default:  "string 2",
			},
		},
		Function: func(ctx *Context) error {
			if !reflect.DeepEqual(ctx.Args, tables[test].expected) {
				// reflect.DeepEqual is slow but should be fine for test cases
				t.Log(fmt.Sprintf("test %d failed", test))
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

	for _, j := range tables {
		r.CommandProcessor(nil, 0, mockMessage("%default "+j.rawArgs), true)
		test++
	}
}
