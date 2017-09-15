package model

// https://github.com/ssut/pocketnpm/blob/master/db/gorm_backend.go
// 

import (
    "path/filepath" 												// go-core
    "errors" 														// go-core
	"time" 															// go-core
	"path" 															// go-core
	"os" 															// go-core
	"github.com/allegro/bigcache" 									// data-cache-bigcache
	"github.com/jinzhu/gorm" 										// db-sql-gorm
	_ "github.com/jinzhu/gorm/dialects/sqlite" 						// db-sql-gorm-sqlite3
	_ "github.com/jinzhu/gorm/dialects/mysql" 						// db-sql-gorm-mysql
	_ "github.com/jinzhu/gorm/dialects/postgres" 					// db-sql-gorm-postgres
	_ "github.com/jinzhu/gorm/dialects/postgres" 					// db-sql-gorm-postgres
	// "gopkg.in/mgo.v2" 											// db-nosql-mongodb
	// "gopkg.in/mgo.v2/bson" 										// db-nosql-mongodb
	etcd "github.com/coreos/etcd/clientv3" 							// db-kvs-etcd
	//etcd "github.com/coreos/etcd/client" 							// db-kvs-etcd
	"github.com/boltdb/bolt" 										// db-kvs-boltdb
	"github.com/garyburd/redigo/redis" 								// db-kvs-redis
	"github.com/jmcvetta/neoism" 									// db-graph-neo4j
	// "github.com/cayleygraph/cayley" 								// db-graph-cayley
	// "github.com/cayleygraph/cayley/graph" 						// db-graph-cayley
	// _ "github.com/cayleygraph/cayley/graph/bolt" 				// db-graph-cayley
	// "github.com/cayleygraph/cayley/quad" 						// db-graph-cayley
	"github.com/ckaznocha/taggraph" 								// db-graph-taggraph
	"github.com/blevesearch/bleve" 									// data-index-search
	tablib "github.com/agrison/go-tablib" 							// data-processing-tablib
	// jsoniter "github.com/json-iterator/go" 						// data-processing-jsoniter
	// "github.com/davecgh/go-spew/spew" 							// data-debug
	// "github.com/astaxie/beego" 									// web-cms
	"golang.org/x/net/context" 										// web-context
    "github.com/qor/qor" 											// web-qor-admin-ui
    "github.com/qor/admin" 											// web-qor-admin-ui
	"github.com/sirupsen/logrus"									// logs-logrus
	prefixed "github.com/x-cray/logrus-prefixed-formatter" 			// logs-logrus
	"github.com/roscopecoltran/sniperkit-limo/config" 				// go-core
)

var validDataOutput 	= []string{"md","csv","yaml","json","xlsx","xml","tsv","mysql","postgres","html","ascii"} // datasets - formats
var availableLocales 	= []string{"en-US", "fr-FR", "pl-PL"}
// var serviceConfig config.Config
// var cfg *config.Config

type EnhancedTime 		time.Time

/*
type Databases struct {
	Datastore 			map[string]*bolt.DB
	Database   			map[string]*gorm.DB
	SearchIdx 			map[string]*bleve.Index
	KvIdx 				map[string]etcd.KeysAPI
}
*/
// var databases = make(map[string]Service)

var (
	Tables       	= 	[]interface{}{
		&Service{}, 	&Category{}, 																				// service + registry organization
		&Star{}, 		&Readme{}, 		&WikiPage{}, 	&User{},													// vcs content indexation
		&Tag{}, 		&Topic{}, 		&Tree{}, 		&Language{}, 	&LanguageDetected{}, 	&LanguageType{}, 	// vcs repository classification
	}
)

var (
	adminUI 			*admin.Admin
	configuration 		*config.Config 		//
	dbs 				*DatabaseDrivers 	// *DatabaseDrivers //New(true, true)
	log 				= logrus.New()
	tagg 				= taggraph.NewTagGaph()
	keyValMap 			map[string]string
)

// ref. https://github.com/tinrab/go-mmo/blob/master/db/dbobjects_gen.go
//type Database interface {
//	Dial(cfg *Config) error
//	Close()
//}

//	globalSetting := make(map[string]string)
// https://github.com/thesyncim/365a/blob/master/server/app.go
// https://github.com/emotionaldots/arbitrage/blob/master/cmd/arbitrage-db/main.go

func init() {
	log.Out 		= os.Stdout 					// logs
	formatter 		:= new(prefixed.TextFormatter) 	// logs
	log.Formatter 	= formatter 					// logs
	log.Level 		= logrus.DebugLevel 			// logs
}

// New return a new db with config input
//func NewDatabases(conf config.Config) (db *Databases, err error) {
//}

//func GetDatabases(cfg *config.Config, verbose bool, debug bool) (db *DatabaseDrivers, err error) {
func GetDatabases() (*DatabaseDrivers, error) {
	dbs 						:= 	&DatabaseDrivers{}
	var err error
	if dbs, err 	= New(true, true); err != nil {
		log.WithError(err).WithFields(
			logrus.Fields{	"src.file": 		"model/data-db-connector.go", 
							"prefix": 			"dbs-get-all",
							"method.name": 		"GetDatabases(...)", 
							"method.prev": 		"New(...)",
							"db.driver": 		"all", 
							"db.driver.groups": "sql,nosql,kvs", 
							}).Error("error while trying to init all database drivers.")
		return dbs, err
	}
	return dbs, nil
}

//func New(cfg *config.Config, verbose bool, debug bool) (db *DatabaseDrivers, err error) {
func New(verbose bool, debug bool) (*DatabaseDrivers, error) {
	dbs 						:= 	&DatabaseDrivers{}
	adapter := "sqlite3"
	gormCli, err := InitGorm(adapter, "./shared/data/limo/limo.db", true)
	if err != nil {
		log.WithError(err).WithFields(
			logrus.Fields{	"src.file": 		"model/data-db-connector.go", 
							"method.name": 		"New(...)", 
							"db.type": 			"sql", 
							"db.driver": 		"gorm", 
							"db.adpater": 		adapter,
							"method.prev": 		"db.initGorm(...)",
							"prefix": 			"dbs-new",
							"action": 			"InitGorm",
							}).Error("error while trying to init 'Gorm' database driver.")
		return dbs, err
	}
	if err := AutoLoadGorm(gormCli, true, true, true, Tables...); err != nil {
		log.WithError(err).WithFields(
			logrus.Fields{	"src.file": 		"model/data-db-connector.go", 
							"method.name": 		"New(...)", 
							"db.type": 			"sql", 
							"db.driver": 		"gorm", 
							"db.adpater": 		adapter,
							"method.prev": 		"db.autoLoadGorm(...)",
							"prefix": 			"dbs-new",
							"action": 			"InitGorm",
							}).Error("error while trying to init 'Gorm' database driver.")
		return dbs, err
	}
	/*
	if cfg.Database.Seeds.AutoLoad {
		// cfg.Database.Seeds.PrefixPath
		// cfg.Database.Seeds.Format
		filepaths, _ := filepath.Glob("db/seeds/data/*.yml")
		if err := configor.Load(&Seeds, filepaths...); err != nil {
			panic(err)
		}
	}
	*/
	boltCli, err := InitBoltDB("./shared/data/limo/limo.boltdb")
	if err != nil {
		log.WithError(err).WithFields(
			logrus.Fields{	"src.file": 		"model/data-db-connector.go", 
							"method.name": 		"New(...)", 
							"db.type": 			"kvs", 
							"db.driver": 		"bolt", 
							"db.adapter": 		"boltdb", 
							"method.prev": 		"db.initBoltDB(...)",
							"prefix": 			"dbs-new",
							"action": 			"InitBoltDB",
							}).Error("error while trying to init 'BoltDB' database driver")
		return dbs, err
	}
	etcdDefaultHost 	:= []string{"http://127.0.0.1:2379"}
	etcdDefaultTimeout 	:= 1 * time.Second
	etcdCli, err := InitEtcd(etcdDefaultHost, etcdDefaultTimeout, true)
	if err != nil {
		log.WithError(err).WithFields(
			logrus.Fields{	"src.file": 		"model/data-db-connector.go", 
							"prefix": 			"dbs-new",
							"db.type": 			"kvs", 
							"db.driver": 		"etcd", 
							"db.adapter": 		"etcd", 
							"method.prev": 		"db.initEtcd(...)",
							"action": 			"AutoloadDB",
							}).Error("error while trying to auto-load all program the tables")
		return dbs, err
	}

	// graphql
	// client := graphql.NewClient("https://example.com/graphql", nil, nil)

	return &DatabaseDrivers{
			gormCli: 	gormCli,
			boltCli:  	boltCli,
			etcdCli: 	etcdCli,
		}, nil

}

