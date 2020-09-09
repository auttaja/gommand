package gommand

import (
	"github.com/andersfylling/disgord"
	"math"
	"strconv"
	"strings"
	"time"
)

// Creates the embeds for the categories.
func createCategoryEmbeds(ctx *Context, key CategoryInterface, value []CommandInterface) []*disgord.Embed {
	// Defines the number of fields per page.
	FieldsPerPage := 5

	// Ignore commands which the user can't run.
	cmds := make([]CommandInterface, 0, len(value))
	for _, v := range value {
		if CommandHasPermission(ctx, v) == nil {
			cmds = append(cmds, v)
		}
	}
	if len(cmds) == 0 {
		return []*disgord.Embed{}
	}

	// Creates an array full of Discord embeds.
	EmbedsLen := int(math.Ceil(float64(len(cmds)) / float64(FieldsPerPage)))
	Fields := make([][]*disgord.EmbedField, EmbedsLen)
	CurrentField := 0
	current := make([]*disgord.EmbedField, 0, FieldsPerPage)
	for i, v := range cmds {
		if i != 0 && i%FieldsPerPage == 0 {
			Fields[CurrentField] = current
			CurrentField++
			current = make([]*disgord.EmbedField, 0, FieldsPerPage)
		}
		Description := v.GetDescription()
		if Description == "" {
			Description = "No description set."
		}
		current = append(current, &disgord.EmbedField{
			Name:   ctx.Prefix + v.GetName() + " " + v.GetUsage(),
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
		Usage:       "[page/command]",
		ArgTransformers: []ArgTransformer{
			{
				Optional: true,
				Function: AnyTransformer(UIntTransformer, StringTransformer),
			},
		},
		PermissionValidators: []PermissionValidator{
			EMBED_LINKS(CheckBotChannelPermissions),
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
				desc := cmd.GetDescription()
				if CommandHasPermission(ctx, cmd) != nil {
					desc += "\n\n**You do not have permission to run this.**"
				}
				_, _ = ctx.Reply(disgord.Embed{
					Title:       ctx.Prefix + cmdname + " " + cmd.GetUsage(),
					Description: desc,
					Color:       2818303,
				})
				return nil
			}

			// Get the page if it is set. If not, this will default to 0 which is handled.
			page, _ := ctx.Args[0].(uint64)

			// Get the embed pages.
			pages := make([]*disgord.Embed, 0)
			for k, v := range ctx.Router.GetCommandsOrderedByCategory() {
				pages = append(pages, createCategoryEmbeds(ctx, k, v)...)
			}

			// Send the embed pages.
			_ = EmbedsPaginatorWithLifetime(ctx, pages, uint(page), "Use "+ctx.Prefix+"help <page number> to flick between pages.", &EmbedLifetimeOptions{InactiveLifetime: time.Minute * 5})

			// Return no errors.
			return nil
		},
	}
}
