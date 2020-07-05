// +build ignore

package main

import (
	"io/ioutil"
	"os"
	"text/template"
)

func die(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	b, err := ioutil.ReadFile("wait_for_event.tmpl")
	die(err)

	tpl, err := template.New("wait_for_event").Parse(string(b))
	die(err)

	f, err := os.Create("wait_for_event_gen.go")
	die(err)
	defer f.Close()

	die(tpl.Execute(f, []string{
		"ChannelCreate",
		"ChannelUpdate",
		"ChannelDelete",
		"ChannelPinsUpdate",
		"TypingStart",
		"InviteDelete",
		"MessageCreate",
		"MessageUpdate",
		"MessageDelete",
		"MessageDeleteBulk",
		"MessageReactionAdd",
		"MessageReactionRemove",
		"MessageReactionRemoveAll",
		"GuildEmojisUpdate",
		"GuildCreate",
		"GuildUpdate",
		"GuildDelete",
		"GuildBanAdd",
		"GuildBanRemove",
		"GuildMemberAdd",
		"GuildMemberRemove",
		"GuildMemberUpdate",
		"GuildRoleCreate",
		"GuildRoleUpdate",
		"GuildRoleDelete",
		"PresenceUpdate",
		"UserUpdate",
		"VoiceStateUpdate",
		"VoiceServerUpdate",
		"WebhooksUpdate",
		"InviteCreate",
	}))
}
