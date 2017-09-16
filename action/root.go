package action

import (
	"fmt"																							// go-core
	"os"																							// go-core
	"github.com/jinzhu/configor" 																	// cfg-load
	"github.com/jinzhu/gorm" 																		// data-sql
	"github.com/boltdb/bolt" 																		// data-kvs
	"github.com/blevesearch/bleve" 																	// search-idx
	"github.com/roscopecoltran/sniperkit-limo/config" 												// app-config
	"github.com/roscopecoltran/sniperkit-limo/model" 												// data-models
	"github.com/roscopecoltran/sniperkit-limo/service" 												// svc-registry
	"github.com/roscopecoltran/sniperkit-limo/output" 												// data-output
	"github.com/spf13/cobra" 																		// cli-cmd
	"github.com/k0kubun/pp" 																		// debug-print
	//"github.com/davecgh/go-spew/spew" 															// debug-print
	"github.com/sirupsen/logrus" 																	// logs-logrus
	prefixed "github.com/x-cray/logrus-prefixed-formatter"  										// logs-logrus
)

// RootCmd is the root command for limo
var RootCmd = &cobra.Command{
	Use:   fmt.Sprintf("%s", strings.ToLower(config.ProgramName)),
	Short: "An advanced CLI/Web toolkit for managing starred repositories, and help you to keep an innovative projects.",
	Long: fmt.Sprintf("%s allows you to manage your starred repositories on GitHub, GitLab, and Bitbucket. You can tag, display, and search your starred repositories.", config.ProgramName),
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		log.WithError(err).WithFields(
			logrus.Fields{	
				"src.file": 					"action/root.go",
				"prefix": 						"new-instance",
				"method.name": 					"Execute(...)",
				"method.prev": 					"RootCmd.Execute(...)",
				}).Error("error while getting starting the program.")
		os.Exit(-1)
	}
}

func init() {

	log.Out 				= os.Stdout 													// logs 	- output
	formatter 				:= new(prefixed.TextFormatter) 									// logs 	- prefix-formatter
	log.Formatter 			= formatter 													// logs 	- msg themes
	log.Level 				= logrus.DebugLevel 											// logs 	- set the log level

	dirTmp 					:= 	config.GetTmpDir() 											// dir 		- tmp
	options.dir.tmp 		= 	dirTmp 														// options 	- dir - tmp

	flags := RootCmd.PersistentFlags() 														// flags
	flags.StringVarP(&options.language, "language", "l", "", 								"language")
	flags.StringVarP(&options.output, 	"output", 	"o", "color", 							"output type")
	flags.StringVarP(&options.service, 	"service", 	"s", "github", 							"service")
	flags.StringVarP(&options.tag, 		"tag", 		"t", "", 								"tag")
	flags.BoolVarP(&options.verbose, 	"verbose", 	"v", false, 							"verbose output")
	flags.StringVarP(&options.config, 	"config", 	"c", "./shared/conf.d/limo.yaml", 		"Path to the configuration filename")

	if options.verbose {
		log.Level = logrus.InfoLevel
	}

	log.WithFields(
		logrus.Fields{
			"src.file": 				"action/root.go", 
			"method.name": 				"init()", 
			"method.prev": 				"config.GetTmpDir()",
			"var.options.dir": 			options.dir, 
			"var.log.Level": 			log.Level, 
			"var.log": 					log, 
			"var.options": 				options, 
			}).Info("config adjusting defaults to current machine...")

}																																	

