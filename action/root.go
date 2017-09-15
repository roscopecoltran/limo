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


const defaultConfigFilePath 	= 	"~/.config/limo/limo.yaml"

var (
	configuration 				*config.Config 														// cfg-init
	db 							= &model.DatabaseDrivers{} 											// data-drivers
	log 						= logrus.New() 														// logs-logrus
)

var (
	inputFile, outputFile 		string 																// ai-word-embed
	dimension, window     		int 																// ai-word-embed
	learningRate          		float64 															// ai-word-embed
)

//type queueDrivers struct {
	//MQ 				map[string]*nsq.Producer
	//NOSQL   			map[string]*gorm.DB
	//SQL   			map[string]*gorm.DB
	//IDX 				map[string]*bleve.Index
	//KVS 				map[string]etcd.KeysAPI
//}

var queue struct {
	Disabled 	bool 		
	// QueueDrivers                 queueDrivers	
}

// https://github.com/toorop/tmail/blob/master/core/scope.go

var options struct {
	config 		string
	language 	string
	output   	string
	service  	string
	tag      	string
	verbose  	bool
	dir 		struct {
		tmp 	string 
		data 	string
		conf 	string
	}
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

	log.Out = os.Stdout 																	// logs
	formatter := new(prefixed.TextFormatter)
	formatter.FullTimestamp = true
	formatter.SetColorScheme(&prefixed.ColorScheme{ 										// Set specific colors for prefix and timestamp
		PrefixStyle:    "blue+b",
		TimestampStyle: "white+h",
	})
	log.Formatter = formatter

	options.dir.tmp 		:= config.GetTmpDir()

	log.WithFields(
		logrus.Fields{
			"src.file": 			"action/root.go", 
			"method.name": 			"init()", 
			"method.prev": 			"config.GetTmpDir()",
			"var.options.dir": 		options.dir, 
			"var.log": 				log, 
			}).Info("config adjusting defaults to current machine...")

	flags := RootCmd.PersistentFlags()
	flags.StringVarP(&options.language, "language", "l", "", 								"language")
	flags.StringVarP(&options.output, 	"output", 	"o", "color", 							"output type")
	flags.StringVarP(&options.service, 	"service", 	"s", "github", 							"service")
	flags.StringVarP(&options.tag, 		"tag", 		"t", "", 								"tag")
	flags.BoolVarP(&options.verbose, 	"verbose", 	"v", false, 							"verbose output")
	flags.StringVarP(&options.config, 	"config", 	"c", "./shared/conf.d/limo.yaml", 		"Path to the configuration filename")

	if options.verbose {
		log.Level = logrus.InfoLevel
	}
}																																	

/*

// Machine Learning

// GetCommonFlagSet sets the common flags for models.
func GetCommonFlagSet() *flag.FlagSet {
	fs := flag.NewFlagSet(RootCmd.Name(), flag.ContinueOnError)
	fs.StringVarP(&inputFile, "input", "i", "example/input.txt", "Input file path for learning")
	fs.StringVarP(&outputFile, "output", "o", "example/word_vectors.txt", "Output file path for each learned word vector")
	fs.IntVarP(&dimension, "dimension", "d", 10, "Set word vector dimension size")
	fs.IntVarP(&window, "window", "w", 5, "Set window size")
	fs.Float64Var(&learningRate, "lr", 0.025, "Set init learning rate")
	return fs
}

// NewCommon creates the common struct.
func NewCommon() models.Common {
	return models.Common{
		InputFile:    inputFile,
		OutputFile:   outputFile,
		Dimension:    dimension,
		Window:       window,
		LearningRate: learningRate,
	}
}
*/

