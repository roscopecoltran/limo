package action

import (
	"fmt"																					// go-core
	"github.com/roscopecoltran/sniperkit-limo/config" 										// app-config
	"github.com/roscopecoltran/sniperkit-limo/model" 										// data-models
	"github.com/spf13/cobra" 																// cli-cmd
	"github.com/sirupsen/logrus" 															// logs-logrus
	//"github.com/davecgh/go-spew/spew" 													// debug-print
	//"github.com/k0kubun/pp" 																// debug-print
)

var homepage = false

// OpenCmd opens a star's URL in your browser
var OpenCmd = &cobra.Command{
	Use:     "open <star>",
	Short:   "Open a star's URL",
	Long:    "Open a star's URL in your default browser.",
	Example: fmt.Sprintf("  %s open limo", config.ProgramName),
	Run: func(cmd *cobra.Command, args []string) {
		output := getOutput()

		if len(args) == 0 {
			output.Fatal("You must specify a star")
		}

		db, err := getDatabase()
		fatalOnError(err)

		stars, err := model.FuzzyFindStarsByName(db, args[0])
		fatalOnError(err)

		checkOneStar(args[0], stars)

		err = stars[0].OpenInBrowser(homepage)
		fatalOnError(err)
	},
}

func init() {
	log.WithFields(
		logrus.Fields{
			"src.file": 			"action/open.go", 
			"cmd.name": 			"OpenCmd",
			"method.name": 			"init()", 
			"var.options": 			options, 
			}).Info("registering command...")
	OpenCmd.Flags().BoolVarP(&homepage, "homepage", "H", false, "open home page instead of URL")
	RootCmd.AddCommand(OpenCmd)
}
