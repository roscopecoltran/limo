package actions

import (
	"context"
	"fmt"

	"github.com/roscopecoltran/sniperkit-limo/config"
	"github.com/roscopecoltran/sniperkit-limo/model"
	"github.com/roscopecoltran/sniperkit-limo/service"
	"github.com/spf13/cobra"
)

// DumpCmd lets you log in
var DumpCmd = &cobra.Command{
	Use:     "dump",
	Short:   "Dump stars from a service",
	Long:    "Dump your local database with your stars from the service specified by [--service] (default: github).",
	Aliases: []string{"dump", "export", "d"},
	Example: fmt.Sprintf("  %s dump", config.ProgramName),
	Run: func(cmd *cobra.Command, args []string) {
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

		totalCreated, totalDumpd, totalErrors := 0, 0, 0

		for starResult := range starChan {
			if starResult.Error != nil {
				totalErrors++
				output.Error(starResult.Error.Error())
			} else {
				created, err := model.DumpStarInfo(db, starResult.Star, dbSvc)
				if err != nil {
					totalErrors++
					output.Error(fmt.Sprintf("Error %s: %s", *starResult.Star.FullName, err.Error()))
				} else {
					if created {
						totalCreated++
					} else {
						totalDumpd++
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
		output.Info(fmt.Sprintf("\nCreated: %d; Dumped: %d; Errors: %d", totalCreated, totalDumpd, totalErrors))
	},
}

func init() {
	RootCmd.AddCommand(DumpCmd)
}