func New(verbose bool) (*config.Config, *model.DatabaseDrivers, error) {
	if options.config != "" { 																// if no user-defined config file to load, pick defaults filepaths

	}
	configFiles 	:= []string{defaultConfigFilePath}
	if _, err 	:= config.New(configuration, true, true, true, configFiles); err != nil { 	// init new configuration
		log.WithError(err).WithFields(
			logrus.Fields{	
				"src.file": 					"action/root.go",
				"action.type": 					"new-instance",
				"var.verbose": 					verbose,
				"var.options": 					options,
				}).Fatal("error while loading the config files.")
		return configuration, dbs, err
	}
	// init new data clients
	if _, err 		:= dbs.New(true, true); err != nil {
		log.WithError(err).WithFields(
			logrus.Fields{	
				"src.file": 					"action/root.go",
				"action.type": 					"new-instance",
				"var.verbose": 					verbose,
				"var.configuration": 			configuration,
				"var.options": 					options,
				}).Fatal("error while loading the databases drivers.")
		return configuration, dbs, err
	}
	if verbose {
		log.WithFields(
			logrus.Fields{	
				"src.file": 					"action/root.go",
				"action.type": 					"new-instance",
				"var.verbose": 					verbose,
				"var.configuration": 			configuration,
				"var.dbs": 						dbs,
				"var.options": 					options,
				}).Debug("error while loading the new sniperkit instance.")
	}
	return configuration, dbs, nil
}

func NewConfiguration() (*config.Config, error) {
	configFilePath 		:= 	config.FindLocalConfig()
	//configuration 	:= 	&config.Config{}
	if configFilePath != "" {
		log.WithFields(logrus.Fields{	
			"src.file": 						"action/root.go",
			"action.type": 						"new-config",
			"method.name": 						"NewConfiguration(...)",
			"method.prev": 						"config.FindLocalConfig(...)",
			"var.configFilePath": 				configFilePath,
			"var.defaultConfigFilePath": 		defaultConfigFilePath,
			}).Info("found local config files")
	}
	if err := configor.Load(&configuration, configFilePath, defaultConfigFilePath); err != nil {
		log.WithError(err).WithFields(logrus.Fields{	
			"src.file": 						"action/root.go",
			"action.type": 						"new-config",
			"method.name": 						"NewConfiguration(...)",
			"method.prev": 						"configor.Load(...)",
			"var.configFilePath": 				configFilePath,
			"var.defaultConfigFilePath": 		defaultConfigFilePath,
			}).Error("error while loading configs with 'configor'")
	}
	return configuration, nil
}

func getConfiguration() (configuration config.Config, err error) {
	configFilePath 	  := config.FindLocalConfig()
	if configFilePath != "" {
		log.WithFields(logrus.Fields{	
			"src.file": 						"action/root.go",
			"action.type": 						"get-config",
			"method.name": 						"NewConfiguration(...)",
			"method.prev": 						"config.FindLocalConfig(...)",
			"var.configFilePath": 				configFilePath,
			}).Error("found local config file")
	}
	if err := configor.Load(&configuration, configFilePath, defaultConfigFilePath); err != nil {
		log.WithError(err).WithFields(
			logrus.Fields{	
				"src.file": 						"action/root.go",
				"action.type": 						"get-config",
				"method.name": 						"getConfiguration(...)",
				"method.prev": 						"configor.Load(...)",
				"var.configFilePath": 				configFilePath,
				}).Fatal("error while loading configs with 'configor'")
	}
	return configuration, nil
}

func getDatabase() (*gorm.DB, error) {
	if dbs.gormCli 	== nil {
		gormCli, err 	:= 	dbs.New(true, true)
		if err != nil {
			log.WithError(err).WithFields(
				logrus.Fields{	
					"src.file": 					"action/root.go",
					"action.type": 					"get-db-sql",
					"method.name": 					"getDatabase(...)",
					"method.prev": 					"model.GetDatabases(...)",
					"var.dbs": 						dbs,
					}).Error("error while connecting to all active database drivers.")
			return gormCli, err
		}
		dbs.gormCli 	= gormCli
		pp.Println(dbs)
		return gormCli, nil
	}
	return dbs.gormCli, nil
}