func InitDatabases() (*DatabaseDrivers, error) {

	db 				:= &DatabaseDrivers{}
	adapter 		:= "sqlite3"

	if err := db.initGorm(adapter, "./shared/data/limo/limo.db", true); err != nil {
		log.WithError(err).WithFields(
			logrus.Fields{	"src.file": 		"model/data-db-connector.go", 
							"prefix": 			"dbs-init",
							"method.name": 		"InitDatabases(...)", 
							"method.prev": 		"db.initGorm(...)",
							"db.adapter": 		adapter, 
							"db.type": 			"sql",
							"db.driver": 		"gorm", 
							}).Error("error while trying to init 'Gorm' database driver.")
		return db, err
	}

	if err := db.autoLoadGorm(true, true, true, Tables...); err != nil {
		log.WithError(err).WithFields(
			logrus.Fields{	"src.file": 		"model/data-db-connector.go", 
							"prefix": 			"dbs-init",
							"method.name": 		"Init", 
							"method.prev": 		"db.autoLoadGorm(...)",
							"db.adapter": 		adapter, 
							"db.type": 			"sql",
							"db.driver": 		"gorm", 
							}).Error("error while trying to init 'Gorm' database driver.")
		return db, err

	}

	/*
	if cfg.Database.Seeds.AutoLoad {
		// cfg.Database.Seeds.PrefixPath
		// cfg.Database.Seeds.Format
		filepaths, _ := filepath.Glob("db/seeds/data/*.yml")
		if err := configor.Load(&Seeds, filepaths...); err != nil {
			panic(err)
		}
	}
	*/

	if err := db.initBoltDB("./shared/data/limo/limo.boltdb", 0600, &bolt.Options{Timeout: 1 * time.Second}, true); err != nil {
		log.WithError(err).WithFields(
			logrus.Fields{	"file": 			"model/data-db-connector.go", 
							"method.name": 		"Init", 
							"db.adapter": 		"boltdb", 
							"method.prev": 		"db.initBoltDB(...)",
							"prefix": 			"dbs-init",
							"action": 			"InitBoltDB",
							}).Error("error while trying to init 'BoltDB' database driver")
		return db, err
	}

	etcdDefaultHost 	:= []string{"http://127.0.0.1:2379"}
	etcdDefaultTimeout 	:= 1 * time.Second
	if err := db.initEtcd(etcdDefaultHost, etcdDefaultTimeout, true); err != nil {
		log.WithError(err).WithFields(
			logrus.Fields{	"file": 			"model/data-db-connector.go", 
							"method.name": 		"Init", 
							"db.adapter": 		"etcd", 
							"method.prev": 		"db.initEtcd(...)",
							"prefix": 			"dbs-init",							
							"action": 			"AutoloadDB",
							}).Error("error while trying to auto-load all program the tables")
		return db, err
	}
	return db, nil
}

// https://github.com/qor/qor-example/blob/master/db/db.go
// InitDB initializes the database at the specified path
func InitGorm(filepath string, adapter string, verbose bool) (*gorm.DB, error) {
//func InitDB(filepath string, adapter string, verbose bool) (*gorm.DB, error) {
	gormDB, err := gorm.Open(adapter, filepath) 	// Get more config options to setup the SQL database
	if err != nil {
		log.WithError(err).WithFields(
			logrus.Fields{	"prefix": 			"dbs-init",
							"db.adapter": 		adapter,
							"src.file": 		"model/data-db-connector.go", 
							"method.name": 		"InitGorm(...)", 							
							"method.prev": 		"gorm.Open(...)",
							}).Warn("error while init the database with gorm.")
		return nil, err
	}
	gormDB.LogMode(verbose) 	// cfg.App.DebugMode
	/*
	// isAutoMigrate, isTruncate, isAdminUIResource
	if err := gormDB.AutoloadDB(true, true, false, Tables...); err != nil {
		log.WithError(err).WithFields(
			logrus.Fields{	"file": 		"model/data-db-connector.go", 
							"method.name": 	"InitGorm", 
							"adapter": 		adapter, 
							"action": 		"AutoloadDB",
							}).Warn("error while trying to auto-load all program the tables")
	}
	*/
	return gormDB, nil
}

func AutoLoadGorm(clientSQL *gorm.DB, isAutoMigrate bool, isTruncate bool, isAdminUI bool, tables ...interface{}) (error) {
	if isAdminUI {
		adminUI 	= 	admin.New(&qor.Config{DB: clientSQL})
	}
	for _, table := range tables {
		if isTruncate {
			if err := clientSQL.DropTableIfExists(table).Error; err != nil {
				log.WithError(err).WithFields(
					logrus.Fields{	"src.file": 			"model/data-db-connector.go", 
									"prefix": 				"db-gorm",
									"method.name": 			"(db *DatabaseDrivers) autoLoadGorm(...)", 
									"method.prev": 			"db.gormCli.DropTableIfExists(...)",
									"var.db.isTruncate": 	isTruncate,
									"var.db.table": 		table,
									}).Warn("error while trying to drop an existing SQL table")
				return err
			}
		}
		if isAutoMigrate {
			if err := clientSQL.AutoMigrate(table).Error; err != nil {
				log.WithError(err).WithFields(
					logrus.Fields{	"src.file": 				"model/data-db-connector.go", 
									"prefix": 					"db-gorm",
									"method.name": 				"(db *DatabaseDrivers) autoLoadGorm(...)", 
									"method.prev": 				"db.gormCli.AutoMigrate(...)",
									"var.db.isAutoMigrate": 	isAutoMigrate,
									"var.db.table": 			table,
									}).Warn("error while trying to auto-migrate db table")
				return err
			}
		}
		if isAdminUI {
			adminUI.AddResource(table)
			log.WithFields(
				logrus.Fields{	"src.file": 			"model/data-db-connector.go", 
								"method.name": 			"(db *DatabaseDrivers) autoLoadGorm(...)", 
								"method.prev": 			"adminUI.AddResource(...)",
								"var.adminui.status": 	isAdminUI,
								"var.adminui.table": 	table,
								}).Info("adding admin UI resource for the table")
		}
	}
	if isAdminUI {
		if len(adminUI.GetResources()) > 0 {
			for _, resource := range adminUI.GetResources() {	
				log.WithFields(
					logrus.Fields{	"src.file": 				"model/data-db-connector.go", 
									"method.name": 				"(db *DatabaseDrivers) autoLoadGorm(...)", 
									"method.prev": 				"adminUI.GetResources()",
									"prefix": 					"webui-admin",
									"var.adminui.resource": 	resource,
									}).Info("detected new admin UI resource")
			}		
		}
	}
	return nil
}

