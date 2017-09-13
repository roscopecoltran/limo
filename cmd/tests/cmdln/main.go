package main

import (
    "fmt"
    "os"
    "os/exec"
	"github.com/roscopecoltran/sniperkit-limo/utils/cmdline"
	// logs
	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

var log = logrus.New()

func init() {
	formatter := new(prefixed.TextFormatter)
	log.Formatter = formatter
	log.Level = logrus.DebugLevel
}

func main() {
    // here is a somewhat complex command line to execute
	fullCommand := `echo -e 'Starting LS\n===========' && ls -la && echo -e "===========\nI'm Done."`

    // split, so we can run this on the command line...
    cmmd, args := cmdline.Split(fullCommand)

    o, err := os.Command(cmmd, args...).Output()
	if err != nil {
		log.WithError(err).WithFields(logrus.Fields{
			"file":    "main.go",
			"function": "Command",
		}).Fatalf("Could not run the command properly: %#s", err)
	}
    
    // output the current directory
    fmt.Println(string(o))

}