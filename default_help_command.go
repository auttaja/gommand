package gommand

import (
	"github.com/andersfylling/disgord"
	"math"
	"strconv"
	"strings"
)

// Creates the embeds for the categories.
func createCategoryEmbeds(ctx *Context, key CategoryInterface, value []*Command) []*disgord.Embed {
	// Defines the number of fields per page.
	FieldsPerPage := 5

	// Ignore commands which the user can't run.
	cmds := make([]*Command, 0, len(value))
	for _, v := range value {
		if v.HasPermission(ctx) == nil {
			cmds = append(cmds, v)
		}
	}
	if len(cmds) == 0 {
		return []*disgord.Embed{}
	}

	// Creates an array full of Discord embeds.
	EmbedsLen := int(math.Ceil(float64(len(value)) / float64(FieldsPerPage)))
	Fields := make([][]*disgord.EmbedField, EmbedsLen)
	CurrentField := 0
	current := make([]*disgord.EmbedField, 0, FieldsPerPage)
	for i, v := range value {
		if i != 0 && i%FieldsPerPage == 0 {
			Fields[CurrentField] = current
			CurrentField++
			current = make([]*disgord.EmbedField, 0, FieldsPerPage)
		}
		Description := v.Description
		if Description == "" {
			Description = "No description set."
		}
		current = append(current, &disgord.EmbedField{
			Name:   ctx.Prefix + v.Name + " " + v.Usage,
			Value:  Description,
			Inline: false,
		})
	}
	if len(current) != 0 {
		Fields[CurrentField] = current
	}
	Embeds := make([]*disgord.Embed, EmbedsLen)
	for i := range Embeds {
		Title := ""
		Description := ""
		if key == nil {
			Title = "General Commands "
			Description = "These commands have not been assigned a category yet."
		} else {
			Title = key.GetName() + " "
			Description = key.GetDescription()
		}
		Title += "[" + strconv.Itoa(i+1) + "/" + strconv.Itoa(EmbedsLen) + "]"
		Embeds[i] = &disgord.Embed{
			Title:       Title,
			Description: Description,
			Color:       2818303,
			Fields:      Fields[i],
		}
	}

	// Return the embeds.
	return Embeds
}

// This sets the default help command.
func defaultHelpCommand() *Command {
	return &Command{
		Name:        "help",
		Description: "Used to get help for a command.",
		Usage:       "[command]",
		ArgTransformers: []ArgTransformer{
			{
				Optional: true,
				Function: StringTransformer,
			},
		},
		Function: func(ctx *Context) error {
			// Get a single command if it is set.
			cmdname, ok := ctx.Args[0].(string)
			if ok {
				cmdname = strings.ToLower(cmdname)
				cmd := ctx.Router.GetCommand(cmdname)
				if cmd == nil {
					_, _ = ctx.Reply(disgord.Embed{
						Title:       "Command not found:",
						Description: "The command \"" + cmdname + "\" was not found.",
						Color:       16711704,
					})
					return nil
				}
				desc := cmd.Description
				if cmd.HasPermission(ctx) != nil {
					desc += "\n\n**You do not have permission to run this.**"
				}
				_, _ = ctx.Reply(disgord.Embed{
					Title:       ctx.Prefix + cmdname + " " + cmd.Usage,
					Description: desc,
					Color:       2818303,
				})
				return nil
			}

			// Get the embed pages.
			pages := make([]*disgord.Embed, 0)
			for k, v := range ctx.Router.GetCommandsOrderedByCategory() {
				pages = append(pages, createCategoryEmbeds(ctx, k, v)...)
			}

			// Send the embed pages.
			_ = EmbedsPaginator(ctx, pages)

			// Return no errors.
			return nil
		},
	}
}
