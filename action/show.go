package action

import (
	"fmt"                                             // go-core
	"github.com/roscopecoltran/sniperkit-limo/config" // app-config
	"github.com/roscopecoltran/sniperkit-limo/model"  // data-models
	"github.com/sirupsen/logrus"                      // logs-logrus
	"github.com/spf13/cobra"                          // cli-cmd
	//"github.com/davecgh/go-spew/spew" 															// debug-print
	//"github.com/k0kubun/pp" 																		// debug-print
)

// ShowCmd shows the version
var ShowCmd = &cobra.Command{
	Use:     "show <star>",
	Short:   "Show a star's details",
	Long:    "Show details about the star identified by <star>.",
	Example: fmt.Sprintf("  %s show limo", config.ProgramName),
	Run: func(cmd *cobra.Command, args []string) {
		output := getOutput()

		if len(args) == 0 {
			output.Fatal("You must specify a star")
		}

		db, err := getDatabase()
		fatalOnError(err)

		stars, err := model.FuzzyFindStarsByName(db, args[0])
		fatalOnError(err)

		for _, star := range stars {
			err = star.LoadTags(db)
			if err != nil {
				output.Error(err.Error())
			} else {
				output.Star(&star)
				output.Info("")
			}
		}
	},
}

func init() {
	log.WithFields(
		logrus.Fields{
			"prefix":      "app-action",
			"src.file":    "action/show.go",
			"cmd.name":    "ShowCmd",
			"method.name": "init()",
			"var.options": options,
		}).Info("registering command...")
	RootCmd.AddCommand(ShowCmd)
}
