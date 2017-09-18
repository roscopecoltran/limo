package action

import (
	"fmt"                                              // go-core
	"github.com/blevesearch/bleve"                     // search-idx
	"github.com/boltdb/bolt"                           // data-kvs
	"github.com/jinzhu/configor"                       // cfg-load
	"github.com/jinzhu/gorm"                           // data-sql
	"github.com/roscopecoltran/sniperkit-limo/config"  // app-config
	"github.com/roscopecoltran/sniperkit-limo/model"   // data-models
	"github.com/roscopecoltran/sniperkit-limo/output"  // data-output
	"github.com/roscopecoltran/sniperkit-limo/service" // svc-registry
	"os"                                               // go-core
	"runtime"                                          // go-core
	"strings"                                          // go-core
	"time"                                             // go-core
	// cfg_util "github.com/roscopecoltran/sniperkit-limo/utils/config" 							// utils-cfg
	"github.com/spf13/cobra" // cli-cmd
	"github.com/spf13/viper" // cli-cmd
	//"github.com/k0kubun/pp" 																		// debug-print
	//"github.com/davecgh/go-spew/spew" 															// debug-print
	"github.com/sirupsen/logrus"                           // logs-logrus
	prefixed "github.com/x-cray/logrus-prefixed-formatter" // logs-logrus
)

// RootCmd is the root command for limo
var RootCmd = &cobra.Command{
	Use:   fmt.Sprintf("%s", strings.ToLower(config.ProgramName)),
	Short: "An advanced CLI/Web toolkit for managing starred repositories, and help you to keep an innovative projects.",
	Long:  fmt.Sprintf("%s allows you to manage your starred repositories on GitHub, GitLab, and Bitbucket. You can tag, display, and search your starred repositories.", config.ProgramName),
	PreRun: func(cmd *cobra.Command, args []string) {
		log.Printf("Cobra.PreRun")
	},
	Run: func(cmd *cobra.Command, args []string) {
		log.Printf("Cobra.Run")
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		log.WithError(err).WithFields(
			logrus.Fields{
				"src.file":    "action/root.go",
				"prefix":      "new-instance",
				"method.name": "Execute(...)",
				"method.prev": "RootCmd.Execute(...)",
			}).Error("error while getting starting the program.")
		// os.Exit(debugExit(true))
		os.Exit(-1)
	}
}

var cfgFile string

// to remove later
var options struct {
	config   string
	language string
	output   string
	service  string
	tag      string
	verbose  bool
	dir      struct {
		tmp  string
		data string
		conf string
	}
}

var (
	dbs           = model.GetDrivers()
	configuration *config.Config // cfg-init
	log           = logrus.New() // logs-logrus
)

func init() {

	runtime.GOMAXPROCS(runtime.NumCPU())

	log.Out = os.Stdout                      // logs 	- output
	formatter := new(prefixed.TextFormatter) // logs 	- prefix-formatter
	log.Formatter = formatter                // logs 	- msg themes
	log.Level = logrus.DebugLevel            // logs 	- set the log level

	dirTmp := config.GetTmpDir() // dir 		- tmp
	options.dir.tmp = dirTmp     // options 	- dir - tmp

	flags := RootCmd.PersistentFlags() // flags
	flags.StringVarP(&options.language, "language", "l", "", "language")
	flags.StringVarP(&options.output, "output", "o", "color", "output type")
	flags.StringVarP(&options.service, "service", "s", "github", "service")
	flags.StringVarP(&options.tag, "tag", "t", "", "tag")
	flags.BoolVarP(&options.verbose, "verbose", "v", false, "verbose output")
	//flags.StringVarP(&options.config, 	"config", 	"c", "./shared/conf.d/limo.yaml", 		"Path to the configuration filename")
	RootCmd.PersistentFlags().StringVar(&options.config, "config", "", "config file (default is $HOME/.config/limo/limo.yaml)")

	if options.verbose {
		log.Level = logrus.InfoLevel
	}

	// cobra doesn't invoke these callbacks until *after* checking for the help flag, so
	// we can't use this...
	// see https://github.com/spf13/cobra/blob/6ed17b5128e8932c9ecd4c3970e8ea5e60a418ac/command.go#L590
	// cobra.OnInitialize(onInitialize)
	onInitialize()
	log.Printf("config: %v", viper.AllSettings())

	log.WithFields(
		logrus.Fields{
			"prefix":          "app-action",
			"src.file":        "action/root.go",
			"method.name":     "init()",
			"method.prev":     "config.GetTmpDir()",
			"var.options.dir": options.dir,
			"var.log.Level":   log.Level,
			"var.log":         log,
			"var.options":     options,
		}).Info("config adjusting defaults to current machine...")

	// dbs =

	// dbs 	= model.New(model.DefaultSql, model.DefaultKvs, model.DefaultGraphs)
	// pb := model.NewPocketBase(&conf.DB)

}

