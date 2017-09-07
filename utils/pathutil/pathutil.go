package pathutil

import (
	// "encoding/json"
	//"fmt"
	//"os"
	"os"
)

// https://github.com/jakdept/dir/blob/master/dir_test.go
// watch + create

type RootFolder struct {
	Repo string
	Path string
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