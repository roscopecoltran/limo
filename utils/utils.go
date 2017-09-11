package utils

import (
	"fmt"
	"os"
	"time"
	// "github.com/davecgh/go-spew/spew"
	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

var VERBOSE = false

var	log 	= logrus.New()

func init() {

	// logs
	log.Out = os.Stdout
	// log.Formatter = new(prefixed.TextFormatter)

	formatter := new(prefixed.TextFormatter)
	formatter.FullTimestamp = true

	// Set specific colors for prefix and timestamp
	formatter.SetColorScheme(&prefixed.ColorScheme{
		PrefixStyle:    "blue+b",
		TimestampStyle: "white+h",
	})

	log.Formatter = formatter

}

// public
func CurrentTime() time.Time {
	return time.Now()
}

func FormatAsDate(t time.Time) string {
	year, month, day := t.Date()

	return fmt.Sprintf("%d/%d/%d", month, day, year)
}

func FormatAsDateTime(t time.Time) string {
	year, month, day := t.Date()

	return fmt.Sprintf("%d/%d/%d @ %d:%d:%d", month, day, year, t.Hour(), t.Minute(), t.Second())
}

func ParseAsDate(timeString string) string {
	stringTime, err := time.Parse("2017-02-03T12:00:00Z07:00", timeString)
	if err != nil {
		Printf("ERROR parsing time: %s, message: %s\n", timeString, err.Error())
		return FormatAsDate(time.Now())
	}

	return FormatAsDate(stringTime)
}

func iterateMonth(ctime time.Time) []string {
	months := []string{}
	now := time.Now()
	for d := ctime; now.After(d); d = d.AddDate(0, 1, 0) {
		months = append(months, fmt.Sprintf("%s", d)[0:7])
	}

	return months
}

func Println(args ...interface{}) (int, error) {
	if VERBOSE {
		return fmt.Println(args...)
	}

	return 0, nil
}

func Printf(format string, args ...interface{}) (int, error) {
	if VERBOSE {
		return fmt.Printf(format, args...)
	}

	return 0, nil
}

/*
func dumpSpew(obj interface{}) string {
	cfg := spew.ConfigState{
		Indent: indentStr,
	}
	return cfg.Sdump(obj)
}
*/