package gommand

import (
	"context"
	"fmt"
	"sync"

	"github.com/andersfylling/disgord"
)

// This is used to represent all of the current menus.
var menuCache = map[disgord.Snowflake]*EmbedMenu{}

// This is the thread lock for the menu cache.
var menuCacheLock = sync.RWMutex{}

// MenuInfo contains the information about the menu.
type MenuInfo struct {
	Author string
	Info   []string
}

// MenuButton is the datatype containing information about the button.
type MenuButton struct {
	Emoji       string
	Name        string
	Description string
}

// MenuReaction represents the button and the function which it triggers.
type MenuReaction struct {
	Button   *MenuButton
	Function func(ChannelID, MessageID disgord.Snowflake, _ *EmbedMenu, client disgord.Session)
}

// MenuReactions are all of the reactions which the menu has.
type MenuReactions struct {
	ReactionSlice []MenuReaction
}

// EmbedMenu is the base menu.
type EmbedMenu struct {
	Reactions *MenuReactions
	parent    *EmbedMenu
	Embed     *disgord.Embed
	MenuInfo  *MenuInfo
}

// Add is used to add a menu reaction.
func (mr *MenuReactions) Add(reaction MenuReaction) {
	Slice := append(mr.ReactionSlice, reaction)
	mr.ReactionSlice = Slice
}

// AddParentMenu is used to add a new parent menu.
func (e *EmbedMenu) AddParentMenu(Menu *EmbedMenu) {
	e.parent = Menu
}

// Display is used to show a menu. This is un-protected so that people can write their own things on top of embed menus, but you probably want to use ctx.DisplayEmbedMenu(menu).
func (e *EmbedMenu) Display(ChannelID, MessageID disgord.Snowflake, client disgord.Session) error {
	menuCacheLock.Lock()
	if len(e.Reactions.ReactionSlice) == 0 {
		delete(menuCache, MessageID)
	} else {
		menuCache[MessageID] = e
	}
	menuCacheLock.Unlock()

	EmbedCopy := e.Embed.DeepCopy().(*disgord.Embed)
	Fields := make([]*disgord.EmbedField, 0)
	for _, k := range e.Reactions.ReactionSlice {
		Fields = append(Fields, &disgord.EmbedField{
			Name:   fmt.Sprintf("%s %s", k.Button.Emoji, k.Button.Name),
			Value:  k.Button.Description,
			Inline: false,
		})
	}
	EmbedCopy.Fields = append(EmbedCopy.Fields, Fields...)

	_, err := client.UpdateMessage(context.TODO(), ChannelID, MessageID).SetContent("").SetEmbed(EmbedCopy).Execute()
	if err != nil {
		return err
	}
	for _, k := range e.Reactions.ReactionSlice {
		err := client.CreateReaction(context.TODO(), ChannelID, MessageID, k.Button.Emoji)
		if err != nil {
			return err
		}
	}
	return nil
}

// ChildMenuOptions are options which can be set when creating a child menu.
type ChildMenuOptions struct {
	Embed        *disgord.Embed `json:"embed"`
	Button       *MenuButton    `json:"button"`
	BeforeAction func()         `json:"-"`
	AfterAction  func()         `json:"-"`
}

// NewChildMenu is used to create a new child menu.
func (e *EmbedMenu) NewChildMenu(options *ChildMenuOptions) *EmbedMenu {
	NewEmbedMenu := &EmbedMenu{
		Reactions: &MenuReactions{
			ReactionSlice: []MenuReaction{},
		},
		Embed:    options.Embed,
		MenuInfo: e.MenuInfo,
	}
	NewEmbedMenu.parent = e
	Reaction := MenuReaction{
		Button: options.Button,
		Function: func(ChannelID, MessageID disgord.Snowflake, _ *EmbedMenu, client disgord.Session) {
			if options.BeforeAction != nil {
				options.BeforeAction()
			}
			_ = client.DeleteAllReactions(context.TODO(), ChannelID, MessageID)
			_ = NewEmbedMenu.Display(ChannelID, MessageID, client)
			if options.AfterAction != nil {
				options.AfterAction()
			}
		},
	}
	e.Reactions.Add(Reaction)
	return NewEmbedMenu
}

// AddBackButton is used to add a back Button to the page.
func (e *EmbedMenu) AddBackButton() {
	Reaction := MenuReaction{
		Button: &MenuButton{
			Description: "Goes back to the parent menu.",
			Name:        "Back",
			Emoji:       "⬆",
		},
		Function: func(ChannelID, MessageID disgord.Snowflake, _ *EmbedMenu, client disgord.Session) {
			_ = client.DeleteAllReactions(context.TODO(), ChannelID, MessageID)
			_ = e.parent.Display(ChannelID, MessageID, client)
		},
	}
	e.Reactions.Add(Reaction)
}

// NewEmbedMenu is used to create a new menu handler.
func NewEmbedMenu(embed *disgord.Embed, ctx *Context) *EmbedMenu {
	var reactions []MenuReaction
	menu := &EmbedMenu{
		Reactions: &MenuReactions{
			ReactionSlice: reactions,
		},
		Embed: embed,
		MenuInfo: &MenuInfo{
			Author: ctx.Message.Author.ID.String(),
			Info:   []string{},
		},
	}
	return menu
}

// This is used to handle menu reactions.
func handleMenuReactionEdit(s disgord.Session, evt *disgord.MessageReactionAdd) {
	go func() {
		// Get the menu if it exists.
		menuCacheLock.RLock()
		menu := menuCache[evt.MessageID]
		if menu == nil {
			menuCacheLock.RUnlock()
			return
		}
		menuCacheLock.RUnlock()

		// Remove the reaction.
		_ = s.DeleteUserReaction(context.TODO(), evt.ChannelID, evt.MessageID, evt.UserID, evt.PartialEmoji)

		// Check the author of the reaction.
		if menu.MenuInfo.Author != evt.UserID.String() {
			return
		}

		for _, v := range menu.Reactions.ReactionSlice {
			if v.Button.Emoji == evt.PartialEmoji.Name {
				v.Function(evt.ChannelID, evt.MessageID, menu, s)
				return
			}
		}
	}()
}

// Handle messages being deleted to stop memory leaks.
func handleEmbedMenuMessageDelete(s disgord.Session, evt *disgord.MessageDelete) {
	go func() {
		menuCacheLock.Lock()
		delete(menuCache, evt.MessageID)
		menuCacheLock.Unlock()
	}()
}