func InitBoltDB(filePath string) (*bolt.DB, error) {
	// Get more config options to setup the bucket or the queue of tasks
	boltDB, err := bolt.Open(filepath, 0600, &bolt.Options{Timeout: time.Second})
	if err != nil {
		log.WithError(err).WithFields(
			logrus.Fields{	"prefix": 				"db-boltdb",
							"method.name": 			"InitBoltDB(...)", 
							"method.prev": 			"bolt.Open(...)",
							"db.adapter": 			"boltdb", 
							"src.file": 			"model/data-db-connector.go", 
							"var.bolt.filepath": 	filepath, 
							"var.bolt.options": 	&bolt.Options{Timeout: time.Second},
							}).Warn("error while init the database with boltDB.")
		return nil, err
	}
	return boltDB, err
}

func InitEtcd(hosts []string, timeOut time.Duration, verbose bool) (etcd.KeysAPI, error) {
	ectdConfig := etcd.Config{
		Endpoints:               hosts,
		Transport:               etcd.DefaultTransport,
		HeaderTimeoutPerRequest: timeOut,
	}
	cli, err := etcd.New(ectdConfig)
	if err != nil {
		log.WithError(err).WithFields(
			logrus.Fields{	"prefix": 				"kvs-etcd",
							"method.name": 			"InitEtcd(...)", 
							"method.prev": 			"etcd.New(...)",
							"db.adapter": 			"etcd", 
							"var.etcd.cfg": 		cfg,
							"src.file": 			"model/data-db-connector.go", 
							}).Warn("error while init the client connection with Etcd Key/Value store.")
		return err
	}
	etcdClient := etcd.NewKeysAPI(cli)
	_, err = etcdClient.Get(context.Background(), "/foo", nil)
	if err != nil && err.Error() == etcd.ErrClusterUnavailable.Error() {
		log.WithError(err).WithFields(
			logrus.Fields{	"method.name": 			"InitEtcd(...)", 
							"db.adapter": 			"etcd", 
							"prefix": 				"kvs-etcd",
							"method.prev": 			"etcdClient.Get(...)",
							"var.etcd.cfg": 		cfg,
							"msg.error":  			etcd.ErrClusterUnavailable.Error(),
							"src.file": 			"model/data-db-connector.go", 
							}).Warn("error while init the client connection with Etcd Key/Value store.")
		return err
	}
	return nil
}

/*
// ref. https://github.com/ssut/pocketnpm/blob/master/db/database.go
func InitBigCache(filepath string) (*bolt.DB, error) {
	cacheConfig := bigcache.DefaultConfig(time.Duration(config.CacheLifetime) * time.Minute)
	cacheConfig.MaxEntrySize = 8192
	cacheConfig.HardMaxCacheSize = config.MaxCacheSize
	cache, err := bigcache.NewBigCache(cacheConfig)
	if err != nil {
		log.Fatalf("Failed to initialize in-memory cache: %v", err)
	}
}
*/

type DatabaseDrivers struct {
	gormCli 			*gorm.DB
	//dynamodbClient 	*dynamodb.DynamoDB
	boltCli  			*bolt.DB
	bleveIdx 			*bleve.Index
	neo4jCli 			*neoism.Database
	cache 				*bigcache.BigCache
	etcdCli   			*EtcdClientPool
	redisCli 			redis.Conn
}

func (dbs *DatabaseDrivers) New(verbose bool, debug bool) (*DatabaseDrivers, error) {
	adapter := "sqlite3"
	if err := dbs.initGorm(adapter, "./shared/data/limo/limo.db", true); err != nil {
		log.WithError(err).WithFields(
			logrus.Fields{	"src.file": 		"model/data-db-connector.go", 
							"prefix": 			"dbs-new",
							"db.adapater": 		adapter, 
							"db.type": 			"sql", 
							"db.driver": 		"gorm", 
							"method.name": 		"(db *DatabaseDrivers) New(...)", 
							"method.prev": 		"db.initGorm(...)",
							"var.verbose": 		verbose,
							"var.debug": 		debug,
							}).Error("error while trying to init 'Gorm' database driver.")
		return dbs, err
	}
	if err := dbs.autoLoadGorm(true, true, true, Tables...); err != nil {
		log.WithError(err).WithFields(
			logrus.Fields{	"src.file": 		"model/data-db-connector.go", 
							"db.adapater": 		adapter, 
							"db.type": 			"sql", 
							"db.driver": 		"gorm", 
							"method.name": 		"(db *DatabaseDrivers) New(...)", 
							"method.prev": 		"db.autoLoadGorm(...)",
							"prefix": 			"dbs-new",
							"var.verbose": 		verbose,
							"var.debug": 		debug,
							}).Error("error while trying to init 'Gorm' database driver.")
		return dbs, err
	}
	/*
	if cfg.Database.Seeds.AutoLoad {
		// cfg.Database.Seeds.PrefixPath
		// cfg.Database.Seeds.Format
		filepaths, _ := filepath.Glob("db/seeds/data/*.yml")
		if err := configor.Load(&Seeds, filepaths...); err != nil {
			panic(err)
		}
	}
	*/
	if err := dbs.initBoltDB("./shared/data/limo/limo.boltdb", 0600, &bolt.Options{Timeout: 1 * time.Second}, true); err != nil {
		log.WithError(err).WithFields(
			logrus.Fields{	"prefix": 						"dbs-new",
							"src.file": 					"model/data-db-connector.go", 
							"db.adapter": 					"bolt",
							"db.driver": 					"boltdb", 
							"db.type": 						"kvs",
							"method.name": 					"(db *DatabaseDrivers) New(...)", 
							"method.prev": 					"db.initBoltDB(...)",
							"var.bolt.Options.Timeout": 	1 * time.Second,
							"var.bolt.file.permissions": 	0600,
							"var.bolt.file.prefix.path": 	"./shared/data/limo/limo.boltdb",
							}).Error("error while trying to init 'BoltDB' database driver")
		return dbs, err
	}
	etcdDefaultHost 	:= []string{"http://127.0.0.1:2379"}
	etcdDefaultTimeout 	:= 1 * time.Second
	if err := dbs.initEtcd(etcdDefaultHost, etcdDefaultTimeout, true); err != nil {
		log.WithError(err).WithFields(
			logrus.Fields{	"src.file": 		"model/data-db-connector.go", 
							"prefix": 			"dbs-new",
							"method.name": 		"(db *DatabaseDrivers) New(...)", 
							"method.prev": 		"db.initEtcd(...)",
							"db.type": 			"kvs", 
							"db.driver": 		"etcd",
							"db.adapter": 		"etcd",
							}).Error("error while trying to auto-load all program the tables")
		return dbs, err
	}
	return dbs, nil
}

func (dbs *DatabaseDrivers) Close() {
	// others (map ?!)
    dbs.gormCli.Close()
}

