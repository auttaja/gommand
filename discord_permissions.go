package gommand

import "context"

// PermissionCheckSettings is used to define the settings which are used for checking permissions.
type PermissionCheckSettings uint8

const (
	// CheckMembersUserPermissions is used to check the members user permissions.
	CheckMembersUserPermissions PermissionCheckSettings = 1 << iota

	// CheckMembersChannelPermissions is used to check the members channel permissions.
	CheckMembersChannelPermissions

	// CheckBotUserPermissions is used to check the bots user permissions.
	CheckBotUserPermissions

	// CheckBotChannelPermissions is used to check the bots channel permissions.
	CheckBotChannelPermissions
)

// Used to wrap permissions.
func permissionsWrapper(PermissionName string, PermissionsHex uint64) func(Checks PermissionCheckSettings) func(ctx *Context) (string, bool) {
	return func(Checks PermissionCheckSettings) func(ctx *Context) (string, bool) {
		// Get all of the check types.
		checkMemberUser := Checks&CheckMembersUserPermissions != 0
		checkMemberChannel := Checks&CheckMembersChannelPermissions != 0
		checkBotUser := Checks&CheckBotUserPermissions != 0
		checkBotChannel := Checks&CheckBotChannelPermissions != 0

		// If all are false, assume we want to check the member user.
		if !checkBotChannel && !checkBotUser && !checkMemberChannel && !checkMemberUser {
			checkMemberUser = true
		}

		// Create the array of checks we want.
		checks := make([]func(ctx *Context) (string, bool), 0, 1)

		// Run any member checks.
		if checkMemberChannel || checkMemberUser {
			f := func(ctx *Context) (string, bool) {
				guild, err := ctx.Guild()
				if err != nil {
					return err.Error(), false
				}
				if guild.OwnerID == ctx.Message.Author.ID {
					return "", true
				}
				var perms uint64
				if checkMemberChannel {
					c, err := ctx.Channel()
					if err != nil {
						return err.Error(), false
					}
					perms, err = c.GetPermissions(context.TODO(), ctx.Session, ctx.Message.Member)
				} else {
					perms, err = ctx.Message.Member.GetPermissions(context.TODO(), ctx.Session)
				}
				if err != nil {
					return err.Error(), false
				}

				// 0x00000008 is the ADMINISTRATOR permission hex code and bypasses everything.
				if (perms & 0x00000008) == 0x00000008 {
					return "", true
				}

				return "You must have the  \"" + PermissionName + "\" permission to run this command.", (perms & PermissionsHex) == PermissionsHex
			}
			checks = append(checks, f)
		}

		// Run any bot user checks.
		if checkBotChannel || checkBotUser {
			f := func(ctx *Context) (string, bool) {
				member, err := ctx.BotMember()
				if err != nil {
					return err.Error(), false
				}

				var perms uint64
				if checkBotChannel {
					c, err := ctx.Channel()
					if err != nil {
						return err.Error(), false
					}
					perms, err = c.GetPermissions(context.TODO(), ctx.Session, member)
				} else {
					perms, err = member.GetPermissions(context.TODO(), ctx.Session)
				}
				if err != nil {
					return err.Error(), false
				}

				// 0x00000008 is the ADMINISTRATOR permission hex code and bypasses everything.
				if (perms & 0x00000008) == 0x00000008 {
					return "", true
				}

				return "The bot must have the  \"" + PermissionName + "\" permission to run this command.", (perms & PermissionsHex) == PermissionsHex
			}
			checks = append(checks, f)
		}

		// Return something to run through them.
		return func(ctx *Context) (string, bool) {
			for _, check := range checks {
				s, ok := check(ctx)
				if !ok {
					return s, ok
				}
			}
			return "", true
		}
	}
}

// CREATE_INSTANT_INVITE is a wrapper for the Discord permission.
var CREATE_INSTANT_INVITE = permissionsWrapper("Create Instant Invite", 0x00000001)

// KICK_MEMBERS is a wrapper for the Discord permission.
var KICK_MEMBERS = permissionsWrapper("Kick Members", 0x00000002)

// BAN_MEMBERS is a wrapper for the Discord permission.
var BAN_MEMBERS = permissionsWrapper("Ban Members", 0x00000004)

