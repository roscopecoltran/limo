package utils

import (
	"errors"
)

/*
	finder.go
*/
var Found 				= errors.New("Found it!")

/*
	errors.go
*/
var NoPathError 		= errors.New("Could not get home path from env vars HOME or USERPROFILE")
//ErrFileNotSet is used as a return when File needs to have a value, but hasn't been set
var ErrFileNotSet 		= &configError{s: "File has not been set"}

/*
	env.go
*/
var ErrCannotConvert 	= errors.New("cannot convert type")

