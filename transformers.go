package gommand

import "strconv"

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
