package utils

import (
	"errors"
	"io"
	// "log"
	"os"
	"strings"
	"flag"
	// "log"
	//"net/http"
	"path/filepath"
	"regexp"
	// logs
	"github.com/sirupsen/logrus"
)
// https://github.com/jakdept/dir/blob/master/dir_test.go
// watch + create

var (
	dir 		= flag.String("dir", "", "directory to monitor and index")
	indexPath 	= flag.String("index", "wiki.bleve", "path to store index")
	pathFilter 	= flag.String("pathFilter", `\.md$`, "regular expression that file names must match")
	staticEtag 	= flag.String("staticEtag", "", "A static etag value.")
	staticPath 	= flag.String("static", "static/", "Path to the static content")
	bindAddr 	= flag.String("addr", ":8099", "http listen address")
	pathRegexp 	*regexp.Regexp
	NoPathError error
)

type RootFolder struct {
	Repo string
	Path string
}

func init() {
	NoPathError = errors.New("Could not get home path from env vars HOME or USERPROFILE")
}

func defaultPath() string {
	home, err := homePath()
	if err != nil {
		log.WithError(err).WithFields(logrus.Fields{"action": "defaultPath", "step": "homePath", "service": "filesystem"}).Warn("Could not get home path from env vars HOME or USERPROFILE")
		//log.Fatal(err)
	}
	return filepath.Join(home, "src")
}

func relativePath(path string) string {
	if strings.HasPrefix(path, *dir) {
		path = path[len(*dir)+1:]
	}
	return path
}

// NoPathError thrown when home path could not automatically be determined

func homePath() (string, error) {
	value := ""
	for _, key := range []string{"HOME", "USERPROFILE"} {
		value = os.Getenv(key)
		if value != "" {
			return value, nil
		}
	}
	if NoPathError != nil {
		log.WithError(NoPathError).WithFields(logrus.Fields{"action": "homePath", "step": "Getenv", "service": "filesystem"}).Warn("Could not get home path from env vars HOME or USERPROFILE")
	}
	return "", NoPathError
}

func createMissingFolder(folder string) (bool, error) {
	folderExists, err := exists(folder)
	if err != nil {
		return false, err
	}

	if !folderExists {
		err = os.MkdirAll(folder, 0700)
		if err != nil {
			return true, err
		}
	}

	return true, nil
}

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

func appendToFile(filename string, content []byte) error {
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0700)
	if err != nil {
		return err
	}

	defer f.Close()

	if _, err = f.WriteAt(content, 0); err != nil {
		return err
	}

	return nil
}

func dirIsEmpty(name string) (bool, error) {
	f, err := os.Open(name)
	if err != nil {
		return false, err
	}
	defer f.Close()

	_, err = f.Readdirnames(1) // Or f.Readdir(1)
	if err == io.EOF {
		return true, nil
	}
	return false, err // Either not empty or error, suits both cases
}

func rmDir(path string) error {
	if VERBOSE {
		log.Printf("rmDir?: %s\n", path)
	}

	if ok, _ := dirIsEmpty(path); ok {
		if VERBOSE {
			log.Printf("Removing %s\n", path)
		}
		return os.Remove(path)
	}

	return nil
}
