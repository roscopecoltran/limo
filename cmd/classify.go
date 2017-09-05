package cmd

import (
	//"context"
	"fmt"
	"github.com/hoop33/limo/config"
	//"github.com/hoop33/limo/model"
	//"github.com/hoop33/limo/service"
	"github.com/spf13/cobra"
)

// ClassifyCmd lets you log in
var ClassifyCmd = &cobra.Command{
	Use:     "classify",
	Short:   "Classify stars from a service",
	Long:    "Classify your local database with your stars from the service specified by [--service] (default: github).",
	Aliases: []string{"detect", "classify", "c"},
	Example: fmt.Sprintf("  %s classify", config.ProgramName),
	Run: func(cmd *cobra.Command, args []string) {

		/*
		ctx := context.Background()

		// Get configuration
		cfg, err := getConfiguration()
		fatalOnError(err)

		// Get the database
		db, err := getDatabase()
		fatalOnError(err)

		// Get the search index
		index, err := getIndex()
		fatalOnError(err)

		// Get the specified service
		svc, err := getService()
		fatalOnError(err)

		// Get the database record for the specified service
		serviceName := service.Name(svc)
		dbSvc, _, err := model.FindOrCreateServiceByName(db, serviceName)
		fatalOnError(err)

		// Create a channel to receive stars, since service can page
		starChan := make(chan *model.StarResult, 20)

		// Get the stars for the authenticated user
		go svc.GetStars(ctx, starChan, cfg.GetService(serviceName).Token, "")

		output := getOutput()

		totalCreated, totalClassifyd, totalErrors := 0, 0, 0

		for starResult := range starChan {
			if starResult.Error != nil {
				totalErrors++
				output.Error(starResult.Error.Error())
			} else {
				created, err := model.CreateOrClassifyStar(db, starResult.Star, dbSvc)
				if err != nil {
					totalErrors++
					output.Error(fmt.Sprintf("Error %s: %s", *starResult.Star.FullName, err.Error()))
				} else {
					if created {
						totalCreated++
					} else {
						totalClassifyd++
					}
					err = starResult.Star.Index(index, db)
					if err != nil {
						totalErrors++
						output.Error(fmt.Sprintf("Error %s: %s", *starResult.Star.FullName, err.Error()))
					}
					output.Tick()
				}
			}
		}
		output.Info(fmt.Sprintf("\nCreated: %d; Classifyed: %d; Errors: %d", totalCreated, totalClassifyd, totalErrors))
		*/
		
	},
}

func init() {
	RootCmd.AddCommand(ClassifyCmd)
}