// ADMINISTRATOR is a wrapper for the Discord permission.
var ADMINISTRATOR = permissionsWrapper("Administrator", 0x00000008)

// MANAGE_CHANNELS is a wrapper for the Discord permission.
var MANAGE_CHANNELS = permissionsWrapper("Manage Channels", 0x00000010)

// MANAGE_GUILD is a wrapper for the Discord permission.
var MANAGE_GUILD = permissionsWrapper("Manage Guild", 0x00000020)

// ADD_REACTIONS is a wrapper for the Discord permission.
var ADD_REACTIONS = permissionsWrapper("Add Reactions", 0x00000040)

// VIEW_AUDIT_LOG is a wrapper for the Discord permission.
var VIEW_AUDIT_LOG = permissionsWrapper("View Audit Log", 0x00000080)

// VIEW_CHANNEL is a wrapper for the Discord permission.
var VIEW_CHANNEL = permissionsWrapper("View Channel", 0x00000400)

// SEND_MESSAGES is a wrapper for the Discord permission.
var SEND_MESSAGES = permissionsWrapper("Send Messages", 0x00000800)

// SEND_TTS_MESSAGES is a wrapper for the Discord permission.
var SEND_TTS_MESSAGES = permissionsWrapper("Send Tts Messages", 0x00001000)

// MANAGE_MESSAGES is a wrapper for the Discord permission.
var MANAGE_MESSAGES = permissionsWrapper("Manage Messages", 0x00002000)

// EMBED_LINKS is a wrapper for the Discord permission.
var EMBED_LINKS = permissionsWrapper("Embed Links", 0x00004000)

// ATTACH_FILES is a wrapper for the Discord permission.
var ATTACH_FILES = permissionsWrapper("Attach Files", 0x00008000)

// READ_MESSAGE_HISTORY is a wrapper for the Discord permission.
var READ_MESSAGE_HISTORY = permissionsWrapper("Read Message History", 0x00010000)

// MENTION_EVERYONE is a wrapper for the Discord permission.
var MENTION_EVERYONE = permissionsWrapper("Mention Everyone", 0x00020000)

// USE_EXTERNAL_EMOJIS is a wrapper for the Discord permission.
var USE_EXTERNAL_EMOJIS = permissionsWrapper("Use External Emojis", 0x00040000)

// CONNECT is a wrapper for the Discord permission.
var CONNECT = permissionsWrapper("Connect", 0x00100000)

// SPEAK is a wrapper for the Discord permission.
var SPEAK = permissionsWrapper("Speak", 0x00200000)

// MUTE_MEMBERS is a wrapper for the Discord permission.
var MUTE_MEMBERS = permissionsWrapper("Mute Members", 0x00400000)

// DEAFEN_MEMBERS is a wrapper for the Discord permission.
var DEAFEN_MEMBERS = permissionsWrapper("Deafen Members", 0x00800000)

// MOVE_MEMBERS is a wrapper for the Discord permission.
var MOVE_MEMBERS = permissionsWrapper("Move Members", 0x01000000)

// USE_VAD is a wrapper for the Discord permission.
var USE_VAD = permissionsWrapper("Use Vad", 0x02000000)

// PRIORITY_SPEAKER is a wrapper for the Discord permission.
var PRIORITY_SPEAKER = permissionsWrapper("Priority Speaker", 0x00000100)

// STREAM is a wrapper for the Discord permission.
var STREAM = permissionsWrapper("Stream", 0x00000200)

// CHANGE_NICKNAME is a wrapper for the Discord permission.
var CHANGE_NICKNAME = permissionsWrapper("Change Nickname", 0x04000000)

// MANAGE_NICKNAMES is a wrapper for the Discord permission.
var MANAGE_NICKNAMES = permissionsWrapper("Manage Nicknames", 0x08000000)

// MANAGE_ROLES is a wrapper for the Discord permission.
var MANAGE_ROLES = permissionsWrapper("Manage Roles", 0x10000000)

// MANAGE_WEBHOOKS is a wrapper for the Discord permission.
var MANAGE_WEBHOOKS = permissionsWrapper("Manage Webhooks", 0x20000000)

// MANAGE_EMOJIS is a wrapper for the Discord permission.
var MANAGE_EMOJIS = permissionsWrapper("Manage Emojis", 0x40000000)