func New(verbose bool) (*config.Config, *model.DatabaseDrivers, error) {
	// if no user-defined config file to load, pick defaults filepaths
	if options.config != "" {

	}
	configFiles 	:= []string{defaultConfigFilePath}
	// init new configuration
	if _, err 	:= config.New(configuration, true, true, true, configFiles); err != nil {
		log.WithError(err).WithFields(
			logrus.Fields{	
				"src.file": 		"action/root.go",
				"action.type": 		"load-db",
				"var.verbose": 		verbose,
				}).Fatal("error while loading the config files.")
		return configuration, db, err
	}
	// init new data clients
	if _, err 		:= db.New(true, true); err != nil {
		log.WithError(err).WithFields(
			logrus.Fields{	
				"src.file": 		"action/root.go",
				"action.type": 		"load-db",
				"var.verbose": 		verbose,
				}).Fatal("error while loading the databases drivers.")
		return configuration, db, err
	}
	if verbose {
		log.WithFields(
			logrus.Fields{	
				"src.file": 		"action/root.go",
				"action.type": 		"new-instance",
				"var.cfg": 			cfg,
				"var.db": 			db,
				"var.verbose": 		verbose,
				}).Debug("error while loading the new sniperkit instance.")
	}
	return configuration, db, nil
}

func NewConfiguration() (*config.Config, error) {
	configFilePath 	:= 	config.FindLocalConfig()
	configuration 	:= 	&config.Config{}
	if configFilePath != "" {
		log.WithFields(logrus.Fields{	
			"src.file": 					"action/root.go",
			"action.type": 					"new-config",
			"method.name": 					"NewConfiguration(...)",
			"method.prev": 					"config.FindLocalConfig(...)",
			"var.configFilePath": 			configFilePath,
			"var.defaultConfigFilePath": 	defaultConfigFilePath,
			}).Info("found local config files")
	}
	if err := configor.Load(&configuration, configFilePath, defaultConfigFilePath); err != nil {
		log.WithError(err).WithFields(logrus.Fields{	
			"src.file": 					"action/root.go",
			"action.type": 					"new-config",
			"method.name": 					"NewConfiguration(...)",
			"method.prev": 					"configor.Load(...)",
			"var.configFilePath": 			configFilePath,
			"var.defaultConfigFilePath": 	defaultConfigFilePath,
			}).Error("error while loading configs with 'configor'")
			//}).Fatal("error while loading the config files with configor package.")
	}
	return configuration, nil
}

func getConfiguration() (configuration config.Config, err error) {
	configFilePath := config.FindLocalConfig()
	if configFilePath != "" {
		log.WithFields(logrus.Fields{	
			"src.file": 					"action/root.go",
			"action.type": 					"new-config",
			"method.name": 					"NewConfiguration(...)",
			"method.prev": 					"config.FindLocalConfig(...)",
			"var.configFilePath": 			configFilePath,
			}).Info("found local config file")
	}
	if err := configor.Load(&configuration, configFilePath, defaultConfigFilePath); err != nil {
		log.WithError(err).WithFields(logrus.Fields{	
			"src.file": 					"action/root.go",
			"action.type": 					"new-config",
			"method.name": 					"getConfiguration(...)",
			"method.prev": 					"configor.Load(...)",
			"var.configFilePath": 			configFilePath,
			}).Fatal("error while loading configs with 'configor'")
	}
	return configuration, nil
}

func getDatabase() (*gorm.DB, error) {
	//db 				:= 	&model.DatabaseDrivers{}
	if _, err 		:= 	db.New(true, true); err != nil {
		log.WithError(err).WithFields(
			logrus.Fields{	
				"src.file": 					"action/root.go",
				"action.type": 					"get-db",
				"method.name": 					"getDatabase(...)",
				"method.prev": 					"model.GetDatabases(...)",
				"var.db": 						db,
				}).Fatal("error while connecting to all active database drivers.")
		return db.gormCli, err
	}
	pp.Println(db)
	//db.gormCli = gormCli
	/*
	if db.gormCli == nil {
		cfg, err 		:= getConfiguration()
		if err != nil {
			return nil, err
		}
		//db, err = model.InitDB(cfg.DatabasePath, true)
		gormCli, err 	:= model.InitGorm(cfg.DatabasePath, "sqlite3", options.verbose)
		if err != nil {
			return nil, err
		}
		db.gormCli = gormCli
	}
	*/
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