func getBucket() (*bolt.DB, error) {
	if dbs.boltCli 	== nil {
		cfg, err 	:= getConfiguration()
		if err != nil {
			log.WithError(err).WithFields(
				logrus.Fields{	
					"src.file": 					"action/root.go",
					"action.type": 					"get-bucket",
					"method.prev": 					"getConfiguration(...)",
					"method.name": 					"getBucket(...)",
					"var.cfg": 						cfg,
					}).Error("error while getting configuration.")
			return nil, err
		}
		boltCli, err := model.InitBoltDB(cfg.DatastorePath)
		if err != nil {
			log.WithError(err).WithFields(
				logrus.Fields{	
					"src.file": 					"action/root.go",
					"action.type": 					"get-bucket",
					"method.prev": 					"model.InitBoltDB(...)",
					"method.name": 					"getBucket(...)",
					"var.cfg.DatastorePath": 		cfg.DatastorePath,
					}).Error("error while initializing the BoltDB bucket.") 
			return nil, err
		}
		dbs.boltCli = boltCli
		return boltCli, nil
	}	
	return dbs.boltCli, nil
}

func getIndex() (bleve.Index, error) {
	if dbs.bleveIdx == nil {
		cfg, err := getConfiguration()
		if err != nil {
			log.WithError(err).WithFields(
				logrus.Fields{	
					"src.file": 					"action/root.go",
					"action.type": 					"get-index",
					"method.name": 					"getIndex(...)",
					"method.prev": 					"getConfiguration(...)",
					"var.cfg": 						cfg,
					}).Error("error while getting configuration.") 
			return nil, err
		}
		bleveIdx, err := model.InitIndex(cfg.IndexPath)
		if err != nil {
			log.WithError(err).WithFields(
				logrus.Fields{	
					"src.file": 					"action/root.go",
					"action.type": 					"get-index",
					"method.prev": 					"model.InitIndex(...)",
					"method.name": 					"getIndex(...)",
					"var.cfg.IndexPath": 			cfg.IndexPath,
					}).Error("error while initializing the BoltDB bucket.") 
			return nil, err
		}
		dbs.bleveIdx = bleveIdx
		return bleveIdx, nil
	}
	return dbs.bleveIdx, nil
}

func getOutput() output.Output {
	output 		:= output.ForName(options.output)
	oCfg, err 	:= getConfiguration()
	if err != nil {
		log.WithError(err).WithFields(
			logrus.Fields{	
				"src.file": 						"action/root.go",
				"action.type": 						"get-bucket",
				"method.name": 						"getOutput(...)",
				"method.prev": 						"getConfiguration(...)",
				"var.output": 						output,
				"var.oCfg": 						oCfg,
				}).Error("error while initializing the BoltDB bucket.") 
	} else {
		output.Configure(oCfg.GetOutput(options.output))
	}
	return output
}

func getService() (service.Service, error) {
	return service.ForName(options.service)
}

func checkOneStar(name string, stars []model.Star) {
	output 			:= getOutput()
	starCount 		:= len(stars)
	if starCount == 0 {
		log.WithError(err).WithFields(
			logrus.Fields{	
				"src.file": 						"action/root.go",
				"action.type": 						"check-star-count",
				"method.name": 						"checkOneStar(...)",
				"method.prev": 						"starCount == 0",
				"var.name": 						fmt.Sprintf("%s", name),				
				}).Error("No stars match") 
	}
	if starCount > 1 {
		for _, star := range stars {
			output.StarLine(&star)
		}
		log.WithError(err).WithFields(
			logrus.Fields{	
				"src.file": 						"action/root.go",
				"action.type": 						"check-star-count",
				"method.name": 						"checkOneStar(...)",
				"method.prev": 						"starCount > 1",
				"var.name": 						fmt.Sprintf("%s", name),				
				"var.starCount": 					starCount,				
				}).Errorf("Star '%s' name is ambiguous.", name) 
	}
}

func fatalOnError(err error) {
	if err != nil {
		log.WithError(err).WithFields(
			logrus.Fields{	
				"src.file": 						"action/root.go",
				"action.type": 						"get-bucket",
				"method.name": 						"fatalOnError(...)",
				"var.err": 							err.Error(),				
				}).Fatal("fatal error triggered.") 
	}
}
