package actions

import (
	"fmt"
	"os"
	// "github.com/jinzhu/configor"
	"github.com/blevesearch/bleve"
	"github.com/roscopecoltran/sniperkit-limo/config"
	"github.com/roscopecoltran/sniperkit-limo/model"
	"github.com/roscopecoltran/sniperkit-limo/output"
	"github.com/roscopecoltran/sniperkit-limo/service"
	"github.com/jinzhu/gorm"
	"github.com/boltdb/bolt"
	"github.com/spf13/cobra"
	log "github.com/sirupsen/logrus"
	// prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

var configuration *config.Config
var db *gorm.DB
var bucket *bolt.DB
var index bleve.Index

var options struct {
	config 		string
	language 	string
	output   	string
	service  	string
	tag      	string
	verbose  	bool
}

// RootCmd is the root command for limo
var RootCmd = &cobra.Command{
	Use:   "limo",
	Short: "A CLI for managing starred repositories",
	Long: `limo allows you to manage your starred repositories on GitHub, GitLab, and Bitbucket.
You can tag, display, and search your starred repositories.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		//fmt.Println(err)
		log.WithError(err).WithFields(log.Fields{"config": "Execute"}).Info("error while getting starting the program.")
		os.Exit(-1)
	}
}

func init() {

	// https://github.com/x-cray/logrus-prefixed-formatter
	// log.Formatter = new(prefixed.TextFormatter)
	// log.Level = logrus.DebugLevel

	flags := RootCmd.PersistentFlags()
	flags.StringVarP(&options.language, "language", "l", "", 								"language")
	flags.StringVarP(&options.output, 	"output", 	"o", "color", 							"output type")
	flags.StringVarP(&options.service, 	"service", 	"s", "github", 							"service")
	flags.StringVarP(&options.tag, 		"tag", 		"t", "", 								"tag")
	flags.BoolVarP(&options.verbose, 	"verbose", 	"v", false, 							"verbose output")
	flags.StringVarP(&options.config, 	"config", 	"c", "./config/settings_default.yml", 	"Path to the configuration filename")
}

func getConfiguration() (*config.Config, error) {
	if configuration == nil {
		var err error
		if configuration, err = config.ReadConfig(); err != nil {
			log.WithError(err).WithFields(log.Fields{"config": "getConfiguration"}).Info("error while getting global configuration data.")
			return nil, err
		}
	}
	return configuration, nil
}

/*
func getConfiguration2() (*config.Config, error) {
	if configuration == nil {
		var err error
		if configuration, err = config.ReadConfig(); err != nil {
			return nil, err
		}
	}
	//if err := configor.Load(&config.Config, options.config); err != nil {
	//	return nil, err
	//}
	return configuration, nil
}
*/

func getDatabase() (*gorm.DB, error) {
	if db == nil {
		cfg, err := getConfiguration()
		if err != nil {
			return nil, err
		}
		db, err = model.InitDB(cfg.DatabasePath, options.verbose)
		if err != nil {
			return nil, err
		}
	}
	return db, nil
}

func getBucket() (*bolt.DB, error) {
	if bucket == nil {
		cfg, err := getConfiguration()
		if err != nil {
			log.WithError(err).WithFields(log.Fields{"config": "getBucket"}).Info("error while getting configuration.")
			return nil, err
		}
		bucket, err = model.InitBoltDB(cfg.DatastorePath)
		if err != nil {
			log.WithError(err).WithFields(log.Fields{"config": "getBucket", "cfg.DatastorePath": cfg.DatastorePath}).Infof("error while init the BoltDB bucket at %#s", cfg.DatastorePath)
			return nil, err
		}
	}
	return bucket, nil
}

func getIndex() (bleve.Index, error) {
	if index == nil {
		cfg, err := getConfiguration()
		if err != nil {
			log.WithError(err).WithFields(log.Fields{"config": "getIndex"}).Info("error while getting configuration.")
			return nil, err
		}
		index, err = model.InitIndex(cfg.IndexPath)
		if err != nil {
			log.WithError(err).WithFields(log.Fields{"config": "getIndex", "cfg.IndexPath": cfg.IndexPath}).Infof("error while init the search engine index: %#s", cfg.IndexPath)
			return nil, err
		}
	}
	return index, nil
}

func getOutput() output.Output {
	output := output.ForName(options.output)
	oc, err := getConfiguration()
	if err == nil {
		log.WithError(err).WithFields(log.Fields{"config": "getOutput", "options.output": options.output}).Info("error while getting output options.")
		output.Configure(oc.GetOutput(options.output))
	}
	return output
}

func getService() (service.Service, error) {
	return service.ForName(options.service)
}

func checkOneStar(name string, stars []model.Star) {
	output := getOutput()
	if len(stars) == 0 {
		output.Fatal(fmt.Sprintf("No stars match '%s'", name))
	}
	if len(stars) > 1 {
		output.Error(fmt.Sprintf("Star '%s' ambiguous:\n", name))
		for _, star := range stars {
			output.StarLine(&star)
		}
		output.Fatal("Narrow your search")
	}
}

func fatalOnError(err error) {
	if err != nil {
		getOutput().Fatal(err.Error())
	}
}
