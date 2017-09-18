package action

import (
	"context"                                          // go-core
	"fmt"                                              // go-core
	"github.com/roscopecoltran/sniperkit-limo/config"  // app-config
	"github.com/roscopecoltran/sniperkit-limo/model"   // data-models
	"github.com/roscopecoltran/sniperkit-limo/service" // svc-registry
	"github.com/sirupsen/logrus"                       // logs-logrus
	"github.com/spf13/cobra"                           // cli-cmd
	//"github.com/davecgh/go-spew/spew" 															// debug-print
	//"github.com/k0kubun/pp" 																		// debug-print
)

// SyncCmd lets you log in
var SyncCmd = &cobra.Command{
	Use:     "sync",
	Short:   "Sync stars from a service",
	Long:    "Sync your local database with your stars from the service specified by [--service] (default: github).",
	Aliases: []string{"synchronize", "sync", "update-bucket", "s"},
	Example: fmt.Sprintf("  %s sync", config.ProgramName),
	Run: func(cmd *cobra.Command, args []string) {

		ctx := context.Background()

		// Get configuration
		cfg, err := getConfiguration()
		fatalOnError(err)

		// Get the database
		db, err := getDatabase()
		fatalOnError(err)

		// Get the database
		bucket, err := getBucket()
		fatalOnError(err)

		// Just to use it once, at least, for the moment
		// we can put the config struct in the bucket
		fmt.Println(bucket)

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
		go svc.GetStars(ctx, starChan, cfg.GetService(serviceName).Token, "", true)

		output := getOutput()

		totalCreated, totalSyncd, totalErrors := 0, 0, 0

		for starResult := range starChan {
			if starResult.Error != nil {
				totalErrors++
				output.Error(starResult.Error.Error())
			} else {
				created, err := model.SyncDB(db, starResult.Star, dbSvc)
				if err != nil {
					totalErrors++
					log.WithError(err).WithFields(logrus.Fields{"config": "SyncCmd", "starResult.Star.FullName": *starResult.Star.FullName}).Warnf("error while getting creating/updating a vcs starred project. \n Error %s: %s", *starResult.Star.FullName, err.Error())
					output.Error(fmt.Sprintf("Error %s: %s", *starResult.Star.FullName, err.Error()))
				} else {
					if created {
						totalCreated++
					} else {
						totalSyncd++
					}
					err = starResult.Star.Index(index, db)
					if err != nil {
						totalErrors++
						log.WithError(err).WithFields(logrus.Fields{"config": "SyncCmd", "starResult.Star.Index.FullName": *starResult.Star.FullName}).Warnf("error while getting creating/updating a vcs starred project. \n Error %s: %s", *starResult.Star.FullName, err.Error())
						output.Error(fmt.Sprintf("Error %s: %s", *starResult.Star.FullName, err.Error()))
					}
					output.Tick()
				}
			}
		}
		log.WithFields(logrus.Fields{"config": "SyncCmd", "action": "SyncedStar"}).Infof("\nCreated: %d; Synced: %d; Errors: %d", totalCreated, totalSyncd, totalErrors)
		output.Info(fmt.Sprintf("\nCreated: %d; Synced: %d; Errors: %d", totalCreated, totalSyncd, totalErrors))
	},
}

func init() {
	log.WithFields(
		logrus.Fields{
			"prefix":      "app-action",
			"src.file":    "action/sync.go",
			"cmd.name":    "SyncCmd",
			"method.name": "init()",
			"var.options": options,
		}).Info("registering command...")
	RootCmd.AddCommand(SyncCmd)
}
