package utils

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"github.com/sirupsen/logrus"
)

func Run(command string, args ...string) error {
	cmd := exec.Command(command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return RunCommand(cmd)
}

func RunSilently(command string, args ...string) error {
	cmd := exec.Command(command, args...)
	cmd.Stdout = ioutil.Discard
	cmd.Stderr = ioutil.Discard
	return RunCommand(cmd)
}

func RunInDir(dir, command string, args ...string) error {
	cmd := exec.Command(command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = dir
	return RunCommand(cmd)
}

type RunFunc func(*exec.Cmd) error

var CommandRunner = func(cmd *exec.Cmd) error {
	return cmd.Run()
}

func RunCommand(cmd *exec.Cmd) error {
	log.WithFields(
		logrus.Fields{	"file": "run.go", 
						"method_name": "RunCommand", 
						}).Infoln("executing command-line: ", cmd.Args[0], strings.Join(cmd.Args[1:], " "))
	err := CommandRunner(cmd)
	if err != nil {
		log.WithError(err).WithFields(
			logrus.Fields{	"file": "run.go", 
							"method_name": "RunCommand", 
							}).Warnf("error occured while executing command-line: ", cmd.Args[0], strings.Join(cmd.Args[1:], " "))
		return &RunError{cmd, err}
	}
	return nil
}

type RunError struct {
	Command   *exec.Cmd
	ExecError error
}

func (e *RunError) Error() string {
	return fmt.Sprintf("%s: %s", e.Command.Path, e.ExecError)
}