//func initGorm(db *gorm.DB) {
func (dbs *DatabaseDrivers) initGorm(filepath string, adapter string, verbose bool) (error) {
//(db *DatabaseDrivers) func initGorm(filepath string, adapter string, verbose bool) (error) {
	gormDB, err := gorm.Open(adapter, filepath) 	// Get more config options to setup the SQL database
	if err != nil {
		log.WithError(err).WithFields(
			logrus.Fields{	"src.file": 			"model/data-db-connector.go", 
							"prefix": 				"db-gorm",
							"method.name": 			"(db *DatabaseDrivers) initGorm(...)", 
							"method.prev": 			"gorm.Open(...)",
							"db.adapter": 			adapter, 
							"var.db.verbose": 		verbose,
							}).Warn("error while init the database with gorm.")
		return err
	}
	gormDB.LogMode(verbose) 	// cfg.App.DebugMode
	// isAutoMigrate, isTruncate, isAdminUIResource
	if err := dbs.autoLoadGorm(true, true, false, Tables...); err != nil {
		log.WithError(err).WithFields(
			logrus.Fields{	"src.file": 			"model/data-db-connector.go", 
							"prefix": 				"db-gorm",
							"method.name": 			"(db *DatabaseDrivers) initGorm(...)", 
							"method.prev": 			"db.autoLoadGorm(...)",
							"db.adapter": 			adapter, 
							"var.db.verbose": 		verbose,
							}).Warn("error while trying to auto-load all program the tables")
		return err
	}
	dbs.gormCli = gormDB
	return nil
}

func (dbs *DatabaseDrivers) getGorm() *gorm.DB {
    return dbs.gormCli
}


func (dbs *DatabaseDrivers) closeGorm() {
    dbs.gormCli.Close()
}

func (dbs *DatabaseDrivers) autoLoadGorm(isAutoMigrate bool, isTruncate bool, isAdminUI bool, tables ...interface{}) (error) {
//(db *DatabaseDrivers) func autoLoadGorm(isAutoMigrate bool, isTruncate bool, isAdminUIResource bool, tables ...interface{}) (error) {
	if isAdminUI {
		adminUI 	= 	admin.New(&qor.Config{DB: dbs.gormCli})
	}
	for _, table := range tables {
		if isTruncate {
			if err := dbs.gormCli.DropTableIfExists(table).Error; err != nil {
				log.WithError(err).WithFields(
					logrus.Fields{	"src.file": 			"model/data-db-connector.go", 
									"prefix": 				"db-gorm",
									"method.name": 			"(db *DatabaseDrivers) autoLoadGorm(...)", 
									"method.prev": 			"db.gormCli.DropTableIfExists(...)",
									"var.db.isTruncate": 	isTruncate,
									"var.db.table": 		table,
									}).Warn("error while trying to drop an existing SQL table")
				return err
			}
		}
		if isAutoMigrate {
			if err := dbs.gormCli.AutoMigrate(table).Error; err != nil {
				log.WithError(err).WithFields(
					logrus.Fields{	"src.file": 				"model/data-db-connector.go", 
									"prefix": 					"db-gorm",
									"method.name": 				"(db *DatabaseDrivers) autoLoadGorm(...)", 
									"method.prev": 				"db.gormCli.AutoMigrate(...)",
									"var.db.isAutoMigrate": 	isAutoMigrate,
									"var.db.table": 			table,
									}).Warn("error while trying to auto-migrate db table")
				return err
			}
		}
		if isAdminUI {
			adminUI.AddResource(table)
			log.WithFields(
				logrus.Fields{	"src.file": 			"model/data-db-connector.go", 
								"method.name": 			"(db *DatabaseDrivers) autoLoadGorm(...)", 
								"method.prev": 			"adminUI.AddResource(...)",
								"var.adminui.status": 	isAdminUI,
								"var.adminui.table": 	table,
								}).Info("adding admin UI resource for the table")
		}
	}
	if isAdminUI {
		if len(adminUI.GetResources()) > 0 {
			for _, resource := range adminUI.GetResources() {	
				log.WithFields(
					logrus.Fields{	"src.file": 				"model/data-db-connector.go", 
									"method.name": 				"(db *DatabaseDrivers) autoLoadGorm(...)", 
									"method.prev": 				"adminUI.GetResources()",
									"prefix": 					"webui-admin",
									"var.adminui.resource": 	resource,
									}).Info("detected new admin UI resource")
			}		
		}
	}
	return nil
}

func (dbs *DatabaseDrivers) initBoltDB(filePath string, filePermissions os.FileMode, options *bolt.Options, verbose bool) (error) {
//(db *DatabaseDrivers) func initBoltDB(filePath string, filePermissions string, options *bolt.Options, verbose bool) (error) {
	boltDB, err := bolt.Open(filePath, filePermissions, options)
	if err != nil {
		log.WithError(err).WithFields(
			logrus.Fields{	"method.name": 			"(db *DatabaseDrivers) initBoltDB(...)", 
							"src.file": 			"model/data-db-connector.go",
							"db.engine": 			"boltdb", 
							"var.boltdb.filepath": 	filePath, 
							"var.boltdb.options":  	options,
							"method.prev": 			"bolt.Open(...)",
							}).Warnf("error while init the database with boltDB.")
		return err
	}
	dbs.boltCli = boltDB
	return nil
}

func (dbs *DatabaseDrivers) getBoltDB() *bolt.DB {
    return dbs.boltCli
}

func (dbs *DatabaseDrivers) closeBoltDB() {
    dbs.boltCli.Close()
}

func (dbs *DatabaseDrivers) initRedis(Password string, DbNum int) {
    var err error
    redisCli, err := redis.Dial("tcp", ":6379")
    if err != nil {
        //log.Println("failed to create the client", err)
		log.WithError(err).WithFields(
			logrus.Fields{	"src.file": 		"model/data-db-connector.go", 
							"method.name": 		"(db *DatabaseDrivers) initRedis(...)", 
							"driver": 			"redigo", 
							"adapter": 			"redis", 
							"prefix": 			"dbs-redis",
							"method.prev": 		"redis.Dial(...)",
							}).Errorln("failed to create the client", err)
        return
    }
	dbs.redisCli = redisCli
	/*
    var err2 error
    _, err2 = db.redisCli.Client.Do("SELECT", DbNum)
    if err2 != nil {
		log.WithError(err).WithFields(
			logrus.Fields{	"src.file": 		"model/data-db-connector.go", 
							"method.name": 		"(db *DatabaseDrivers) initRedis(...)", 
							"adapter": 			"redis", 
							"driver": 			"redigo", 
							"prefix": 			"dbs-redis",
							"method.prev": 		" db.redisCli.Client.Do(...)",
							}).Errorln("failed to create the client", err2)
        // log.Println("failed to create the client", err2)
    }
    */
}

func (dbs *DatabaseDrivers) getRedis() redis.Conn {
    return dbs.redisCli
}

func (dbs *DatabaseDrivers) closeRedis() {
    dbs.redisCli.Close()
}

type EtcdCli struct {
	EtcdCli 				*etcd.Client
	inUse   				bool
}

type EtcdClientPool struct {
	connPool 				chan 			EtcdCli
	Address  				string
}

const (
	ETCDCLIENTTIMEOUT      = 30 * time.Second
	ETCD_CLIENT_POOL_COUNT = 10
)

