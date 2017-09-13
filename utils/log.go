package utils

import (
	"os"
	// logs
	"github.com/sirupsen/logrus"
	"github.com/motemen/go-colorine"
)

var logger = &colorine.Logger{
	colorine.Prefixes{
		"git":      colorine.Verbose,
		"hg":       colorine.Verbose,
		"svn":      colorine.Verbose,
		"darcs":    colorine.Verbose,
		"skip":     colorine.Verbose,
		"cd":       colorine.Verbose,
		"resolved": colorine.Verbose,
		"open":    colorine.Warn,
		"exists":  colorine.Warn,
		"warning": colorine.Warn,
		"authorized": colorine.Notice,
		"error": colorine.Error,
		"": colorine.Info,
	},
}

//func Log(prefix, message string) {
//	logger.Log(prefix, message)
//}

func ErrorIf(err error) bool {
	if err != nil {
		// Log("error", err.Error())
		log.WithError(err).WithFields(
			logrus.Fields{	"action": 	"ErrorIf", 
							"file": 	"utils/log.go"}).Errorf("%s\n", err)
		return true
	}
	return false
}

func DieIf(err error) {
	if err != nil {
		// Log("error", err.Error())
		log.WithError(err).WithFields(
			logrus.Fields{	"action": 	"DieIf", 
							"file": 	"utils/log.go"}).Fatalf("%s\n", err)
		os.Exit(1)
	}
}

func PanicIf(err error) {
	if err != nil {
		log.WithError(err).WithFields(
			logrus.Fields{	"action": 	"PanicIf", 
							"file": 	"utils/log.go"}).Fatalf("%s\n", err)
		panic(err)
	}
}