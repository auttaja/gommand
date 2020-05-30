package gommand

import "testing"

// TestMiddleware is used to test the middleware functionality.
func TestMiddleware(t *testing.T) {
	patchMember = false
	r := NewRouter(&RouterConfig{
		PrefixCheck: StaticPrefix("%"),
		Middleware: []Middleware{
			func(ctx *Context) error {
				ctx.MiddlewareParams["b"] = 1
				return nil
			},
		},
	})
	r.SetCommand(&Command{
		Name: "middleware",
		Middleware: []Middleware{
			func(ctx *Context) error {
				ctx.MiddlewareParams["a"] = 1
				return nil
			},
		},
		Function: func(ctx *Context) error {
			m, ok := ctx.MiddlewareParams["a"].(int)
			if !ok || m != 1 {
				t.Log("command middleware fail")
				t.FailNow()
				return nil
			}
			m, ok = ctx.MiddlewareParams["b"].(int)
			if !ok || m != 1 {
				t.Log("global middleware fail")
				t.FailNow()
				return nil
			}
			return nil
		},
	})
	r.AddErrorHandler(func(_ *Context, err error) bool {
		t.Log(err)
		t.FailNow()
		return true
	})
	r.CommandProcessor(nil, mockMessage("%middleware"), true)
}