func (dbs *DatabaseDrivers) initEtcd(etcdAP, hosts []string, timeout time.Duration, verbose bool) error {
	// If ETCD is not provided, we assme that we are in single instance mode.
	// We keep our client data in a map.
	if etcdAP == "" {
		keyValMap = make(map[string]string)
		etcdCli = nil
		return
	}
	etcdCliLocal := EtcdClientPool{Address: etcdAP, connPool: make(chan EtcdCli, ETCD_CLIENT_POOL_COUNT)}
	for i := 0; i < ETCD_CLIENT_POOL_COUNT; i++ {
		var etcdCli EtcdCli
		if err := etcdCli.initEtcdCli(etcdAP); err == nil {
			etcdCliLocal.connPool <- etcdCli
		}
	}
	dbs.etcdCli = &etcdCliLocal

}

func RegisterCallBk(etcdAP, path string, callbk func(interface{}, bool, string, string), cbkctx interface{}) {
	if etcdAP != "" {
		go EtcdV3Monitor(etcdAP, path, callbk, cbkctx)
	}
}

func (cliPool *EtcdClientPool) getCli() EtcdCli {
	cli := <-cliPool.connPool
	return cli
}

func (cliPool *EtcdClientPool) putCli(cli EtcdCli) {
	if cli.EtcdCli == nil {
		log.Crit(ctx, "Releasing invalid client", nil)
	}
	cliPool.connPool <- cli
}

func (cliPool *EtcdClientPool) reconnectCli(cli EtcdCli) error {
	if cli.EtcdCli != nil {
		cli.EtcdCli.Close()
	}
	return cli.init(cliPool.Address)
}

func (cli *EtcdCli) initEtcdCli(etcdAddress string) error {
	dialTimeout 	:= time.Duration(ETCDCLIENTTIMEOUT)
	localCli, err 	:= etcd.New(etcd.Config{
		Endpoints:   []string{etcdAddress},
		DialTimeout: dialTimeout,
	})
	if err != nil {
		log.WithError(err).WithFields(
			logrus.Fields{	"method.name": 			"(dbs *DatabaseDrivers) init(...)", 
							"db.client.engine": 	"etcd", 
							"db.client.version": 	"v3", 
							"var.cfg": 		   		cfg,
							"var.etcdAddress": 		etcdAddress,
							"var.dialTimeout": 		dialTimeout,
							}).Warnf("Failed to connect to ETCD")
		return err
	}
	dbs.EtcdCli 	= localCli
	dbs.inUse 		= false
	return nil
}

func GetKey(key string) string {
	// t1 := time.Now()
	// defer fmt.Println("Get Key Time taken", time.Since(t1))
	if etcdCli == nil {
		if val, found := keyValMap[key]; found == true {
			return val
		}
	} else {
		if val := EtcdV3Get(key); len(val) != 0 {
			return val[0]
		}
	}
	return ""
}

func SetKeyVal(key, val string, leasetime int) {
	// t1 := time.Now()
	// defer fmt.Println("Set Key Time taken", time.Since(t1))
	if etcdCli == nil {
		keyValMap[key] = val
	} else {
		EtcdV3Set(key, val, int64(leasetime))
	}
}

func DeleteKey(key string) {
	if etcdCli != nil {
		EtcdV3Del(key)
	} else {
		delete(keyValMap, key)
	}
}

func EtcdV3Set(key, value string, timeout int64) {
	var err error
	etcli 	:= etcdCli.getCli()
	defer etcdCli.putCli(etcli)
	cli 	:= etcli.EtcdCli
	if timeout != 0 {
		// minimum lease TTL is 5-second
		ectx, cancel := context.WithTimeout(context.Background(), ETCDCLIENTTIMEOUT)
		resp, err1 := cli.Grant(ectx, timeout)
		if err1 != nil {
			log.WithError(err1).WithFields(
				logrus.Fields{	"method.name": 			"(dbs *DatabaseDrivers) init(...)", 
								"db.client.engine": 	"etcd", 
								"db.client.version": 	"v3", 
								//"var.cfg": 		   		cfg,
								"var.key": 				key,
								"var.val": 				value,
								}).Error("ETCD grant set Failed")
			// log.Err(ctx, "ETCD grant set Failed", blog.Fields{"key": key, "val": value})
			return
		}
		_, err = cli.Put(context.Background(), key, value, clientv3.WithLease(resp.ID))
		cancel()
	} else {
		_, err = cli.Put(context.Background(), key, value)
		// cancel()
	}
	if err != nil {
		log.WithError(err).WithFields(
			logrus.Fields{	"method.name": 			"(dbs *DatabaseDrivers) init(...)", 
							"db.client.engine": 	"etcd", 
							"db.client.version": 	"v3", 
							//"var.cfg": 		   		cfg,
							"var.key": 				key,
							"var.val": 				value,
							}).Error("ETCD value set Failed")
		// log.Err(ctx, "ETCD value set Failed", blog.Fields{"key": key, "val": value, "err": err.Error()})
	}
}

func EtcdV3Get(key string) []string {
	etcli := etcdCli.getCli()
	defer etcdCli.putCli(etcli)
	cli := etcli.EtcdCli
	requestTimeout := time.Duration(ETCDCLIENTTIMEOUT)
	lctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	resp, err := cli.Get(lctx, key)
	cancel()
	if err != nil {
		log.WithError(err).WithFields(
			logrus.Fields{	"method.name": 			"(dbs *DatabaseDrivers) init(...)", 
							"db.client.engine": 	"etcd", 
							"db.client.version": 	"v3", 
							"var.key": 				key,
							}).Error("ETCD value get Failed")
		// log.Err(ctx, "ETCD value get Failed", blog.Fields{"key": key, "err": err.Error()})
	}
	var ret []string
	for _, ev := range resp.Kvs {
		ret = append(ret, string(ev.Value))
	}
	return ret
}

func EtcdV3Del(key string) {
	etcli := etcdCli.getCli()
	defer etcdCli.putCli(etcli)
	cli := etcli.EtcdCli
	// delete the keys
	requestTimeout := time.Duration(ETCDCLIENTTIMEOUT)
	lctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	_, err := cli.Delete(lctx, key, clientv3.WithPrefix())
	cancel()
	if err != nil {
		log.WithError(err).WithFields(
			logrus.Fields{	"method.name": 			"(dbs *DatabaseDrivers) init(...)", 
							"db.client.engine": 	"etcd", 
							"db.client.version": 	"v3", 
							"var.key": 				key,
							}).Error("ETCD value delete Failed")
		// log.Err(ctx, "ETCD value delete Failed", blog.Fields{"key": key, "err": err.Error()})
	}
}

func EtcdV3Monitor(etcdAP, keytoWatch string, callbk func(interface{}, bool, string, string), cbkctx interface{}) {
	dialTimeout := time.Duration(ETCDCLIENTTIMEOUT)
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{etcdAP},
		DialTimeout: dialTimeout,
	})
	//We cannot make any progress if client connection fails.
	if err != nil {
		log.WithError(err).WithFields(
			logrus.Fields{	"method.name": 			"(dbs *DatabaseDrivers) init(...)", 
							"db.client.engine": 	"etcd", 
							"db.client.version": 	"v3", 
							"var.keytoWatch": 		keytoWatch,
							"var.dialTimeout": 		dialTimeout,
							"var.etcdAP": 			etcdAP,
							}).Error("ETCD connection failed.")
		// log.Err(ctx, "ETCD connection failed ", blog.Fields{"err": err.Error()})
		time.Sleep(10 * time.Second)
		panic("ETCD connection failed")
	}
	defer cli.Close()
	log.WithFields(
		logrus.Fields{	
			"method.name": 			"(dbs *DatabaseDrivers) init(...)", 
			"db.client.engine": 	"etcd", 
			"db.client.version": 	"v3", 
			"var.keytoWatch": 		keytoWatch,
			"var.dialTimeout": 		dialTimeout,
			"var.etcdAP": 			etcdAP,
			}).Error("Watching for change")
	rch := cli.Watch(context.Background(), keytoWatch, clientv3.WithPrefix())
	for wresp := range rch {
		for _, ev := range wresp.Events {
			bdelete := false
			if ev.Type == clientv3.EventTypeDelete {
				bdelete = true
			}
			callbk(cbkctx, bdelete, string(ev.Kv.Key), string(ev.Kv.Value))
		}
	}
}

