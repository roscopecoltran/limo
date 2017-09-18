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

/*
refs:
	- high_priority:
		- https://github.com/Termina1/starlight
	- links:
		- https://github.com/indraniel/srasearch/blob/master/makeindex/main.go
		- https://github.com/ulrf/ulrf/blob/master/models/svuldb.go
		- https://github.com/Jonexlee/project/blob/1d794ed0db1f47cac807381a468f1baea5e910a6/model/batch/batchInfo.go
		- https://github.com/urandom/readeef
*/

var UpdateCmd = &cobra.Command{ // UpdateCmd lets you log in
	Use:     "update",
	Short:   "Update stars from a service",
	Long:    "Update your local database with your stars from the service specified by [--service] (default: github).",
	Aliases: []string{"update", "up"},
	Example: fmt.Sprintf("  %s update", config.ProgramName),
	Run: func(cmd *cobra.Command, args []string) {

		ctx := context.Background()    // Init the context background
		cfg, err := getConfiguration() // Get configuration
		fatalOnError(err)

		db, err := getDatabase() // Get the SQL database (default: sqlite3)
		fatalOnError(err)

		// kvs, err 							:= getBucket() 										// Get the KV bucket (default: boltDB)
		// fatalOnError(err)
		// fmt.Println(kvs) 																		// Just to use it once, at least, for the moment, we can put the config struct in the bucket

		index, err := getIndex() // Get the search index
		fatalOnError(err)

		svc, err := getService() // Get the specified service
		fatalOnError(err)

		serviceName := service.Name(svc) // Get the database record for the specified service
		dbSvc, _, err := model.FindOrCreateServiceByName(db, serviceName)
		fatalOnError(err)

		starChanCount := 20 // Create a channel to receive stars, since service can page (default: 20)
		starChan := make(chan *model.StarResult, starChanCount)
		go svc.GetStars(ctx, starChan, cfg.GetService(serviceName).Token, "", true) // Get the stars for the authenticated user

		output := getOutput()                              // Get the output options
		totalCreated, totalUpdated, totalErrors := 0, 0, 0 // init new default counters

		for starResult := range starChan {
			if starResult.Error != nil {
				totalErrors++
				output.Error(starResult.Error.Error())
				log.WithError(starResult.Error).WithFields(
					logrus.Fields{
						"src.file":          "action/update.go",
						"method.name":       "UpdateCmd = &cobra.Command{...}",
						"method.prev":       "starResult.Error != nil",
						"var.options":       options,
						"var.totalCreated":  totalCreated,
						"var.totalUpdated":  totalUpdated,
						"var.totalErrors":   totalErrors,
						"var.starChanCount": starChanCount,
						"var.starResult":    *starResult.Star,
					}).Error("error occured while iterating through the starred repository channel.")
			} else {
				// need to use another method like LOAD DATA INTO FILE for MySQL for example, let's think about it and apply an efficient solution to it.
				created, err := model.CreateOrUpdateStar(db, starResult.Star, dbSvc) // create or update new starred repository
				if err != nil {
					totalErrors++
					log.WithError(err).WithFields(
						logrus.Fields{
							"src.file":                      "action/update.go",
							"method.name":                   "UpdateCmd = &cobra.Command{...}",
							"method.prev":                   "model.CreateOrUpdateStar(...)",
							"var.options":                   options,
							"var.totalCreated":              totalCreated,
							"var.totalUpdated":              totalUpdated,
							"var.totalErrors":               totalErrors,
							"var.starChanCount":             starChanCount,
							"var.starResult.Star.RemoteURI": starResult.Star.RemoteURI,
						}).Error("error while getting creating/updating the starred repo into the data-store.")
					//output.Error(fmt.Sprintf("Error %s: %s", *starResult.Star.FullName, err.Error()))
				} else {
					if created {
						totalCreated++ // increment created counter
					} else {
						totalUpdated++ // increment updated counter
					}
					err = starResult.Star.Index(index, db) // index the new content in the full-text search engine
					if err != nil {
						totalErrors++
						log.WithError(err).WithFields(
							logrus.Fields{
								"src.file":                      "action/update.go",
								"method.name":                   "UpdateCmd = &cobra.Command{...}",
								"method.prev":                   "starResult.Star.Index(...)",
								"var.options":                   options,
								"var.totalCreated":              totalCreated,
								"var.totalUpdated":              totalUpdated,
								"var.totalErrors":               totalErrors,
								"var.starChanCount":             starChanCount,
								"var.starResult.Star.RemoteURI": starResult.Star.RemoteURI,
							}).Error("error while indexing the star into the full-text engine.")
						//output.Error(fmt.Sprintf("Error %s: %s", *starResult.Star.FullName, err.Error()))
					}
					// output.Tick() 																// display tick if logrus is disabled or compatible with the tick display
				}
			}
		}
		log.WithFields(
			logrus.Fields{
				"src.file":          "action/update.go",
				"method.name":       "UpdateCmd = &cobra.Command{...}",
				"method.prev":       "starResult.Star.Index(...)",
				"var.options":       options,
				"var.totalCreated":  totalCreated,
				"var.totalUpdated":  totalUpdated,
				"var.totalErrors":   totalErrors,
				"var.starChanCount": starChanCount,
			}).Info("update completed.")
		// output.Info(fmt.Sprintf("\nCreated: %d; Updated: %d; Errors: %d", totalCreated, totalUpdated, totalErrors))
	},
}

func init() {
	log.WithFields(
		logrus.Fields{
			"prefix":      "app-action",
			"src.file":    "action/update.go",
			"cmd.name":    "UpdateCmd",
			"method.name": "init()",
			"var.options": options,
		}).Info("registering command...")
	RootCmd.AddCommand(UpdateCmd)
}
