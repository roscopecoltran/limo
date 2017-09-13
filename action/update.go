package actions

import (
	"context"
	"fmt"
	"github.com/roscopecoltran/sniperkit-limo/config"
	"github.com/roscopecoltran/sniperkit-limo/model"
	"github.com/roscopecoltran/sniperkit-limo/service"
	"github.com/spf13/cobra"
	"github.com/sirupsen/logrus"
	// tablib "github.com/agrison/go-tablib"
	// "github.com/davecgh/go-spew/spew"
	// jsoniter "github.com/json-iterator/go"
	//"github.com/k0kubun/pp"
)

// top priority
// https://github.com/Termina1/starlight

// https://github.com/indraniel/srasearch/blob/master/makeindex/main.go
// https://github.com/ulrf/ulrf/blob/master/models/svuldb.go
// https://github.com/Jonexlee/project/blob/1d794ed0db1f47cac807381a468f1baea5e910a6/model/batch/batchInfo.go
// https://github.com/urandom/readeef

// UpdateCmd lets you log in
var UpdateCmd = &cobra.Command{
	Use:     "update",
	Short:   "Update stars from a service",
	Long:    "Update your local database with your stars from the service specified by [--service] (default: github).",
	Example: fmt.Sprintf("  %s update", config.ProgramName),
	Run: func(cmd *cobra.Command, args []string) {

		ctx := context.Background()

		// Get configuration
		cfg, err := getConfiguration()
		fatalOnError(err)
		// spew.Printf(&cfg)

		// Get the database
		db, err := getDatabase()
		fatalOnError(err)
		// spew.Printf(&db)

		// Get the database
		bucket, err := getBucket()
		fatalOnError(err)
		// spew.Printf(&bucket)

		// Just to use it once, at least, for the moment
		// we can put the config struct in the bucket
		fmt.Println(bucket)

		// Get the search index
		index, err := getIndex()
		fatalOnError(err)
		// spew.Printf(&index)

		// Get the specified service
		svc, err := getService()
		fatalOnError(err)
		// spew.Printf(&svc)

		// Get the database record for the specified service
		serviceName := service.Name(svc)
		dbSvc, _, err := model.FindOrCreateServiceByName(db, serviceName)
		fatalOnError(err)
		// spew.Printf(&dbSvc)

		// Create a channel to receive stars, since service can page
		starChan := make(chan *model.StarResult, 20)

		// Get the stars for the authenticated user
		go svc.GetStars(ctx, starChan, cfg.GetService(serviceName).Token, "", true)

		output := getOutput()

		totalCreated, totalUpdated, totalErrors := 0, 0, 0

		for starResult := range starChan {
			if starResult.Error != nil {
				totalErrors++
				output.Error(starResult.Error.Error())
			} else {
				created, err := model.CreateOrUpdateStar(db, starResult.Star, dbSvc)
				if err != nil {
					totalErrors++
					log.WithError(err).WithFields(logrus.Fields{"config": "UpdateCmd", "starResult.Star.FullName": *starResult.Star.FullName}).Warnf("error while getting creating/updating a vcs starred project. \n Error %s: %s", *starResult.Star.FullName, err.Error())
					output.Error(fmt.Sprintf("Error %s: %s", *starResult.Star.FullName, err.Error()))
				} else {
					if created {
						totalCreated++
					} else {
						totalUpdated++
					}
					err = starResult.Star.Index(index, db)
					if err != nil {
						totalErrors++
						log.WithError(err).WithFields(logrus.Fields{"config": "UpdateCmd", "starResult.Star.Index.FullName": *starResult.Star.FullName}).Warnf("error while getting creating/updating a vcs starred project. \n Error %s: %s", *starResult.Star.FullName, err.Error())
						output.Error(fmt.Sprintf("Error %s: %s", *starResult.Star.FullName, err.Error()))
					}
					// output.Tick()
				}
			}
		}
		log.WithFields(logrus.Fields{"config": "UpdateCmd", "action": "SyncedStar"}).Infof("\nCreated: %d; Synced: %d; Errors: %d", totalCreated, totalUpdated, totalErrors)
		output.Info(fmt.Sprintf("\nCreated: %d; Updated: %d; Errors: %d", totalCreated, totalUpdated, totalErrors))
	},
}

// functions to update languages, topics

func init() {
	RootCmd.AddCommand(UpdateCmd)
}