/*
func (dbs *DatabaseDrivers) initEtcd(hosts []string, timeout time.Duration, verbose bool) error {
//(db *DatabaseDrivers) func initEtcd(hosts []string, timeout time.Second, verbose bool) error {
	cfg := etcd.Config{
		Endpoints:               hosts,
		Transport:               etcd.DefaultTransport,
		HeaderTimeoutPerRequest: timeout,
	}
	cli, err := etcd.New(cfg)
	if err != nil {
		log.WithError(err).WithFields(
			logrus.Fields{	"db": "InstantiateEtcd", 
							"engine": "etcd", 
							"cfg": cfg,
							}).Warnf("Failed to connect to ETCD")
		return err
	}
	etcdClient := etcd.NewKeysAPI(cli)
	_, err = etcdClient.Get(context.Background(), "/sniperkit", nil)
	if err != nil && err.Error() == etcd.ErrClusterUnavailable.Error() {
		log.WithError(err).WithFields(
			logrus.Fields{	"db": "InstantiateEtcd", 
							"engine": "etcd", 
							"cfg": cfg,
							}).Warnf("error while init the database with boltDB.")
		return err
	}
	dbs.etcdCli 		= cli
	dbs.etcdKeys 		= etcdClient
	return nil
}
*/

func (dbs *DatabaseDrivers) getEtcd() etcd.KeysAPI {
    return dbs.etcdCli
}

func (dbs *DatabaseDrivers) closeEtcd() {
    dbs.etcdCli.Close()
}

//func (o *DatabaseDrivers) LoadDefaults() {
//}

//func (dbs *DatabaseDrivers) LoadDefaults() {
//}

// https://github.com/skyrunner2012/xormplus/blob/master/xorm/dataset.go
// NewDataset creates a new Dataset.
func NewDataset(headers []string) *tablib.Dataset {
	return tablib.NewDataset(headers)
}

// NewStarDump(ds)
func NewDump(content []byte, dumpPrefixPath string, dumpType string, dataFormat []string) (error) {
	ds, err := tablib.LoadJSON(content)
	if err != nil {
		log.WithError(err).WithFields(logrus.Fields{"method": "NewStarDump", "call": "LoadJSON"}).Debug("failed to load LoadJSON() with content")
		// panic(err)
		return err
	}
	if err := os.MkdirAll(dumpPrefixPath, 0777); err != nil {
		log.WithError(err).WithFields(logrus.Fields{"method": "NewStarDump", "call": "MkdirAll"}).Debugf("MkdirAll error on %#s", dumpPrefixPath)
		// panic(err)
		return err
	}
	for _, t := range dataFormat {
		filePath  := path.Join(dumpPrefixPath, dumpType+"."+t) // fmt.Sprintf("%s/%s", dumpPrefixPath, "repository.yaml") // will create a function
		file, err := os.Create(filePath)
		if err != nil {
			log.WithError(err).WithFields(logrus.Fields{"method": "NewStarDump", "call": "WriteTo"}).Errorf("%#v Write to %#v", t, filePath)
			// panic(err)
			return err
		}
		defer file.Close()
		switch df := t; df {
		case "json":
			json, err := ds.JSON()
			if err != nil {
				// panic(err)
				return errors.New("Error while converting data to "+df+" format")
			}
			json.WriteTo(file)
			// log.WithFields(logrus.Fields{"method": "NewStarDump", "call": "WriteTo"}).Debugf("%#v Write to %#v",  df, filePath)
		case "yaml":
			yaml, err := ds.YAML()
			if err != nil {
				// panic(err)
				return errors.New("Error while converting data to "+df+" format")
			}
			yaml.WriteTo(file)
			// log.WithFields(logrus.Fields{"method": "NewStarDump", "call": "WriteTo"}).Debugf("%#v Write to %#v",  df, filePath)
		case "csv":
			csv, err := ds.CSV()
			if err != nil {
				// panic(err)
				return errors.New("Error while converting data to "+df+" format")
			}
			csv.WriteTo(file)
			// log.WithFields(logrus.Fields{"method": "NewStarDump", "call": "WriteTo"}).Debugf("%#v Write to %#v",  df, filePath)
		case "xml":
			xml, err := ds.XML()
			if err != nil {
				// panic(err)
				return errors.New("Error while converting data to "+df+" format")
			}
			xml.WriteTo(file)
			// log.WithFields(logrus.Fields{"method": "NewStarDump", "call": "WriteTo"}).Debugf("%#v Write to %#v",  df, filePath)
		case "markdown":
			ascii := ds.Tabular("markdown")
			if ascii == nil {
				// panic(err)
				return errors.New("Error while converting data to "+df+" format")
			}
			ascii.WriteTo(file)
			// log.WithFields(logrus.Fields{"method": "NewStarDump", "call": "WriteTo"}).Debugf("%#v Write to %#v",  df, filePath)
		default:
			return errors.New("Unsupported data format: "+df)
		}
		file.Close()
	}
	return nil
}

func TimeToMicroseconds(t time.Time) int64 {
	return t.Unix()*int64(time.Second/time.Microsecond) + int64(t.Nanosecond())/int64(time.Microsecond)
}

/*

// https://github.com/Termina1/starlight/blob/93bd58b4c4795ca12b9fc849db9e4e3b0c668ca4/star_repo.go

func ReindexReposMongoDB(coll *mgo.Collection, batch int) *mgo.Iter {
  return coll.Find(bson.M{}).Batch(batch).Iter()
}

func RepoUpdateMongoDB(coll *mgo.Collection, name string, repo *StarRepo) {
  coll.Upsert(bson.M{"name": name}, repo)
}
*/

/*
// Open opens the connection to the bolt database defined by path.
func (d *DatabaseDriver) OpenBoltDriver(path string) error {
	if d.bucket != nil {
		return errors.New("store alread open")
	}

	store, err := bolt.Open(path, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.WithError(err).WithFields(logrus.Fields{"db": "OpenBoltDriver", "engine": "boltdb", "path": path}).Warnf("error while opening the boltDB bucket.")
		return err
	}

	err = store.Update(func(tx *bolt.Tx) error {
		buckets := [][]byte{
			store.importsBucket,
		}
		for _, bucket := range buckets {
			_, err := tx.CreateBucketIfNotExists(bucket)
			if err != nil {
				log.WithError(err).WithFields(logrus.Fields{"db": "OpenBoltDriver", "method": "CreateBucketIfNotExists", "engine": "boltdb", "path": path}).Warnf("error while creating the boltDB bucket if missing.")
				return err
			}
		}

		return nil
	})

	d.bucket = store
	return nil
}

// Close closes the underlying database.
func (d *DatabaseDriver) CloseBoltDriver() error {
	if d.bucket != nil {
		err := d.bucket.Close()
		d.bucket = nil
		log.WithError(err).WithFields(logrus.Fields{"db": "CloseBoltDriver", "method": "d.bucket.Close()", "engine": "boltdb"}).Warnf("error while closing the boltDB bucket.")
		return err
	}
	return nil
}
*/