func New(verbose bool) error {
	//if options.config != "" { 																// if no user-defined config file to load, pick defaults filepaths
	//}
	configFiles := []string{defaultConfigFilePath}
	if _, err := config.New(configuration, true, true, true, configFiles); err != nil { // init new configuration
		log.WithError(err).WithFields(
			logrus.Fields{
				"prefix":      "app-action",
				"src.file":    "action/root.go",
				"action.type": "new-instance",
				"var.verbose": verbose,
				"var.options": options,
			}).Fatal("error while loading the config files.")
		return err
	}
	if verbose {
		log.WithFields(
			logrus.Fields{
				"prefix":            "app-action",
				"src.file":          "action/root.go",
				"action.type":       "new-instance",
				"var.verbose":       verbose,
				"var.configuration": configuration,
				"var.options":       options,
			}).Debug("error while loading the new sniperkit instance.")
	}
	return nil
}

func NewConfiguration() (*config.Config, error) {
	configFilePath := config.FindLocalConfig()
	//configuration 	:= 	&config.Config{}
	if configFilePath != "" {
		log.WithFields(logrus.Fields{
			"prefix":                    "app-action",
			"src.file":                  "action/root.go",
			"action.type":               "new-config",
			"method.name":               "NewConfiguration(...)",
			"method.prev":               "config.FindLocalConfig(...)",
			"var.configFilePath":        configFilePath,
			"var.defaultConfigFilePath": defaultConfigFilePath,
		}).Info("found local config files")
	}
	if err := configor.Load(&configuration, configFilePath, defaultConfigFilePath); err != nil {
		log.WithError(err).WithFields(logrus.Fields{
			"prefix":                    "app-action",
			"src.file":                  "action/root.go",
			"action.type":               "new-config",
			"method.name":               "NewConfiguration(...)",
			"method.prev":               "configor.Load(...)",
			"var.configFilePath":        configFilePath,
			"var.defaultConfigFilePath": defaultConfigFilePath,
		}).Error("error while loading configs with 'configor'")
	}
	return configuration, nil
}

func getConfiguration() (configuration config.Config, err error) {
	configFilePath := config.FindLocalConfig()
	if configFilePath != "" {
		log.WithFields(logrus.Fields{
			"prefix":             "app-action",
			"src.file":           "action/root.go",
			"action.type":        "get-config",
			"method.name":        "NewConfiguration(...)",
			"method.prev":        "config.FindLocalConfig(...)",
			"var.configFilePath": configFilePath,
		}).Warn("found local config file")
	}
	if err := configor.Load(&configuration, configFilePath); err != nil {
		log.WithError(err).WithFields(
			logrus.Fields{
				"prefix":             "app-action",
				"src.file":           "action/root.go",
				"action.type":        "get-config",
				"method.name":        "getConfiguration(...)",
				"method.prev":        "configor.Load(...)",
				"var.configFilePath": configFilePath,
				//"var.defaultConfigFilePath": 		defaultConfigFilePath,
			}).Fatal("error while loading configs with 'configor'")
	}
	return configuration, nil
}

func getDatabase() (*gorm.DB, error) {
	cfg, err := getConfiguration()
	if err != nil {
		log.WithError(err).WithFields(
			logrus.Fields{
				"prefix":      "app-action",
				"src.file":    "action/root.go",
				"action.type": "get-index",
				"method.prev": "getIndex(...)",
				"method.name": "getConfiguration(...)",
				"var.cfg":     cfg,
			}).Error("error while getting configuration.")
		return nil, err
	}
	//pp.Print(cfg)
	if !dbs.Initialized {
		if !dbs.Gorm.Ok {
			if err := dbs.Init("sqlite3", map[string]bool{"boltdb": true}, map[string]bool{}); err != nil {
				log.WithError(err).WithFields(
					logrus.Fields{
						"prefix":      "app-action",
						"src.file":    "action/root.go",
						"action.type": "dbs-init",
						"method.prev": "dbs.Init(...)",
						"method.name": "getDatabase(...)",
						"var.cfg":     cfg,
						"var.dbs":     dbs,
					}).Error("error while initializing the BoltDB bucket.")
			}
		}
	}
	//pp.Print(dbs)
	return dbs.Gorm.Cli, nil
}

