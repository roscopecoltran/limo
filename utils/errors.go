package utils

import (
	"fmt"
	"os"

	"github.com/jwaldrip/tint"
)

func exitWithMsg(msgs ...interface{}) {
	fmt.Println(msgs...)
	os.Exit(1)
}

func exitIfErr(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, tint.Colorize(fmt.Sprintf("ERROR: %s", err), tint.Red))
	}
}

func exitIfErrWithMsg(err error, msgs ...interface{}) {
	if err != nil {
		exitWithMsg(msgs...)
	}
}

type configError struct {
	field string
	s string
}

// Error returns a formatted string with the full error message
func (e *configError) Error() string {
	return fmt.Sprintf("%s: %s", e.field, e.s)
}

// GetField returns the field in error
func (e *configError) GetField() string {
	return e.field
}

// IsFieldError returns true if the particular error is related to a field in the configuration file
func (e *configError) IsFieldError() bool {
	return e.field != ""
}

//ErrFileNotSet is used as a return when File needs to have a value, but hasn't been set
var ErrFileNotSet = &configError{s: "File has not been set"}