/*
// init caley
func InitCaleyGraph(filepath string) (*bolt.DB, error) {
	// init command(s)
	if err != nil {
		log.WithError(err).WithFields(logrus.Fields{"db": "InitCaleyGraph", "engine": "").Warnf("error while init the full text search engine service.")
		return nil, err
	}
	return db, err
}
*/

/*
// init neo4j
func InitNeo4J(filepath string) (*bolt.DB, error) {
	// init command(s)
	if err != nil {
		log.WithError(err).WithFields(logrus.Fields{"db": "InitNeo4J", "engine": "neo4j", "drivers": "").Warnf("error while init the full text search engine service.")
		return nil, err
	}
	return db, err
}
*/

/*
// init dgraph
func InitDGraph(filepath string) (*bolt.DB, error) {
	// init command(s)
	if err != nil {
		log.WithError(err).WithFields(logrus.Fields{"db": "InitDGraph", "engine": "dgraph", "drivers": "").Warnf("error while init the full text search engine service.")
		return nil, err
	}
	return db, err
}
*/

/*
// init elasticsearch
func InitElasticsearch(filepath string) (*bolt.DB, error) {
	// init command(s)
	if err != nil {
		log.WithError(err).WithFields(logrus.Fields{"db": "InitElasticsearch", "engine": "elasticsearch", "drivers": "").Warnf("error while init the full text search engine service.")
		return nil, err
	}
	return db, err
}
*/

/*
// init sphinxsearch (for corpuses and dictionaries)
func InitSphinxSearch(filepath string) (*bolt.DB, error) {
	// init command(s)
	if err != nil {
		log.WithError(err).WithFields(logrus.Fields{"db": "InitSphinxSearch", "engine": "sphinxsearch", "drivers": "").Warnf("error while init the full text search engine service.")
		return nil, err
	}
	return db, err
}
*/

/*
// init mongodb
func InitMongoDB(filepath string) (*bolt.DB, error) {
	// init command(s)
	if err != nil {
		log.WithError(err).WithFields(logrus.Fields{"db": "InitMongoDB", "engine": "mongodb", "drivers": "").Warnf("error while connecting to the NoSQL data-store service.")
		return nil, err
	}
	return db, err
}
*/

/*
// init cassandra
func InitCassandraDB(filepath string) (*bolt.DB, error) {
	// init command(s)
	if err != nil {
		log.WithError(err).WithFields(logrus.Fields{"db": "InitCassandraDB", "engine": "cassandra", "drivers": "").Warnf("error while init the key/value store service.")
		return nil, err
	}
	return db, err
}
*/

/*
// init redis
func InitRedis(filepath string) (*bolt.DB, error) {
	// init command(s)
	if err != nil {
		log.WithError(err).WithFields(logrus.Fields{"db": "InitRedis", "engine": "redis", "drivers": "").Warnf("error while init the key/value store service.")
		return nil, err
	}
	return db, err
}
*/

/*
// init webdis
func InitWebdis(filepath string) (*bolt.DB, error) {
	// init command(s)
	if err != nil {
		log.WithError(err).WithFields(logrus.Fields{"db": "InitWebdis", "engine": "webdis", "drivers": "").Warnf("error while init the key/value store service.")
		return nil, err
	}
	return db, err
}
*/

/*
// init memcached
func InitMemcached(filepath string) (*bolt.DB, error) {
	// init command(s)
	if err != nil {
		log.WithError(err).WithFields(logrus.Fields{"db": "InitMemcached", "engine": "memcached", "drivers": "").Warnf("error while init the key/value store service.")
		return nil, err
	}
	return db, err
}
*/



