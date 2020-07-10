package gommand

import (
	"context"
	"errors"
	"strconv"

	"github.com/andersfylling/disgord"
)

// EmbedsPaginator is used to paginate together several embeds.
func EmbedsPaginator(ctx *Context, Pages []*disgord.Embed, InitialPage uint, NoButtonTextContent string) error {
	// Check the permissions which the bot has permission to use embed menus in this channel.
	c, err := ctx.Channel()
	if err != nil {
		return err
	}
	m, err := ctx.BotMember()
	if err != nil {
		return err
	}
	perms, err := c.GetPermissions(context.TODO(), ctx.Session, m)
	if err != nil {
		return err
	}
	UseEmbedMenus := (perms&disgord.PermissionManageMessages)==disgord.PermissionManageMessages&&(perms&disgord.PermissionAddReactions)==disgord.PermissionAddReactions

	// Get the pages length.
	PagesLen := len(Pages)

	// If initial page is greater than the length of pages, make it 0.
	if int(InitialPage) > PagesLen {
		InitialPage = 1
	}

	// If length of pages is 0, throw an error.
	if PagesLen == 0 {
		return errors.New("pages length is 0")
	}

	// Defines the current page.
	CurrentPage := 1

	// Modify the embed to be ready.
	PrepareEmbed := func(em *disgord.Embed) *disgord.Embed {
		em = em.DeepCopy().(*disgord.Embed)
		em.Footer = &disgord.EmbedFooter{
			Text: "Page " + strconv.Itoa(CurrentPage) + "/" + strconv.Itoa(PagesLen),
		}
		return em
	}

	// Defines the display and last page.
	var DisplayPage *EmbedMenu
	var DisplayEmbed *disgord.Embed
	var LastPage *EmbedMenu

	// Iterate through the pages.
	for i, em := range Pages {
		// Prepare the embed.
		em = PrepareEmbed(em)

		// Make the embed menu if we are using this method.
		if UseEmbedMenus {
			if LastPage == nil {
				LastPage = NewEmbedMenu(em, ctx)
			} else {
				PageBefore := LastPage
				LastPage = LastPage.NewChildMenu(&ChildMenuOptions{
					Embed: em,
					Button: &MenuButton{
						Emoji:       "▶️",
						Name:        "Forward",
						Description: "Goes forward a page.",
					},
				})
				LastPage.Reactions.Add(MenuReaction{
					Button: &MenuButton{
						Emoji:       "◀️",
						Name:        "Back",
						Description: "Goes back a page.",
					},
					Function: func(ChannelID, MessageID disgord.Snowflake, _ *EmbedMenu, client disgord.Session) {
						_ = PageBefore.Display(ChannelID, MessageID, client)
					},
				})
			}
			if uint(i)+1 == InitialPage {
				DisplayPage = LastPage
			}
		} else if uint(i)+1 == InitialPage {
			DisplayEmbed = em
		}

		// Add to the page.
		CurrentPage++
	}

	// Return displaying the embed.
	if UseEmbedMenus {
		return ctx.DisplayEmbedMenu(DisplayPage)
	} else {
		return func() (err error) { _, err = ctx.Reply(NoButtonTextContent, DisplayEmbed); return }()
	}
}