func getBucket() (*bolt.DB, error) {
	cfg, err := getConfiguration()
	if err != nil {
		log.WithError(err).WithFields(
			logrus.Fields{
				"prefix":      "app-action",
				"src.file":    "action/root.go",
				"action.type": "get-index",
				"method.name": "getIndex(...)",
				"method.prev": "getConfiguration(...)",
				"var.cfg":     cfg,
			}).Error("error while getting configuration.")
		return nil, err
	}
	if !dbs.Bolt.Ok {
		bucket, err := model.InitBoltDB(cfg.DatastorePath)
		if err != nil {
			log.WithError(err).WithFields(
				logrus.Fields{
					"prefix":                "app-action",
					"src.file":              "action/root.go",
					"action.type":           "get-bucket",
					"method.prev":           "model.InitBoltDB(...)",
					"method.name":           "getBucket(...)",
					"var.cfg.DatastorePath": cfg.DatastorePath,
				}).Error("error while initializing the BoltDB bucket.")
			return bucket, err
		}
		dbs.Bolt.Ok = true
		dbs.Bolt.Cli = bucket
		return dbs.Bolt.Cli, nil
	}
	return dbs.Bolt.Cli, nil
}

func getIndex() (bleve.Index, error) {
	cfg, err := getConfiguration()
	if err != nil {
		log.WithError(err).WithFields(
			logrus.Fields{
				"prefix":      "app-action",
				"src.file":    "action/root.go",
				"action.type": "get-index",
				"method.name": "getIndex(...)",
				"method.prev": "getConfiguration(...)",
				"var.cfg":     cfg,
			}).Error("error while getting configuration.")
		return nil, err
	}
	if !dbs.Bleve.Ok {
		bleveIdx, err := model.InitIndex(cfg.IndexPath)
		if err != nil {
			log.WithError(err).WithFields(
				logrus.Fields{
					"prefix":            "app-action",
					"src.file":          "action/root.go",
					"action.type":       "get-index",
					"method.prev":       "model.InitIndex(...)",
					"method.name":       "getIndex(...)",
					"var.cfg.IndexPath": cfg.IndexPath,
				}).Error("error while initializing the BoltDB bucket.")
			return bleveIdx, err
		}
		dbs.Bleve.Ok = true
		dbs.Bleve.Cli = bleveIdx
		return dbs.Bleve.Cli, nil
	}
	return dbs.Bleve.Cli, nil
}

func getOutput() output.Output {
	output := output.ForName(options.output)
	oCfg, err := getConfiguration()
	if err != nil {
		log.WithError(err).WithFields(
			logrus.Fields{
				"prefix":      "app-action",
				"src.file":    "action/root.go",
				"action.type": "get-bucket",
				"method.name": "getOutput(...)",
				"method.prev": "getConfiguration(...)",
				"var.output":  output,
				"var.oCfg":    oCfg,
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
	output := getOutput()
	starCount := len(stars)
	if starCount == 0 {
		log.WithFields(
			logrus.Fields{
				"prefix":      "app-action",
				"src.file":    "action/root.go",
				"action.type": "check-star-count",
				"method.name": "checkOneStar(...)",
				"method.prev": "starCount == 0",
				"var.name":    fmt.Sprintf("%s", name),
			}).Error("No stars match")
	}
	if starCount > 1 {
		for _, star := range stars {
			output.StarLine(&star)
		}
		log.WithFields(
			logrus.Fields{
				"prefix":        "app-action",
				"src.file":      "action/root.go",
				"action.type":   "check-star-count",
				"method.name":   "checkOneStar(...)",
				"method.prev":   "starCount > 1",
				"var.name":      fmt.Sprintf("%s", name),
				"var.starCount": starCount,
			}).Errorf("Star '%s' name is ambiguous.", name)
	}
}

func fatalOnError(err error) {
	if err != nil {
		log.WithError(err).WithFields(
			logrus.Fields{
				"prefix":      "app-action",
				"src.file":    "action/root.go",
				"action.type": "get-bucket",
				"method.name": "fatalOnError(...)",
				"var.err":     err.Error(),
			}).Fatal("fatal error triggered.")
	}
}

func onInitialize() {
	initializeConfig()
	initializeRuntime()
}

func initializeConfig() {
	// cfg_util.InitConfig(cfgFile)
}

func initializeRuntime() {
	/*
		go func() {
			gRPC := server.NewKedsRPCServer()
			gRPC.Cobra = server.NewCobra(RootCmd)
			gRPC.Start()
		}()
	*/
	//TODO this feels hacky...need a more reliable way to determine that the plugins have loaded
	//and the server has started
	time.Sleep(3 * time.Second)
}
