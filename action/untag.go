package action

import (
	"fmt"																							// go-core
	"github.com/roscopecoltran/sniperkit-limo/config" 												// app-config
	"github.com/roscopecoltran/sniperkit-limo/model" 												// data-models
	"github.com/spf13/cobra" 																		// cli-cmd
	"github.com/sirupsen/logrus" 																	// logs-logrus
	//"github.com/davecgh/go-spew/spew" 															// debug-print
	//"github.com/k0kubun/pp" 																		// debug-print
)

var UntagCmd = &cobra.Command{ 																				// UntagCmd tags a star
	Use:     "untag <star> [tag]...",
	Short:   "Untag a star", 
	Long:    "Untag the star identified by <star> with the tags specified by [tag], or all if [tag] not specified.",
	Example: fmt.Sprintf("  %s untag limo gui", config.ProgramName),
	Aliases: []string{"untag", "utag", "ut"},
	Run: func(cmd *cobra.Command, args []string) {
		output := getOutput()

		if len(args) == 0 {
			output.Fatal("You must specify a star (and optionally a tag)")
		}

		db, err := getDatabase()
		fatalOnError(err)

		stars, err := model.FuzzyFindStarsByName(db, args[0])
		fatalOnError(err)

		checkOneStar(args[0], stars)

		output.StarLine(&stars[0])

		if len(args) == 1 {

			fatalOnError(stars[0].RemoveAllTags(db)) 														// Untag all

			output.Info(fmt.Sprintf("Removed all tags"))

		} else {

			fatalOnError(stars[0].LoadTags(db))

			for _, tagName := range args[1:] {

				tag, err := model.FindTagByName(db, tagName)
				if err != nil {

					output.Error(err.Error())

				} else if tag == nil {

					output.Error(fmt.Sprintf("Tag '%s' does not exist", tagName))

				} else if !stars[0].HasTag(tag) {

					output.Error(fmt.Sprintf("'%s' isn't tagged with '%s'", *stars[0].FullName, tagName))

				} else {

					err = stars[0].RemoveTag(db, tag)
					if err != nil {

						output.Error(err.Error())

					} else {

						output.Info(fmt.Sprintf("Removed tag '%s'", tag.Name))

					}

				}
			}
		}
	},
}

func init() {
	log.WithFields(
		logrus.Fields{
			"src.file": 			"action/untag.go", 
			"cmd.name": 			"UntagCmd",
			"method.name": 			"init()", 
			"var.options": 			options, 
			}).Info("registering command...")
	RootCmd.AddCommand(UntagCmd)
}
