package gommand

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/andersfylling/disgord"
)

// This is used to represent all of the current menus.
var menuCache = map[disgord.Snowflake]*EmbedMenu{}

// This is the thread lock for the menu cache.
var menuCacheLock = sync.RWMutex{}

// This is used to hold the lifetime information (if applicable) for each active menu.
var menuLifetimeCache = map[disgord.Snowflake]*EmbedLifetimeOptions{}

// This is the thread lock for the lifetimeCache.
var menuLifetimeCacheLock = sync.RWMutex{}

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

	myID disgord.Snowflake
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
	_ = client.Channel(ChannelID).Message(MessageID).DeleteAllReactions()

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
		if k.Button.Name == "" || k.Button.Description == "" {
			continue
		}

		emojiFormatted := k.Button.Emoji
		if strings.Contains(k.Button.Emoji, ":") {
			if strings.HasPrefix(k.Button.Emoji, "a:") {
				emojiFormatted = "<" + emojiFormatted + ">"
			} else {
				emojiFormatted = "<:" + emojiFormatted + ">"
			}
		}
		Fields = append(Fields, &disgord.EmbedField{
			Name:   fmt.Sprintf("%s %s", emojiFormatted, k.Button.Name),
			Value:  k.Button.Description,
			Inline: false,
		})
	}
	EmbedCopy.Fields = append(EmbedCopy.Fields, Fields...)

	msgRef := client.Channel(ChannelID).Message(MessageID)
	_, err := msgRef.Update().SetContent("").SetEmbed(EmbedCopy).Execute()
	if err != nil {
		return err
	}
	for _, k := range e.Reactions.ReactionSlice {
		err := msgRef.Reaction(k.Button.Emoji).Create()
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
		myID:     e.myID,
	}
	NewEmbedMenu.parent = e
	Reaction := MenuReaction{
		Button: options.Button,
		Function: func(ChannelID, MessageID disgord.Snowflake, _ *EmbedMenu, client disgord.Session) {
			if options.BeforeAction != nil {
				options.BeforeAction()
			}
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
			_ = e.parent.Display(ChannelID, MessageID, client)
		},
	}
	e.Reactions.Add(Reaction)
}

// AddExitButton is used to add an Exit button to the page, deleting the menu if pressed.
func (e *EmbedMenu) AddExitButton() {
	Reaction := MenuReaction{
		Button: &MenuButton{
			Description: "Exits the current menu.",
			Name:        "Exit",
			Emoji:       "❌",
		},
		Function: func(ChannelID, MessageID disgord.Snowflake, _ *EmbedMenu, client disgord.Session) {
			_ = client.Channel(ChannelID).Message(MessageID).Delete()
		},
	}
	e.Reactions.Add(Reaction)
}

// NewEmbedMenu is used to create a new menu handler.
func NewEmbedMenu(embed *disgord.Embed, ctx *Context) *EmbedMenu {
	var reactions []MenuReaction
	menu := &EmbedMenu{
		myID: ctx.BotUser.ID,
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

// EmbedLifetimeOptions represents the options used to control the lifetime of a menu, all fields are optional.
type EmbedLifetimeOptions struct {
	// The maximum time after which an embed is created for which it can exist.
	// After this time has passed, if it hasn't already then the menu will be deleted.
	MaximumLifetime time.Duration

	// The maximum time after an embed was last interacted with (reacted to) that it can exist for.
	// After this time has passed, if it hasn't already then the menu will be deleted.
	InactiveLifetime time.Duration

	// The function called before the menu should be deleted.
	BeforeDelete func()

	// The function called after the message is deleted (only if deleted through exceeded lifetime).
	// Called regardless of errors returned when deleting the message.
	AfterDelete func()

	maxLifetimeTimer *time.Timer

	inactiveTimer *time.Timer
}

// Start inits the timers for the lifetime.
func (l *EmbedLifetimeOptions) Start(ChannelID, MessageID disgord.Snowflake, client disgord.Session) {
	if l.MaximumLifetime == time.Duration(0) && l.InactiveLifetime == time.Duration(0) {
		// This is blank, don't bother caching / updating it.
		return
	}

	menuLifetimeCacheLock.Lock()
	if l.MaximumLifetime > time.Duration(0) {
		// init the maxLifetimeTimer if the MaximumLifetime is a positive non-zero value.
		l.maxLifetimeTimer = time.AfterFunc(l.MaximumLifetime, func() {
			if l.BeforeDelete != nil {
				l.BeforeDelete()
			}
			err := client.Channel(ChannelID).Message(MessageID).Delete()
			if err != nil {
				// If there was an error deleting the message, remove from the menu cache anyway.
				menuCacheLock.Lock()
				delete(menuCache, MessageID)
				menuCacheLock.Unlock()

				menuLifetimeCacheLock.Lock()
				delete(menuLifetimeCache, MessageID)
				menuLifetimeCacheLock.Unlock()
			}
			if l.AfterDelete != nil {
				l.AfterDelete()
			}
		})
	}

	if l.InactiveLifetime > time.Duration(0) {
		l.inactiveTimer = time.AfterFunc(l.InactiveLifetime, func() {
			if l.BeforeDelete != nil {
				l.BeforeDelete()
			}
			err := client.Channel(ChannelID).Message(MessageID).Delete()
			if err != nil {
				// If there was an error deleting the message, remove from the menu cache anyway.
				menuCacheLock.Lock()
				delete(menuCache, MessageID)
				menuCacheLock.Unlock()

				menuLifetimeCacheLock.Lock()
				delete(menuLifetimeCache, MessageID)
				menuLifetimeCacheLock.Unlock()
			}
			if l.AfterDelete != nil {
				l.AfterDelete()
			}
		})
	}

	menuLifetimeCache[MessageID] = l
	menuLifetimeCacheLock.Unlock()
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
		if evt.UserID == menu.myID {
			// This is by me! Do not delete!
			return
		}
		_ = s.Channel(evt.ChannelID).Message(evt.MessageID).Reaction(evt.PartialEmoji).DeleteUser(evt.UserID)

		// Check the author of the reaction.
		if menu.MenuInfo.Author != evt.UserID.String() {
			return
		}

		for _, v := range menu.Reactions.ReactionSlice {
			standardized := ""
			if evt.PartialEmoji.ID == 0 {
				standardized = evt.PartialEmoji.Name
			} else {
				standardized = evt.PartialEmoji.Name + ":" + evt.PartialEmoji.ID.String()
			}

			// We use HasSuffix here because of the "a:" that might be attached.
			if strings.HasSuffix(v.Button.Emoji, standardized) {
				menuLifetimeCacheLock.Lock()
				if lifetime, ok := menuLifetimeCache[evt.MessageID]; ok && lifetime.inactiveTimer != nil {
					if lifetime.inactiveTimer.Stop() {
						// Only reset the timer if it's still "active".
						_ = lifetime.inactiveTimer.Reset(lifetime.InactiveLifetime)
					}
				}
				menuLifetimeCacheLock.Unlock()
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

		menuLifetimeCacheLock.Lock()
		if lifetime, ok := menuLifetimeCache[evt.MessageID]; ok {
			if lifetime.inactiveTimer != nil {
				_ = lifetime.inactiveTimer.Stop()
			}
			if lifetime.maxLifetimeTimer != nil {
				_ = lifetime.maxLifetimeTimer.Stop()
			}
			delete(menuLifetimeCache, evt.MessageID)
		}
		menuLifetimeCacheLock.Unlock()
	}()
}
