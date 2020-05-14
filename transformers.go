package gommand

import (
	"context"
	"github.com/andersfylling/snowflake/v4"
	"strconv"
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

// IntTransformer is used to transform an arg to a unsigned integer if possible.
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
	id := getUserMention(&StringIterator{Text: Arg})
	if id == nil {
		return
	}
	user, e := ctx.Session.GetUser(context.TODO(), snowflake.ParseSnowflakeString(*id))
	if e == nil {
		err = nil
	}
	return
}

// MemberTransformer is used to transform a member if possible.
func MemberTransformer(ctx *Context, Arg string) (member interface{}, err error) {
	err = &InvalidTransformation{Description: "This was not a valid user ID or mention of someone in this guild."}
	id := getUserMention(&StringIterator{Text: Arg})
	if id == nil {
		return
	}
	member, e := ctx.Session.GetMember(context.TODO(), ctx.Message.GuildID, snowflake.ParseSnowflakeString(*id))
	if e == nil {
		err = nil
	}
	return
}
