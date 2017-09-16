package model

/*
import (
	//"errors"
	//"fmt"
	////"log"
	//"strings"
    // "github.com/qor/qor"
    // "github.com/qor/admin"
	"encoding/json"
	"fmt"
	// "os"
	"path"
	"github.com/jinzhu/gorm"
	// jsoniter "github.com/json-iterator/go"
	// tablib "github.com/agrison/go-tablib"
	// "github.com/davecgh/go-spew/spew"
	// fuzz "github.com/google/gofuzz"
	"github.com/roscopecoltran/sniperkit-limo/utils"
	// "github.com/sirupsen/logrus"
)

// https://github.com/hekar/gitmark/blob/master/gitmark/bookmark.go

type Bookmark struct {
	gorm.Model
	Repo  string
	Title string
	Url   string
}

// RepoResult wraps a repo and an error
type BookmarkResult struct {
	Bookmark  *Bookmark
	Error error
}

func AppendBookmark(db *gorm.DB, rootFolder utils.RootFolder, bookmark *Bookmark) (string, error) {
	appendJson := false
	appendContent := []byte{0}
	if appendJson {
		var err error
		appendContent, err = json.Marshal(bookmark)
		if err != nil {
			return "", err
		}
	} else {
		appendContent = []byte(fmt.Sprintf("\n* [%s](%s)", bookmark.Title, bookmark.Url))
	}
	folder := path.Join(rootFolder.Path, bookmark.Repo)
	_, err := utils.createMissingFolder(folder)
	if err != nil {
		return "", err
	}
	filename := path.Join(folder, "README.md")
	err = utils.appendToFile(filename, appendContent)
	if err != nil {
		return "", err
	}
	return filename, nil
}
*/

