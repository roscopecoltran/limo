package cmd

import (
	"fmt"

	"github.com/hoop33/limo/config"
	"github.com/hoop33/limo/model"
	"github.com/spf13/cobra"
)

var adders = map[string]func([]string){
	"star": addStar,
	"tag":  addTag,
	"topic": addTopic,
	"language": addLanguage,
	"academic": addAcademic,
	"readme": addReadme,
	"tree": addTree,
	"package": addPackage,
}

// AddCmd adds stars and tags
var AddCmd = &cobra.Command{
	Use:     "add <star|tag> <name>...",
	Short:   "Add star(s), tag(s), academic(s), readme(s), package(s), language(s) or topic",
	Long:    "Add star(s) or tag(s). Adding a tag adds it to your local database. Adding a star stars the repository on the specified service.",
	Example: fmt.Sprintf("  %s add tag vim database\n  %s add star hoop33/limo --service github", config.ProgramName, config.ProgramName),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			getOutput().Fatal("You must specify star or tag and values")
		}
		if fn, ok := adders[args[0]]; ok {
			fn(args[1:])
		} else {
			getOutput().Fatal(fmt.Sprintf("'%s' not valid", args[0]))
		}
	},
}

// start - custom

func addAcademic(values []string) {
	getOutput().Fatal("Not yet implemented")
}

func addPackage(values []string) {
	getOutput().Fatal("Not yet implemented")
}

func addTopic(values []string) {
	getOutput().Fatal("Not yet implemented")
}

func addReadme(values []string) {
	getOutput().Fatal("Not yet implemented")
}

func addLanguage(values []string) {
	getOutput().Fatal("Not yet implemented")
}

func addTree(values []string) {
	getOutput().Fatal("Not yet implemented")
}

// end - custom

func addStar(values []string) {
	getOutput().Fatal("Not yet implemented")
}

func addTag(values []string) {
	output := getOutput()

	db, err := getDatabase()
	fatalOnError(err)

	for _, value := range values {
		tag, created, err := model.FindOrCreateTagByName(db, value)
		if err != nil {
			output.Error(err.Error())
		} else {
			if created {
				output.Info(fmt.Sprintf("Created tag '%s'", tag.Name))
			} else {
				output.Error(fmt.Sprintf("Tag '%s' already exists", tag.Name))
			}
		}
	}
}

func init() {
	RootCmd.AddCommand(AddCmd)
}
