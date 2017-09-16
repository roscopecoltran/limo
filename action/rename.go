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

// RenameCmd renames a tag
var RenameCmd = &cobra.Command{
	Use:     "rename <tag> <name>",
	Aliases: []string{"mv"},
	Short:   "Rename a tag",
	Long:    "Rename the tag with name <tag> to <name>.",
	Example: fmt.Sprintf("  %s rename www web", config.ProgramName),
	Run: func(cmd *cobra.Command, args []string) {
		output := getOutput()

		if len(args) < 2 {
			output.Fatal("You must specify a tag and a new name")
		}

		db, err := getDatabase()
		fatalOnError(err)

		tag, err := model.FindTagByName(db, args[0])
		fatalOnError(err)

		if tag == nil {
			output.Fatal(fmt.Sprintf("Tag '%s' not found", args[0]))
		}

		fatalOnError(tag.Rename(db, args[1]))

		output.Info(fmt.Sprintf("Renamed tag '%s' to '%s'", args[0], tag.Name))
	},
}

func init() {
	log.WithFields(
		logrus.Fields{
			"src.file": 			"action/rename.go", 
			"cmd.name": 			"RenameCmd",
			"method.name": 			"init()", 
			"var.options": 			options, 
			}).Info("registering command...")
	RootCmd.AddCommand(RenameCmd)
}
