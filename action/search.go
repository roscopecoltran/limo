package action

import (
	"fmt"                                             // go-core
	"github.com/blevesearch/bleve"                    // data-idx
	"github.com/roscopecoltran/sniperkit-limo/config" // app-config
	"github.com/roscopecoltran/sniperkit-limo/model"  // data-models
	"github.com/sirupsen/logrus"                      // logs-logrus
	"github.com/spf13/cobra"                          // cli-cmd
	"strconv"                                         // go-core
	"strings"                                         // go-core
	//"github.com/davecgh/go-spew/spew" 															// debug-print
	//"github.com/k0kubun/pp" 																		// debug-print
)

/*
refs:
	links:
	- https://github.com/cioc/decentralizedSearch/blob/master/providers/stackoverflow/stackoverflow.go
	- https://github.com/yieldbot/ferret/blob/master/providers/github/github.go
	- https://github.com/piger/corpus/blob/master/file_walk.go
	- https://github.com/smnalex/stealth
	- https://github.com/keimoon/cerebro/blob/master/search/reddit.go
	- https://github.com/zjucx/SearchEngine/blob/master/main.go
	- https://github.com/google/zoekt/blob/master/cmd/zoekt-git-index/main.go
	- github.com/BenjaminCh/app-store
*/

// SearchCmd does a full-text search
var SearchCmd = &cobra.Command{
	Use:     "search <search string>",
	Aliases: []string{"find", "query", "q"},
	Short:   "Search stars",
	Long:    "Perform a full-text search on your stars",
	Example: fmt.Sprintf("  %s search robust", config.ProgramName),
	Run: func(cmd *cobra.Command, args []string) {

		output := getOutput()
		if len(args) == 0 {

			log.WithFields(logrus.Fields{
				"actions":   "SearchCmd",
				"len(args)": len(args),
			}).Warnf("You must specify a search string")

			output.Fatal("You must specify a search string")
		}

		index, err := getIndex()
		fatalOnError(err)

		query := bleve.NewMatchQuery(strings.Join(args, " "))
		request := bleve.NewSearchRequest(query)
		results, err := index.Search(request)
		fatalOnError(err)

		db, err := getDatabase()
		fatalOnError(err)

		for _, hit := range results.Hits {
			ID, err := strconv.Atoi(hit.ID)
			if err != nil {
				output.Error(err.Error())
			} else {
				star, err := model.FindStarByID(db, uint(ID))
				if err != nil {
					output.Error(err.Error())
				} else {
					output.Inline(fmt.Sprintf("(%f) ", hit.Score))
					output.StarLine(star)
				}
			}
		}
	},
}

func init() {
	log.WithFields(
		logrus.Fields{
			"prefix":      "app-action",
			"src.file":    "action/search.go",
			"cmd.name":    "SearchCmd",
			"method.name": "init()",
			"var.options": options,
		}).Info("registering command...")
	RootCmd.AddCommand(SearchCmd)
}
