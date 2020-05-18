package gommand

import (
	"context"
	"errors"
	"github.com/andersfylling/disgord"
	"github.com/andersfylling/snowflake/v4"
	"strconv"
)

// EmbedsPaginator is used to paginate together several embeds.
func EmbedsPaginator(ctx *Context, Pages []*disgord.Embed) error {
	// Get the pages length.
	PagesLen := len(Pages)

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

	// Defines the first and last page.
	var FirstPage *EmbedMenu
	var LastPage *EmbedMenu

	// Iterate through the pages.
	for _, em := range Pages {
		// Prepare the embed.
		em = PrepareEmbed(em)

		// Make the embed menu.
		if LastPage == nil {
			FirstPage = NewEmbedMenu(em, ctx)
			LastPage = FirstPage
		} else {
			PageBefore := LastPage
			LastPage = LastPage.NewChildMenu(em, MenuButton{
				Emoji:       "▶️",
				Name:        "Forward",
				Description: "Goes forward a page.",
			})
			LastPage.Reactions.Add(MenuReaction{
				Button: MenuButton{
					Emoji:       "◀️",
					Name:        "Back",
					Description: "Goes back a page.",
				},
				Function: func(ChannelID, MessageID snowflake.Snowflake, _ *EmbedMenu, client disgord.Session) {
					_ = client.DeleteAllReactions(context.TODO(), ChannelID, MessageID)
					_ = PageBefore.Display(ChannelID, MessageID, client)
				},
			})
		}

		// Add to the page.
		CurrentPage++
	}

	// Return displaying the embed.
	return ctx.DisplayEmbedMenu(FirstPage)
}
