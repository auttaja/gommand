package gommand

import "testing"

// TestOptional is used to test the optional argument parser.
func TestOptional(t *testing.T) {
	r := NewRouter(&RouterConfig{
		PrefixCheck: StaticPrefix("%"),
	})

	// Handle one optional argument.
	r.SetCommand(&Command{
		Name: "oneoptional",
		ArgTransformers: []ArgTransformer{
			{
				Function: UIntTransformer,
			},
			{
				Function: StringTransformer,
				Optional: true,
			},
		},
		Function: func(ctx *Context) error {
			strexists := ctx.Args[0].(uint64) == 1
			s, ok := ctx.Args[1].(string)
			if ok {
				if !strexists {
					t.Log("string exists when it shouldn't: ", s)
					t.FailNow()
					return nil
				}
				if s != "test" {
					t.Log("string is", s)
					t.FailNow()
					return nil
				}
			} else if strexists {
				t.Log("string doesn't exist.")
				t.FailNow()
				return nil
			}
			return nil
		},
	})

	// Handle multiple optional arguments.
	r.SetCommand(&Command{
		Name: "multioptional",
		ArgTransformers: []ArgTransformer{
			{
				Function: UIntTransformer,
			},
			{
				Function: StringTransformer,
				Optional: true,
			},
			{
				Function: StringTransformer,
				Optional: true,
			},
		},
		Function: func(ctx *Context) error {
			strexists := ctx.Args[0].(uint64) == 1
			s, ok := ctx.Args[1].(string)
			if ok {
				if !strexists {
					t.Log("string exists when it shouldn't: ", s)
					t.FailNow()
					return nil
				}
				if s != "test" {
					t.Log("string is", s)
					t.FailNow()
					return nil
				}
			} else if strexists {
				t.Log("string doesn't exist.")
				t.FailNow()
				return nil
			}
			s, ok = ctx.Args[2].(string)
			if ok {
				if !strexists {
					t.Log("string exists when it shouldn't: ", s)
					t.FailNow()
					return nil
				}
				if s != "123" {
					t.Log("string is", s)
					t.FailNow()
					return nil
				}
			} else if strexists {
				t.Log("string doesn't exist.")
				t.FailNow()
				return nil
			}
			return nil
		},
	})

	// Handle optional remainders.
	r.SetCommand(&Command{
		Name: "remainder",
		ArgTransformers: []ArgTransformer{
			{
				Function: UIntTransformer,
			},
			{
				Function:  StringTransformer,
				Optional:  true,
				Remainder: true,
			},
		},
		Function: func(ctx *Context) error {
			strexists := ctx.Args[0].(uint64) == 1
			s, ok := ctx.Args[1].(string)
			if ok {
				if !strexists {
					t.Log("string exists when it shouldn't: ", s)
					t.FailNow()
					return nil
				}
				if s != "test" {
					t.Log("string is", s)
					t.FailNow()
					return nil
				}
			} else if strexists {
				t.Log("string doesn't exist.")
				t.FailNow()
				return nil
			}
			return nil
		},
	})

	// Handle greedy remainders.
	r.SetCommand(&Command{
		Name: "greedy",
		ArgTransformers: []ArgTransformer{
			{
				Function: UIntTransformer,
			},
			{
				Function: IntTransformer,
				Optional: true,
				Greedy:   true,
			},
		},
		Function: func(ctx *Context) error {
			argsexist := ctx.Args[0].(uint64) == 1
			x, ok := ctx.Args[1].([]interface{})
			if ok {
				if !argsexist {
					t.Log("array exists when it shouldn't: ", x)
					t.FailNow()
					return nil
				}
				sum := 0
				for _, v := range x {
					sum += v.(int)
				}
				if sum != 4 {
					t.Log("array adds up to", sum)
					t.FailNow()
					return nil
				}
			} else if argsexist {
				t.Log("array doesn't exist.")
				t.FailNow()
				return nil
			}
			return nil
		},
	})

	// Fail the test on error.
	r.AddErrorHandler(func(_ *Context, err error) bool {
		t.Log(err)
		t.FailNow()
		return true
	})

	// Run the commands.
	r.CommandProcessor(nil, mockMessage("%oneoptional 0"), true)
	r.CommandProcessor(nil, mockMessage("%oneoptional 1 test"), true)
	r.CommandProcessor(nil, mockMessage("%multioptional 0"), true)
	r.CommandProcessor(nil, mockMessage("%multioptional 1 test \"123\""), true)
	r.CommandProcessor(nil, mockMessage("%remainder 0"), true)
	r.CommandProcessor(nil, mockMessage("%remainder 1 test"), true)
	r.CommandProcessor(nil, mockMessage("%greedy 0"), true)
	r.CommandProcessor(nil, mockMessage("%greedy 1 1 1  1 1"), true)
}
