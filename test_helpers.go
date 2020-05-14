package gommand

import "github.com/andersfylling/disgord"

func mockMessage(content string) *disgord.Message {
	return &disgord.Message{
		Lockable:        disgord.Lockable{},
		Author:          &disgord.User{Bot: false},
		Timestamp:       disgord.Time{},
		EditedTimestamp: disgord.Time{},
		Content:         content,
		Type:            disgord.MessageTypeDefault,
		GuildID:         1,
		Activity:        disgord.MessageActivity{},
		Application:     disgord.MessageApplication{},
	}
}
