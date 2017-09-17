package model

// https://github.com/ssut/pocketnpm/blob/master/db/gorm_backend.go
// 

import (
    // "path/filepath" 												// go-core
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
	"github.com/cayleygraph/cayley" 								// db-graph-cayley
	// "github.com/cayleygraph/cayley/graph" 						// db-graph-cayley
	// _ "github.com/cayleygraph/cayley/graph/bolt" 				// db-graph-cayley
	// "github.com/cayleygraph/cayley/quad" 						// db-graph-cayley
	"github.com/ckaznocha/taggraph" 								// db-graph-taggraph
	"github.com/blevesearch/bleve" 									// data-index-search
	tablib "github.com/agrison/go-tablib" 							// data-processing-tablib
	// jsoniter "github.com/json-iterator/go" 						// data-processing-jsoniter
	// "github.com/davecgh/go-spew/spew" 							// data-debug
	"github.com/k0kubun/pp" 										// debug-print
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

// var DBS  = DatabaseDrivers

// ref. https://github.com/tinrab/go-mmo/blob/master/db/dbobjects_gen.go
//type Database interface {
//	Dial(cfg *Config) error
//	Close()
//}

//	globalSetting := make(map[string]string)
// https://github.com/thesyncim/365a/blob/master/server/app.go
// https://github.com/emotionaldots/arbitrage/blob/master/cmd/arbitrage-db/main.go

var (
	DefaultSql 				= "sqlite3"
	DefaultKvs 				= map[string]bool{"boltdb": true, "etcd": true}
	DefaultGraphs 			= map[string]bool{"neo4j": true,  "cayley": true}
)

type GormRes struct {
	Ok 				bool  				`default:"false" json:"-" yaml:"-"`
	Cli 			*gorm.DB 			`json:"-" yaml:"-"`
}

type BoltRes struct {
	Ok 				bool  				`default:"false" json:"-" yaml:"-"`
	Cli 			*bolt.DB 			`json:"-" yaml:"-"`
}

type DatabaseDrivers struct {

	Initialized 		bool 				`default:"false" json:"-" yaml:"-"`
	Gorm struct {
		Ok 				bool  				`default:"false" json:"-" yaml:"-"`
		Cli 			*gorm.DB 			`json:"-" yaml:"-"`
	}

	Bolt struct {
		Ok 				bool  				`default:"false" json:"-" yaml:"-"`
		Cli 			*bolt.DB 			`json:"-" yaml:"-"`		
	}

	Redis struct {
		Ok 				bool  				`default:"false" json:"-" yaml:"-"`
		Cli 			redis.Conn 			`json:"-" yaml:"-"`
	}

	Etcd struct {
		Ok 				bool  				`default:"false" json:"-" yaml:"-"`
		Cli 			*etcd.Client 		`json:"-" yaml:"-"`
		Kvc 			etcd.KV 			`json:"-" yaml:"-"`
		//pool 			*EtcdClientPool 	`json:"-" yaml:"-"`
	}

	Bleve struct {
		Ok 				bool  				`default:"false" json:"-" yaml:"-"`
		Cli 			bleve.Index 		`json:"-" yaml:"-"`
	}

	Neo4j struct {
		Ok 				bool  				`default:"false" json:"-" yaml:"-"`
		Cli 			*neoism.Database 	`json:"-" yaml:"-"`
	}

	Cayley struct {
		Ok 				bool  				`default:"false" json:"-" yaml:"-"`
		Cli 			*cayley.Handle 		`json:"-" yaml:"-"`		
	}

	BigCache struct {
		Ok 				bool  				`default:"false" json:"-" yaml:"-"`
		Cli 			*bigcache.BigCache 	`json:"-" yaml:"-"`
	}

}

var (
	RootConf 			*config.Config 				//
	RootDrivers 		*DatabaseDrivers 			//
	//RootDrivers 		*DatabaseDrivers 			//
	AdminUI 			struct {
		Ok 				bool  						`default:"false" json:"-" yaml:"-"`
		Res				*admin.Admin 				`default:"false" json:"-" yaml:"-"`
	}
	log 				= logrus.New()
	tagg 				= taggraph.NewTagGaph()
	keyValMap 			map[string]string
)

func init() {
	log.Out 			= os.Stdout 					// logs
	formatter 			:= new(prefixed.TextFormatter) 	// logs
	log.Formatter 		= formatter 					// logs
	log.Level 			= logrus.DebugLevel 			// logs
	RootDrivers 		= &DatabaseDrivers{Initialized: false}
	if ! RootDrivers.Initialized {
		RootDrivers  	= New(DefaultSql, DefaultKvs, DefaultGraphs)
		log.WithFields(
			logrus.Fields{	
				"src.file": 				"model/data-core.go", 
				"method.name": 				"init(...)", 
				"method.prev": 				"New(...)",
				"prefix": 					"dbs-init",
				"var.Drivers.Gorm.Ok": 		RootDrivers.Gorm.Ok,
				"var.Drivers.Bolt.Ok": 		RootDrivers.Bolt.Ok,
				//"var.Drivers.Ectd.Ok": 	Drivers.Ectd.Ok,
				//"var.Drivers.Caley.Ok": 	Drivers.Caley.Ok,
				//"var.Drivers.Neo4j.Ok": 	Drivers.Neo4j.Ok,
				}).Debug("error while trying to init all datastore drivers.")			
	}
}

func GetDrivers() *DatabaseDrivers {
	pp.Print(RootDrivers)
	return RootDrivers
}

//func New(cfg *config.Config, verbose bool, debug bool) (db *DatabaseDrivers, err error) {
func New(sqlEngine string, kvEngines map[string]bool, graphEngines map[string]bool) *DatabaseDrivers {

	var gormClientStatus, boltClientStatus bool
	log.WithFields(
		logrus.Fields{	"src.file": 				"model/data-core.go", 
						"method.name": 				"New(...)", 
						"var.sqlEngine": 			sqlEngine, 
						"var.kvEngines": 			kvEngines, 
						"var.graphEngines": 		graphEngines,
						"struct.dbs": 				RootDrivers,
						"prefix": 					"dbs-list",
						}).Info("error while trying to init 'Gorm' database driver.")

	gormClient, err  := InitGorm("./shared/data/limo/limo.db", "sqlite3")
	if err == nil {
		gormClientStatus = true
	} else {
		log.WithError(err).WithFields(
			logrus.Fields{	"src.file": 			"model/data-core.go", 
							"method.name": 			"New(...)", 
							"db.type": 				"sql", 
							"db.driver": 			"gorm", 
							"db.adpater": 			sqlEngine,
							"struct.dbs": 			RootDrivers,
							"method.prev": 			"db.initGorm(...)",
							"prefix": 				"dbs-new",
							}).Error("error while trying to init 'Gorm' database driver.")
		//return drivers
	}
	// gormRes 	:= &GormRes{Ok: gormClientStatus, Cli: gormClient}
	RootDrivers.Gorm.Ok 			= gormClientStatus
	RootDrivers.Gorm.Cli 			= gormClient

	boltClient, err 	:= InitBoltDB("./shared/data/limo/limo.boltdb")
	if err == nil {
		boltClientStatus = true
	} else {
		log.WithError(err).WithFields(
			logrus.Fields{	"src.file": 			"model/data-core.go", 
							"method.name": 			"New(...)", 
							"db.type": 				"kvs", 
							"db.driver": 			"bolt", 
							"db.adapter": 			"boltdb", 
							"method.prev": 			"db.initBoltDB(...)",
							"var.kvEngines": 		kvEngines,
							"struct.dbs": 			RootDrivers,
							"prefix": 				"dbs-new",
							"action": 				"InitBoltDB",
							}).Error("error while trying to init 'BoltDB' database driver")
		// return drivers
	}
	//boltRes 	:= &BoltRes{Ok: boltClientStatus, Cli: boltClient}
	RootDrivers.Bolt.Ok 			= boltClientStatus
	RootDrivers.Bolt.Cli 			= boltClient

	/*
	if err := AutoLoadGorm(gormClient, true, true, true, Tables...); err != nil {
		log.WithError(err).WithFields(
			logrus.Fields{	"src.file": 		"model/data-core.go", 
							"method.name": 		"New(...)", 
							"db.type": 			"sql", 
							"db.driver": 		"gorm", 
							"db.adpater": 		sqlEngine,
							"method.prev": 		"db.autoLoadGorm(...)",
							"prefix": 			"dbs-new",
							}).Error("error while trying to init 'Gorm' database driver.")
		return err
	}
	*/

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

	/*
	etcdDefaultHost 	:= []string{"http://127.0.0.1:2379"}
	etcdDefaultTimeout 	:= 1 * time.Second
	etcdCli, err := InitEtcd(etcdDefaultHost, etcdDefaultTimeout)
	if err != nil {
		log.WithError(err).WithFields(
			logrus.Fields{	"src.file": 		"model/data-core.go", 
							"prefix": 			"dbs-new",
							"db.type": 			"kvs", 
							"db.driver": 		"etcd", 
							"db.adapter": 		"etcd", 
							"method.prev": 		"db.initEtcd(...)",
							"action": 			"AutoloadDB",
							}).Error("error while trying to auto-load all program the tables")
		return dbs, err
	}
	*/

	// graphql
	// client := graphql.NewClient("https://example.com/graphql", nil, nil)
	return nil

}

// https://github.com/qor/qor-example/blob/master/db/db.go
// InitDB initializes the database at the specified path
func InitGorm(filepath string, adapter string) (*gorm.DB, error) {
//func InitDB(filepath string, adapter string, verbose bool) (*gorm.DB, error) {
	gormDB, err := gorm.Open(adapter, filepath) 	// Get more config options to setup the SQL database
	if err != nil {
		log.WithError(err).WithFields(
			logrus.Fields{	"prefix": 			"dbs-init",
							"db.adapter": 		adapter,
							"src.file": 		"model/data-core.go", 
							"method.name": 		"InitGorm(...)", 							
							"method.prev": 		"gorm.Open(...)",
							}).Warn("error while init the database with gorm.")
		return gormDB, err
	}
	gormDB.LogMode(false) 	// cfg.App.DebugMode
	RootDrivers.Gorm.Ok 			= true
	RootDrivers.Gorm.Cli 			= gormDB
	/*
	// isAutoMigrate, isTruncate, isAdminUIResource
	if err := gormDB.AutoloadDB(true, true, false, Tables...); err != nil {
		log.WithError(err).WithFields(
			logrus.Fields{	"file": 		"model/data-core.go", 
							"method.name": 	"InitGorm", 
							"adapter": 		adapter, 
							"action": 		"AutoloadDB",
							}).Warn("error while trying to auto-load all program the tables")
	}
	*/
	return RootDrivers.Gorm.Cli, nil
}

func AutoLoadGorm(db *gorm.DB, isAutoMigrate bool, isTruncate bool, isAdminUI bool, tables ...interface{}) error {
	AdminUI.Ok = isAdminUI
	if isAdminUI {
		AdminUI.Res = 	admin.New(&qor.Config{DB: db})
	}
	for _, table := range tables {
		if isTruncate {
			if err := db.DropTableIfExists(table).Error; err != nil {
				log.WithError(err).WithFields(
					logrus.Fields{	"src.file": 			"model/data-core.go", 
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
			if err := db.AutoMigrate(table).Error; err != nil {
				log.WithError(err).WithFields(
					logrus.Fields{	"src.file": 				"model/data-core.go", 
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
			AdminUI.Res.AddResource(table)
			log.WithFields(
				logrus.Fields{	"src.file": 					"model/data-core.go", 
								"method.name": 					"(db *DatabaseDrivers) autoLoadGorm(...)", 
								"method.prev": 					"adminUI.AddResource(...)",
								"var.adminui.status": 			isAdminUI,
								"var.adminui.table": 			table,
								}).Info("adding admin UI resource for the table")
		}
	}
	if isAdminUI {
		if len(AdminUI.Res.GetResources()) > 0 {
			for _, resource := range AdminUI.Res.GetResources() {	
				log.WithFields(
					logrus.Fields{	"src.file": 				"model/data-core.go", 
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
	boltDB, err 			:= bolt.Open(filePath, 0600, &bolt.Options{Timeout: time.Second})
	if err != nil {
		log.WithError(err).WithFields(
			logrus.Fields{	"prefix": 				"db-boltdb",
							"method.name": 			"InitBoltDB(...)", 
							"method.prev": 			"bolt.Open(...)",
							"db.adapter": 			"boltdb", 
							"src.file": 			"model/data-core.go", 
							"var.bolt.filepath": 	filePath, 
							"var.bolt.options": 	&bolt.Options{Timeout: time.Second},
							}).Warn("error while init the database with boltDB.")
		return boltDB, err
	}
	RootDrivers.Bolt.Ok 			= true
	RootDrivers.Bolt.Cli 			= boltDB
	return RootDrivers.Bolt.Cli, nil
}

func InitEtcd(hosts []string, timeOut time.Duration) (*etcd.Client, error) {
	ectdConfig := etcd.Config{
		Endpoints:               hosts,
		DialTimeout: 			 timeOut,
	}
	ectdClient, err := etcd.New(ectdConfig)
	if err != nil {
		log.WithError(err).WithFields(
			logrus.Fields{	"prefix": 				"kvs-etcd",
							"method.name": 			"InitEtcd(...)", 
							"method.prev": 			"etcd.New(...)",
							"db.adapter": 			"etcd", 
							"var.etcd.cfg": 		ectdConfig,
							"src.file": 			"model/data-core.go", 
							}).Warn("error while init the client connection with Etcd Key/Value store.")
		return ectdClient, err
	}
	defer ectdClient.Close()
	etcdKvc := etcd.NewKV(ectdClient)
	_, err = etcdKvc.Get(context.Background(), "/foo", nil)
	if err != nil {
		log.WithError(err).WithFields(
			logrus.Fields{	"method.name": 			"InitEtcd(...)", 
							"db.adapter": 			"etcd", 
							"prefix": 				"kvs-etcd",
							"method.prev": 			"etcdClient.Get(...)",
							"var.etcd.cfg": 		ectdConfig,
							"src.file": 			"model/data-core.go", 
							}).Warn("error while init the client connection with Etcd Key/Value store.")
		return ectdClient, err
	}
	RootDrivers.Etcd.Ok 			= true
	RootDrivers.Etcd.Cli 			= ectdClient
	RootDrivers.Etcd.Kvc 			= etcdKvc
	return RootDrivers.Etcd.Cli, nil
}

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