/*
func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// ref. https://github.com/BalkanTech/goilerplate/blob/master/databases/databases.go

//NewGormConnection reads the provided config and returns an active Gorm database connection or an error
func NewGormConnection(c *config.Config) (*gorm.DB, error) {
	if 	GetType(&c) != "sqlite3" &&
		GetType(&c) != "postgres" &&
		GetType(&c) != "mysql" {
		return nil, ErrNotGorm
	}
	s, err := GetDBConnectionString()
	if err != nil {
		return nil, err
	}
	db, err := gorm.Open(GetType(), s)
	if cfg.Debug {
		db.LogMode(true)
	}
	return db, err
}

type MongoDBConnection struct {
	Session *mgo.Session
	DB      *mgo.Database
}

//NewMgoConnection reads the provided config and returns an active MGO session or an error
func NewMgoConnection(c *config.Config) (*MongoDBConnection, error) {
	if GetType(&c) != "mongodb" {
		return nil, ErrNotMongoDB
	}
	s, err := GetDBConnectionString()
	if err != nil {
		return nil, err
	}
	session, err := mgo.Dial(s)
	if err != nil {
		return nil, err
	}
	mode, err := GetMongoMode()
	if err != nil {
		return nil, err
	}
	session.SetMode(mode, true)
	db := session.DB("") // The DB name has been provided via the dial string
	return &MongoDBConnection{Session: session, DB: db}, nil
}


func initWithGorm(db *gorm.DB) {
  	db.AutoMigrate(	&cfg.Config{}, 
  					&cfg.SMTPConfig{}, 
  					&cfg.LogConfig{}, 
  					//&cfg.Directories{}, 
  					&cfg.ServiceConfig{}, 
  					&cfg.EndpointConfig{},
  					&cfg.Backend{},
  					&cfg.EngineConfig{},
  					&utils.OptionsSift{})
}

func initWithMGO(db *MongoDBConnection) {

}

func initWithBoltDB(db *bolt.DB) {

}

func DatabaseInit(c *config.Config) (e error) {
	if !IsValidDatabaseType() {
		return fmt.Errorf("Invalid database type in config file")
	}
	if IsGorm() {
		DB, err := NewGormConnection(c)
		if err != nil {
			return err
		}
		initWithGorm(DB)
		AdminInit(DB)
	}
	if IsMGO() {
		DB, err := NewMgoConnection(c)
		if err != nil {
			return err
		}
		initWithMGO(DB)
	}
	if IsBoltDB() {
		DB, err := NewBoltdbConnection(c)
		if err != nil {
			return err
		}
		initWithBoltDB(DB)
	}
	return nil
}

// GetType returns the database type in lowercase
func GetType(c *config.Config) string {
	return strings.ToLower(c.Database.Adapter)
}

func getMongoDBConnectionString(c *config.Config) (string, error)  {
	if GetType(&c) != "mongodb" {
		return "", &configError{"Database:Type", "Field not or incorrectly set"}
	}
	if c.Database.Host == "" {
		return "", &configError{"Database:Host", "Field not set"}
	}
	if c.Database.Name == "" {
		return "", &configError{"Database:Name", "Field not set"}
	}
	str := "mongodb://"
	if c.Database.User != "" {
		str += fmt.Sprintf("%s:%s@%s/%s", c.Database.User, c.Database.Password, c.Database.Host, c.Database.Name)
	} else {
		str += fmt.Sprintf("%s/%s", c.Database.Host, c.Database.Name)
	}
	return str, nil
}

// GetMongoMode returns a mgo.Mode based upon the settings of the configuration file. The default mode is mgo.Strong
func GetMongoMode(c *config.Config) (mgo.Mode, error)  {
	if GetType(&c) != "mongodb" {
		return -1, &configError{"Database: Adapter", "Field not or incorrectly set"}
	}
	switch(strings.ToLower(c.Database.Mode)) {
		case "primary":
			return mgo.Primary, nil
		case "primary_preferred":
			return mgo.PrimaryPreferred, nil
		case "secondary":
			return mgo.Secondary, nil
		case "secondary_preferred":
			return mgo.SecondaryPreferred, nil
		case "nearest":
			return mgo.Nearest, nil
		case "eventual":
			return mgo.Eventual, nil
		case "monotonic":
			return mgo.Monotonic, nil
		case "strong":
			return mgo.Strong, nil
		default:
			return mgo.Strong, nil
	}
}

func getBoltDBConnectionString(c *config.Config) (string, error) {
	if GetType(&c) != "boltdb" {
		return "", &configError{"Database: Adapter", "Field not or incorrectly set"}
	}
	if c.Database.Name == "" {
		return "", &configError{"Database: Name", "Field not set"}
	}
	return path.Join(c.Dirs.Shared, c.Dirs.Data, c.App.Info.Name, fmt.Sprintf("%s.boltdb", c.Database.Name)), nil
}

func getLocalDatabaseDirectoryString(c *config.Config) (string, error) {
	if 	GetType(&c) != "sqlite3" &&
		GetType(&c) != "boltdb" {
		return "", &configError{"Database: Adapter", "Not a local/embedded database"}
	}
	if c.Database.Name == "" {
		return "", &configError{"Database: Name", "Field not set"}
	}
	return path.Join(c.Dirs.Shared, c.Dirs.Data, c.App.Info.Name), nil
}

func setLocalDatabaseDirectory() (bool) {
	localDatabaseDirectory, err := getLocalDatabaseDirectoryString()
	if err != nil {
		return false
	}
	if err := os.MkdirAll(localDatabaseDirectory, 0700); err != nil {
		return false
	}
	return true
}

func getSqLite3ConnectionString(c *config.Config) (string, error) {
	if GetType(&c) != "sqlite3" {
		return "", &configError{"Database: Adapter", "Field not or incorrectly set"}
	}
	if c.Database.Name == "" {
		return "", &configError{"Database: Name", "Field not set"}
	}
	dbPath := path.Join(c.Dirs.Shared, c.Dirs.Data, c.App.Info.Name, fmt.Sprintf("%s.sqlite3", c.Database.Name))
	return dbPath, nil
}

func getPostgresConnectionString(c *config.Config) (string, error) {
	if GetType(&c) != "postgres" {
		return "", &configError{"Database: Adapter", "Field not or incorrectly set"}
	}
	if c.Database.Host == "" {
		return "", &configError{"Database: Host", "Field not set"}
	}
	if c.Database.Name == "" {
		return "", &configError{"Database: Name", "Field not set"}
	}
	ssl := "disable"
	if c.Database.SSLMode {
		ssl = "enable"
	}
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=%s", c.Database.Host, c.Database.User, c.Database.Password, c.Database.Name, ssl), nil
}

func getMySQLConnectionString(c *config.Config) (string, error) {
	if GetType(&c) != "mysql" {
		return "", &configError{"Database: Adapter", "Field not or incorrectly set"}
	}
	if c.Database.Name == "" {
		return "", &configError{"Database: Name", "Field not set"}
	}
	if c.Database.Local == "" {
		c.Database.Local = "Local"
	}
	if c.Database.Charset == "" {
		c.Database.Charset = "utf8"
	}
	if strings.ToLower(c.Database.Host) == "localhost" || c.Database.Host == "127.0.0.1" {
		c.Database.Host = ""
	}
	return fmt.Sprintf("%s:%s@%s/%s?charset=%s&parseTime=%t&loc=%s", c.Database.User, c.Database.Password, c.Database.Host, c.Database.Name, c.Database.Charset, c.Database.ParseTime, c.Database.Local), nil
}

// GetDBConnectionString creates and returns a formatted database connection string
(c *config.Config) func GetDBConnectionString() (string, error) {
	switch GetType(&c) {
		case "boltdb":
			return getBoltDBConnectionString()
		case "sqlite3":
			return getSqLite3ConnectionString()
		case "postgres":
			return getPostgresConnectionString()
		case "mysql":
			return getMySQLConnectionString()
		case "mongodb":
			return getMongoDBConnectionString()
		default:
			return "", &configError{"Database:Type", "Unsported database type"}
	}
}

// IsGorm returns true or false whether the database type is set to use a Gorm type of database
(c *config.Config) func IsGorm() bool {
	dbtype := GetType(c)
	return dbtype == "sqlite3" || dbtype == "postgres" || dbtype == "mysql"
}

// IsMGO returns true or false whether the database type is set to use a MongoDB type of database
(c *config.Config) func IsMGO() bool {
	return GetType(c) == "mongodb"
}

(c *config.Config) func IsBoltDB() bool {
	return GetType(c) == "boltdb"
}

// IsMGO returns true or false whether the database type is set to use a MongoDB type of database
(c *config.Config) func IsDirDatabase() bool {
	//return cfg.GetType() == "mongodb"
	return setLocalDatabaseDirectory()
}

// IsValidType returns true if the Type is set to a valid value, false if set to a false value
func IsValidDatabaseType() bool {
	return IsGorm() || IsMGO() || IsBoltDB()
}
*/

/*
// Open opens the connection to the bolt database defined by path.
func (d *Driver) OpenBoltDriver(path string) error {
	if d.bucket != nil {
		return errors.New("store alread open")
	}

	store, err := bolt.Open(path, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.WithError(err).WithFields(logrus.Fields{"db": "OpenBoltDriver", "engine": "boltdb", "path": path}).Warnf("error while opening the boltDB bucket.")
		return err
	}

	err = store.Update(func(tx *bolt.Tx) error {
		buckets := [][]byte{
			importsBucket,
		}
		for _, bucket := range buckets {
			_, err := tx.CreateBucketIfNotExists(bucket)
			if err != nil {
				log.WithError(err).WithFields(logrus.Fields{"db": "OpenBoltDriver", "method": "CreateBucketIfNotExists", "engine": "boltdb", "path": path}).Warnf("error while creating the boltDB bucket if missing.")
				return err
			}
		}

		return nil
	})

	d.bucket = store
	return nil
}

// Close closes the underlying database.
func (d *Driver) CloseBoltDriver() error {
	if d.bucket != nil {
		err := d.bucket.Close()
		d.bucket = nil
		log.WithError(err).WithFields(logrus.Fields{"db": "CloseBoltDriver", "method": "d.bucket.Close()", "engine": "boltdb"}).Warnf("error while closing the boltDB bucket.")
		return err
	}
	return nil
}
*/

/*
// move this piece of code into an admin dedicated file.
// qor admin - web ui
// qo beego - 
func InitAdmin(db *gorm.DB) (error) {
	// Initalize
	adm, err := admin.New(&qor.Config{DB: &db.DB})
	if err != nil {
		log.WithError(err).WithFields(logrus.Fields{"db": "InitAdmin", "action": "admin.New").Warnf("error while init the admin webui powered by qor-admin.")
		return err
	}
	adm.AddResource(&db.User{}, &admin.Config{Menu: []string{"Limo"}})
	mux, err := http.NewServeMux()
	if err != nil {
		log.WithError(err).WithFields(logrus.Fields{"db": "InitAdmin", "action": "NewServeMux").Warnf("error while init the mux web-server.")
		return err
	}
	adm.MountTo("/admin", mux)
	beego.Handler("/admin", mux)
	beego.Handler("/admin/*", mux)
	beego.Run()
}
*/
