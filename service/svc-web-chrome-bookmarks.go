package service

// # Page generation config
// chrome_bookmarks_path="/Users/$USER/Library/Application Support/Google/Chrome/Default/Bookmarks"
// output_file_name="index.html"
// root_folder="Supermarks"

// # Upload config
// scp_destination="matt@man1.biz:www/man1.biz/supermarks/"

// ref. https://raw.githubusercontent.com/man1/Supermarks/master/supermarks.go

import (
	"encoding/json"
	"flag"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"time"
	"github.com/sirupsen/logrus"
)

// MARK: Constants
const ChromePrimaryTemplateName 			= "main"
const ChromeTemplateFile 					= "templates.html"

// MARK: Config defaults and global state
const DefaultChromeBookmarksFile 			= "/Users/$USER/Library/Application Support/Google/Chrome/Default/Bookmarks"
const DefaultChromeBookmarksOutputFile 		= "index.html"
const DefaultChromeBookmarksRootFolderName 	= "Bookmarks Bar"

var ChromeGlobalConfig ChromeBookmarkConfig

type ChromeBookmarksConfig struct {
  ChromeBookmarksFile string
  OutputFile string
  RootFolderName string
}

// MARK: Structs to represent page & bookmark data
type ChromePageData struct {
  RootNodes []ChromeBookmarkNode
  Updated time.Time
}

type ChromeBookmarkNode struct {
  Title, URL string
  Children []ChromeBookmarkNode
}

// MARK: Generic types for JSON parsing, for easy reference
type ChromeJSON map[string]interface{}
type ChromeJSONArr []interface{}

// MARK: The meat & potatoes
// Retrieve JSON for desired bookmarks from the filesystem
func getChromeJSON() ChromeJSON {
  bookmarksFile, readError := ioutil.ReadFile(ChromeGlobalConfig.ChromeBookmarksFile)
  check(readError)
  var bookmarksJSON ChromeJSON
  unmarshalError := json.Unmarshal(bookmarksFile, &bookmarksJSON)
  check(unmarshalError)
  return bookmarksJSON
}

// Transform bookmark JSON into a PageData struct
func pageDataFromChromeJSON(data ChromeJSON) ChromePageData {
  nodes := bookmarkChromeNodesFromJSON(data, ChromeGlobalConfig.RootFolderName, false)
  return ChromePageData{nodes, time.Now()}
}

// Parse JSON into recursively defined bookmark nodes. If `rootFound` is
// false, this will traverse down to the folder with name `rootFolder`
// before starting to build the BookmarkNode struct.
func bookmarkChromeNodesFromJSON(data ChromeJSON, rootFolder string, rootFound bool) []ChromeBookmarkNode {
  rootFoundOrIsRoot := rootFound || isBookmarkChromeWithName(data, rootFolder)
  var name string
  var URL string
  children := []ChromeBookmarkNode{}
  for key, val := range data {
    switch key {
    case "roots", "bookmark_bar":
      // Dive down through Chrome's root bookmark file node(s)
      valJSON, isJSON := val.(map[string]interface{})
      if isJSON {
        return bookmarkChromeNodesFromJSON(JSON(valJSON), rootFolder, rootFound)
      } else {
        log.Fatal("Failed to parse root node JSON")
      }
    case "name":
      // Set the name for this node (applies for both folders & link bookmarks)
      name = val.(string)
    case "url":
      // Set the URL for this node, if exists (applies only for link bookmarks)
      URL = val.(string)
    case "children":
      // Recursively parse the array of JSON children for this folder node
      valJSONArr, isJSONArr := val.([]interface{})
      if isJSONArr {
        nodeArr := ChromeJSONArr(valJSONArr)
        for i := range nodeArr {
          node := nodeArr[i]
          nodeJSON, isJSON := node.(map[string]interface{})
          if isJSON {
            newNodes := bookmarkChromeNodesFromJSON(ChromeJSON(nodeJSON), rootFolder, rootFoundOrIsRoot)
            children  = append(children, newNodes...)
          } else {
            log.WithFields(logrus.Fields{  
                        "file":         "service/svc-web-chrome.go", 
                        "method_name":  "bookmarkChromeNodesFromJSON(data ChromeJSON, rootFolder string, rootFound bool)", 
                        "driver":       "chrome", 
                        "feature":      "bookmarks",
                        "action":       "ChromeJSONArr(valJSONArr)",
                      }).Fatal("Failed to parse child node JSON")
           // log.Fatal("Failed to parse child node JSON")
          }
        }
      } else {
        log.WithFields(logrus.Fields{  
                    "file":         "service/svc-web-chrome.go", 
                    "method_name":  "bookmarkChromeNodesFromJSON(data ChromeJSON, rootFolder string, rootFound bool)", 
                    "driver":       "chrome", 
                    "feature":      "bookmarks",
                    "action":       "ChromeJSONArr(valJSONArr)",
                  }).Fatal("Failed to parse child array JSON")
        //log.Fatal("Failed to parse child array JSON")
      }
    }
  }
  // If the root folder has been found previously, return this node. Otherwise,
  // only return its children, such that traversal to root folder will continue
  // without yet starting to build the BookmarkNode struct for the page.
  if rootFound {
    return []ChromeBookmarkNode{ChromeBookmarkNode{name, URL, children}}
  } else {
    return children
  }
}

// Identify if the name of the top level JSON node is equal to `name`
func isBookmarkChromeWithName(data ChromeJSON, name string) bool {
  for key, val := range data {
    switch key {
    case "name":
      return val == name
    }
  }
  return false
}

// Write out a file generated from template.html using the provided PageData
func generateChromeBookmarkPage(pageContents ChromePageData) {
  // Setup the template
  pageTemplate, templateCreationError := template.ParseFiles(ChromeTemplateFile)
  check(templateCreationError)
  // Write out the templated HTML
  file, fileError 	:= os.Create(ChromeGlobalConfig.OutputFile)
  check(fileError)
  defer file.Close()
  templateUseErr 	:= pageTemplate.ExecuteTemplate(file, ChromePrimaryTemplateName, pageContents)
  check(templateUseErr)
}

// Check an error
func check(err error) {
  if err != nil {
    log.Fatal(err)
  }
}

// MARK: Main
func mainChromeBookmark() {
	// Parse flags into ChromeGlobalConfig
	bookmarksFilePtr 			:= flag.String("bookmarks", DefaultChromeBookmarksFile, 			"Path to Chrome's bookmarks file")
	outputFilePtr 				:= flag.String("output", 	DefaultChromeBookmarksOutputFile, 		"Name for output file, e.g. 'bookmarks.html'")
	rootFolderPtr 				:= flag.String("root", 		DefaultChromeBookmarksRootFolderName, 	"Name of the root bookmark folder to parse")
	flag.Parse()
	ChromeGlobalConfig 	 		= ChromeBookmarksConfig{*bookmarksFilePtr, *outputFilePtr, *rootFolderPtr}
	// Generate the page
	bookmarks 					:= getChromeJSON()
	pd 							:= pageDataFromChromeJSON(bookmarks)
	generateChromeBookmarkPage(pd)
}