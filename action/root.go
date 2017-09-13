package actions

import (
	"fmt"
	"os"
	"github.com/jinzhu/configor"
	"github.com/blevesearch/bleve"
	"github.com/roscopecoltran/sniperkit-limo/config"
	"github.com/roscopecoltran/sniperkit-limo/model"
	"github.com/roscopecoltran/sniperkit-limo/output"
	"github.com/roscopecoltran/sniperkit-limo/service"
	"github.com/jinzhu/gorm"
	"github.com/boltdb/bolt"
	"github.com/spf13/cobra"
	"github.com/sirupsen/logrus"
	// "github.com/davecgh/go-spew/spew"
	// "github.com/k0kubun/pp"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

var defaultConfigFilePath 	string = "~/.config/limo/limo.yaml"
var configuration 			*config.Config
var	db 						*model.DatabaseDrivers
var	log 					= logrus.New()

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
		log.WithError(err).WithFields(logrus.Fields{"config": "Execute"}).Info("error while getting starting the program.")
		os.Exit(-1)
	}
}

func init() {

	// logs
	log.Out = os.Stdout
	// log.Formatter = new(prefixed.TextFormatter)

	formatter := new(prefixed.TextFormatter)
	formatter.FullTimestamp = true

	// Set specific colors for prefix and timestamp
	formatter.SetColorScheme(&prefixed.ColorScheme{
		PrefixStyle:    "blue+b",
		TimestampStyle: "white+h",
	})

	log.Formatter = formatter

	tmpDir := config.GetTmpDir()
	log.WithFields(logrus.Fields{"action": "init", "step": "getTmpDir"}).Infof("tmp dir located at: %#s", tmpDir)

	flags := RootCmd.PersistentFlags()
	flags.StringVarP(&options.language, "language", "l", "", 								"language")
	flags.StringVarP(&options.output, 	"output", 	"o", "color", 							"output type")
	flags.StringVarP(&options.service, 	"service", 	"s", "github", 							"service")
	flags.StringVarP(&options.tag, 		"tag", 		"t", "", 								"tag")
	flags.BoolVarP(&options.verbose, 	"verbose", 	"v", false, 							"verbose output")
	flags.StringVarP(&options.config, 	"config", 	"c", "./config/settings_default.yml", 	"Path to the configuration filename")

	if options.verbose {
		log.Level = logrus.InfoLevel
	}

	db, err := model.InitDatabases()
	if err != nil {
		log.WithError(err).WithFields(
			logrus.Fields{	"action": 	"init"}).Fatal("error while loading the config files with configor package.")
	}

}

func getConfiguration() (configuration config.Config, err error) {
	configFilePath := config.FindLocalConfig()
	if configFilePath != "" {
		log.WithFields(
			logrus.Fields{	"config": 	"getConfiguration"}).Infof("FOUND configuration data to load: %#s", configFilePath)
	}
	if err := configor.Load(&configuration, configFilePath, defaultConfigFilePath); err != nil {
		log.WithError(err).WithFields(
			logrus.Fields{	"config": 	"getConfiguration"}).Fatal("error while loading the config files with configor package.")
	}
	return configuration, nil
}

func getDatabase() (*gorm.DB, error) {
	db, _ := model.GetDatabases()
	if db.gormCli == nil {
		cfg, err := getConfiguration()
		if err != nil {
			return nil, err
		}
		//db, err = model.InitDB(cfg.DatabasePath, true)
		gormCli, err := model.InitDB(cfg.DatabasePath, "sqlite3", options.verbose)
		if err != nil {
			return nil, err
		}
		db.gormCli = gormCli
	}
	return db.gormCli, nil
}

/*
func getDatabase() (*gorm.DB, error) {
	if db == nil {
		cfg, err := getConfiguration()
		if err != nil {
			return nil, err
		}
		//db, err = model.InitDB(cfg.DatabasePath, true)
		db, err = model.InitDB(cfg.DatabasePath, options.verbose)
		if err != nil {
			return nil, err
		}
	}
	return db, nil
}
*/

func getBucket() (*bolt.DB, error) {
	db, _ := model.GetDatabases()
	if db.boltCli == nil {
		cfg, err := getConfiguration()
		if err != nil {
			log.WithError(err).WithFields(
				logrus.Fields{	"config": "getBucket"}).Info("error while getting configuration.")
			return nil, err
		}
		boltCli, err := model.InitBoltDB(cfg.DatastorePath)
		if err != nil {
			log.WithError(err).WithFields(
				logrus.Fields{	"config": 			 "getBucket", 
								"cfg.DatastorePath": cfg.DatastorePath}).Warnf("error while init the BoltDB bucket at %#s", cfg.DatastorePath)
			return nil, err
		}
		db.boltCli = boltCli
	}	
	return db.boltCli, nil
}

func getIndex() (bleve.Index, error) {
	db, _ := model.GetDatabases()
	if db.bleveIdx == nil {
		cfg, err := getConfiguration()
		if err != nil {
			log.WithError(err).WithFields(
				logrus.Fields{"config": "getIndex"}).Info("error while getting configuration.")
			return nil, err
		}
		bleveIdx, err := model.InitIndex(cfg.IndexPath)
		if err != nil {
			log.WithError(err).WithFields(
				logrus.Fields{	"config": "getIndex", 
								"cfg.IndexPath": cfg.IndexPath}).Warnf("error while init the search engine index: %#s", cfg.IndexPath)
			return nil, err
		}
		db.bleveIdx = bleveIdx
	}
	return db.bleveIdx, nil
}

func getOutput() output.Output {
	output := output.ForName(options.output)
	oc, err := getConfiguration()
	if err != nil {
		log.WithError(err).WithFields(
			logrus.Fields{	"config": 		  "getOutput", 
							"options.output": options.output}).Warnf("error while getting output options.")
	} else {
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
