package action 

import (
	"fmt"
	//"os"
	"github.com/pkg/errors"
	"github.com/Code-Hex/exit" 																		// debug-exit
)

type exiter interface {
	ExitCode() int
}

type causer interface {
	Cause() error
}

func unwrapErrors(err error) (int, error) {
	for e := err; e != nil; {
		switch e.(type) {
		case exiter:
			return e.(exiter).ExitCode(), e
		case causer:
			e = e.(causer).Cause()
		default:
			return 1, e // default error
		}
	}
	return 0, nil
}

func getExitCode(err error) int {
	for e := err; e != nil; {
		switch e.(type) {
		case exiter:
			return e.(exiter).ExitCode()
		case causer:
			e = e.(causer).Cause()
		}
	}
	return exit.OK
}

func fileOperate() error {
	if err := doSomething1(); err != nil {
		return err
	}
	return nil
}

func doSomething3() error {
	return exit.MakeOSFile(errors.New("Failed to operate files"))
}

func doSomething2() error {
	if err := doSomething3(); err != nil {
		return errors.Wrap(err, "Second error")
	}
	return nil
}

func doSomething1() error {
	if err := doSomething2(); err != nil {
		return errors.Wrap(err, "Third error")
	}
	return nil
}

func debugExit(trace bool) int {
	if err := fileOperate(); err != nil {
		var code int
		if trace {
			code = getExitCode(err)
			fmt.Printf("%+v\n", err)
		} else {
			code, err = unwrapErrors(err)
			fmt.Printf("%s\n", err.Error())
		}
		return code
	}

	return exit.OK
}