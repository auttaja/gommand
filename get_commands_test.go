package gommand

import "testing"

// TestGetCommands is used to test getting commands.
func TestGetCommands(t *testing.T) {
	patchMember = false
	r := NewRouter(&RouterConfig{})
	r.SetCommand(&Command{
		Name: "a",
		Function: func(ctx *Context) error {
			return nil
		},
	})
	cat := &Category{Name: "a"}
	r.SetCommand(&Command{
		Name:     "b",
		Aliases:  []string{"c"},
		Category: cat,
		Function: func(ctx *Context) error {
			return nil
		},
	})
	r.RemoveCommand(r.GetCommand("help"))
	l := len(r.GetAllCommands())
	if l != 2 {
		t.Log("command count should be", l)
		t.FailNow()
		return
	}
	ordered := r.GetCommandsOrderedByCategory()
	if ordered[nil][0].Name != "a" {
		t.Log("command is invalid")
		t.FailNow()
		return
	}
	if ordered[cat][0].Name != "b" {
		t.Log("command is invalid")
		t.FailNow()
		return
	}
}
