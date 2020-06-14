package gommand

import (
	"context"
	"github.com/andersfylling/disgord"
	"github.com/auttaja/fastparse"
	"io"
	"strconv"
	"strings"
	"time"
)

// StringTransformer just takes the argument and returns it.
func StringTransformer(_ *Context, Arg string) (interface{}, error) {
	return Arg, nil
}

// IntTransformer is used to transform an arg to a integer if possible.
func IntTransformer(_ *Context, Arg string) (interface{}, error) {
	i, err := strconv.Atoi(Arg)
	if err != nil {
		return nil, &InvalidTransformation{Description: "Could not transform the argument to an integer."}
	}
	return i, nil
}

// UIntTransformer is used to transform an arg to a unsigned integer if possible.
func UIntTransformer(_ *Context, Arg string) (interface{}, error) {
	i, err := strconv.ParseUint(Arg, 10, 64)
	if err != nil {
		return nil, &InvalidTransformation{Description: "Could not transform the argument to an unsigned integer."}
	}
	return i, nil
}

// UserTransformer is used to transform a user if possible.
func UserTransformer(ctx *Context, Arg string) (user interface{}, err error) {
	err = &InvalidTransformation{Description: "This was not a valid user ID or mention."}
	id := getMention(strings.NewReader(Arg), '@', false)
	if id == nil {
		return
	}
	user, e := ctx.Session.GetUser(context.TODO(), disgord.ParseSnowflakeString(*id))
	if e == nil {
		err = nil
	}
	return
}

// MemberTransformer is used to transform a member if possible.
func MemberTransformer(ctx *Context, Arg string) (member interface{}, err error) {
	err = &InvalidTransformation{Description: "This was not a valid user ID or mention of someone in this guild."}
	id := getMention(strings.NewReader(Arg), '@', false)
	if id == nil {
		return
	}
	member, e := ctx.Session.GetMember(context.TODO(), ctx.Message.GuildID, disgord.ParseSnowflakeString(*id))
	if e == nil {
		err = nil
	}
	return
}

// ChannelTransformer is used to transform a channel if possible.
func ChannelTransformer(ctx *Context, Arg string) (channel interface{}, err error) {
	err = &InvalidTransformation{Description: "This was not a valid channel ID or mention of a channel in this guild."}
	id := getMention(strings.NewReader(Arg), '#', false)
	if id == nil {
		return
	}
	channel, e := ctx.Session.GetChannel(context.TODO(), disgord.ParseSnowflakeString(*id))
	if e == nil {
		err = nil
	}
	return
}

// Gets ID's from the URL if possible.
func getMessageIds(manager *fastparse.ParserManager, start string, iterator io.ReadSeeker) []string {
	urlStart := strings.NewReader(start)
	ob := make([]byte, 1)
	for {
		b, e := urlStart.ReadByte()
		if e != nil {
			break
		}
		_, e = iterator.Read(ob)
		if e != nil {
			return nil
		}
		if ob[0] != b {
			return nil
		}
	}
	p := manager.Parser(iterator)
	defer p.Done()
	s, _ := p.Remainder()
	split := strings.Split(s, "/")
	if len(split) != 3 && len(split) != 4 {
		return nil
	}
	return split
}

// MessageURLTransformer is used to transform a message URL to a message if possible.
func MessageURLTransformer(ctx *Context, Arg string) (message interface{}, err error) {
	err = &InvalidTransformation{Description: "This is not a valid message URL or a message which the bot cannot access."}
	iterator := strings.NewReader(Arg)
	a := getMessageIds(ctx.Router.parserManager, "https://discordapp.com/channels/", iterator)
	if a == nil {
		_, _ = iterator.Seek(0, io.SeekStart)
		a = getMessageIds(ctx.Router.parserManager, "https://discord.com/channels/", iterator)
		if a == nil {
			return
		}
	}
	message, e := ctx.Session.GetMessage(context.TODO(), disgord.ParseSnowflakeString(a[1]), disgord.ParseSnowflakeString(a[2]))
	if e == nil {
		err = nil
	}
	return
}

var str2bool = map[string]bool{
	"y":     true,
	"yes":   true,
	"1":     true,
	"n":     false,
	"no":    false,
	"0":     false,
	"true":  true,
	"false": false,
}

// BooleanTransformer is used to transform an argument into a boolean if possible.
func BooleanTransformer(_ *Context, Arg string) (interface{}, error) {
	boolean, ok := str2bool[strings.ToLower(Arg)]
	if !ok {
		return nil, &InvalidTransformation{Description: "This is not a valid boolean representation."}
	}
	return boolean, nil
}

// RoleTransformer is used to transform a role if possible.
func RoleTransformer(ctx *Context, Arg string) (role interface{}, err error) {
	err = &InvalidTransformation{Description: "This was not a valid role ID or mention of a role in this guild."}
	id := getMention(strings.NewReader(Arg), '@', true)
	if id == nil {
		return
	}
	roles, e := ctx.Session.GetGuildRoles(context.TODO(), ctx.Message.GuildID)
	if e != nil {
		return
	}
	for _, role = range roles {
		if role.(*disgord.Role).ID.String() == *id {
			err = nil
			return
		}
	}
	return
}

// DurationTransformer is used to transform a duration if possible.
func DurationTransformer(_ *Context, Arg string) (duration interface{}, err error) {
	err = &InvalidTransformation{Description: "This was not a valid duration."}
	duration, e := time.ParseDuration(Arg)
	if e == nil {
		err = nil
	}
	return
}

// AnyTransformer takes multiple transformers and tries to find one which works.
func AnyTransformer(Transformers ...func(ctx *Context, Arg string) (interface{}, error)) func(ctx *Context, Arg string) (item interface{}, err error) {
	return func(ctx *Context, Arg string) (item interface{}, err error) {
		err = &InvalidTransformation{Description: "Unable to transform the argument properly."}
		for _, v := range Transformers {
			res, e := v(ctx, Arg)
			if e == nil {
				item = res
				err = nil
				return
			}
		}
		return
	}
}
