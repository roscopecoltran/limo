package model

// https://github.com/ssut/pocketnpm/blob/master/db/gorm_backend.go
//

import (
	// "path/filepath" 															// go-core
	"errors"                                     // go-core
	"github.com/allegro/bigcache"                // data-cache-bigcache
	"github.com/jinzhu/gorm"                     // db-sql-gorm
	_ "github.com/jinzhu/gorm/dialects/mysql"    // db-sql-gorm-mysql
	_ "github.com/jinzhu/gorm/dialects/postgres" // db-sql-gorm-postgres
	_ "github.com/jinzhu/gorm/dialects/postgres" // db-sql-gorm-postgres
	_ "github.com/jinzhu/gorm/dialects/sqlite"   // db-sql-gorm-sqlite3
	"os"                                         // go-core
	"path"                                       // go-core
	"time"                                       // go-core
	// "gopkg.in/mgo.v2" 															// db-nosql-mongodb
	// "gopkg.in/mgo.v2/bson" 													// db-nosql-mongodb
	etcd "github.com/coreos/etcd/clientv3" // db-kvs-etcd
	//etcd "github.com/coreos/etcd/client" 										// db-kvs-etcd
	"github.com/boltdb/bolt"                  // db-kvs-boltdb
	"github.com/bradfitz/gomemcache/memcache" // db-kvs-memcache
	"github.com/cayleygraph/cayley"           // db-graph-cayley
	"github.com/garyburd/redigo/redis"        // db-kvs-redis
	"github.com/jmcvetta/neoism"              // db-graph-neo4j
	// "github.com/cayleygraph/cayley/graph" 									// db-graph-cayley
	// _ "github.com/cayleygraph/cayley/graph/bolt" 							// db-graph-cayley
	// "github.com/cayleygraph/cayley/quad" 									// db-graph-cayley
	"github.com/blevesearch/bleve"  // data-index-search
	"github.com/ckaznocha/taggraph" // db-graph-taggraph
	//elastic "gopkg.in/olivere/elastic.v2" 									// data-index-search
	tablib "github.com/agrison/go-tablib" // data-processing-tablib
	elastic "gopkg.in/olivere/elastic.v5" // data-index-search
	// jsoniter "github.com/json-iterator/go" 									// data-processing-jsoniter
	// "github.com/davecgh/go-spew/spew" 										// data-debug
	"github.com/k0kubun/pp" // debug-print
	// "github.com/astaxie/beego" 												// web-cms
	"github.com/qor/admin"                                 // web-qor-admin-ui
	"github.com/qor/qor"                                   // web-qor-admin-ui
	"github.com/roscopecoltran/sniperkit-limo/config"      // app-config
	"github.com/sirupsen/logrus"                           // logs-logrus
	prefixed "github.com/x-cray/logrus-prefixed-formatter" // logs-logrus
	"golang.org/x/net/context"                             // web-context
	//"github.com/roscopecoltran/sniperkit-limo/model/drivers/rdb/gorm" 		// dbs-client-sql
	//"github.com/roscopecoltran/sniperkit-limo/model/drivers/rdb/xorm" 		// dbs-client-sql
	//"github.com/roscopecoltran/sniperkit-limo/model/drivers/rdb/xorm-plus" 	// dbs-client-sql
	//"github.com/roscopecoltran/sniperkit-limo/model/drivers/kvs/boltdb" 		// dbs-client-kvs
	//etcd "github.com/roscopecoltran/sniperkit-limo/model/drivers/kvs/etcd/v2" // dbs-client-kvs
	//etcd "github.com/roscopecoltran/sniperkit-limo/model/drivers/kvs/etcd/v3" // dbs-client-kvs
	//"github.com/roscopecoltran/sniperkit-limo/model/drivers/kvs/redis" 		// dbs-client-kvs
	//"github.com/roscopecoltran/sniperkit-limo/model/drivers/graph/neo4j" 		// dbs-client-graph
	//"github.com/roscopecoltran/sniperkit-limo/model/drivers/graph/cayley" 	// dbs-client-graph
	//"github.com/roscopecoltran/sniperkit-limo/model/drivers/graph/dgraph" 	// dbs-client-graph
	//"github.com/roscopecoltran/sniperkit-limo/model/drivers/docs/mongodb" 	// dbs-client-docs
	//"github.com/roscopecoltran/sniperkit-limo/model/drivers/docs/cassandra" 	// dbs-client-docs
	//"github.com/roscopecoltran/sniperkit-limo/model/drivers/cache/ramdisk" 	// dbs-client-memory
	//"github.com/roscopecoltran/sniperkit-limo/model/drivers/cache/memcache" 	// dbs-client-memory
)

var validDataOutput = []string{"md", "csv", "yaml", "json", "xlsx", "xml", "tsv", "mysql", "postgres", "html", "ascii"} // datasets - formats
var availableLocales = []string{"en-US", "fr-FR", "pl-PL"}

var (
	Tables = []interface{}{
		&Service{}, &ExternalURL{},
		&User{}, &UserInfoVCS{}, &Email{}, &Address{},
		&Star{}, &Tag{}, &Topic{}, &PatternEntry{},
		&Tree{}, &TreeEntry{}, &Readme{},
		&Language{}, &LanguageType{}, // &Detection{},
	}
)

var (
	DefaultSql         = "sqlite3"
	DefaultKvs         = map[string]bool{"boltdb": true, "etcd": true}
	DefaultGraphs      = map[string]bool{"neo4j": true, "cayley": true}
	Default_Date_Short = "0000-00-00 00:00:00 -0000 UTC"
	Default_Date_Long  = "0000-00-00T00:00:00Z"
)

const Default_VCS_Github_Domain = "github.com"
const Default_VCS_Gitlab_Domain = "gitlab.com"
const Default_VCS_Bitbucket_Domain = "bitbucket.org"

/*
type Databases struct {
	Datastore 			map[string]*bolt.DB
	Database   			map[string]*gorm.DB
	SearchIdx 			map[string]*bleve.Index
	KvIdx 				map[string]etcd.KeysAPI
}
*/

type DatabaseDrivers struct {
	Initialized bool   `default:"false" json:"initialized" yaml:"initialized"`
	Mode        string `default:"dev" json:"mode" yaml:"mode"`

	Gorm struct {
		Ok  bool     `default:"false" json:"status" yaml:"status"`
		Cli *gorm.DB `json:"-" yaml:"-"`
		//Cluster 		map[string]*gorm.DB 			`json:"-" yaml:"-"`
		//Config 		map[string]ConfigGormDB 		`json:"config" yaml:"config"`
		Version float64 `json:"version,omitempty" yaml:"version,omitempty"`
	} `json:"gorm,omitempty" yaml:"gorm,omitempty"`

	// https://github.com/banlong/bleve/blob/master/storage/storage.go
	//
	Bolt struct {
		Ok  bool     `default:"false" json:"status" yaml:"status"`
		Cli *bolt.DB `json:"-" yaml:"-"`
		//Cluster map[string]*bolt.DB `json:"-" yaml:"-"`
		//Config 		map[string]ConfigBoltDB 		`json:"config" yaml:"config"`
		Version float64 `json:"version,omitempty" yaml:"version,omitempty"`
	} `json:"boltdb,omitempty" yaml:"boltdb,omitempty"`

	// https://github.com/mickeyinfoshan/gengo/blob/master/interfaces/selector.go
	/*
		MongoDB struct {
			Ok 				bool  							`default:"false" json:"status" yaml:"status"`
			Cli 			*mgo.Session 					`json:"-" yaml:"-"`
			//Cluster 		map[string]*mgo.Session 		`json:"-" yaml:"-"`
			//Config 		map[string]ConfigMongoDB 		`json:"config" yaml:"config"`
			Version 		float64 						`json:"version,omitempty" yaml:"version,omitempty"`
		} `json:"mongodb,omitempty" yaml:"mongodb,omitempty"`
	*/

	Redis struct {
		Ok  bool       `default:"false" json:"status" yaml:"status"`
		Cli redis.Conn `json:"-" yaml:"-"`
		//Cluster 		map[string]redis.Conn 			`json:"-" yaml:"-"`
		//Config 		map[string]ConfigRedisKVS 		`json:"config" yaml:"config"`
		Version float64 `json:"version,omitempty" yaml:"version,omitempty"`
		//Host			string 							`json:"host,omitempty" yaml:"host,omitempty"`
		//Port			string 							`json:"port,omitempty" yaml:"port,omitempty"`
		TimeOut float64 `default:"2.0" json:"timeout,omitempty" yaml:"timeout,omitempty"`
	} `json:"redis,omitempty" yaml:"redis,omitempty"`

	// https://github.com/bradfitz/gomemcache
	Memcache struct {
		Ok  bool             `default:"false" json:"status" yaml:"status"`
		Cli *memcache.Client `json:"-" yaml:"-"`
		//Cluster 		map[string]memcache.Client 		`json:"-" yaml:"-"`
		//Config 		map[string]ConfigMemcache 		`json:"config" yaml:"config"`
		Version float64 `json:"version,omitempty" yaml:"version,omitempty"`
		//Host			string 							`json:"host,omitempty" yaml:"host,omitempty"`
		//Port			string 							`json:"port,omitempty" yaml:"port,omitempty"`
		TimeOut float64 `default:"2.0" json:"timeout,omitempty" yaml:"timeout,omitempty"`
	} `json:"memcache,omitempty" yaml:"memcache,omitempty"`

	Etcd struct {
		Ok  bool         `default:"false" json:"status" yaml:"status"`
		Cli *etcd.Client `json:"-" yaml:"-"`
		Kvc etcd.KV      `json:"-" yaml:"-"`
		//Config 		map[string]ConfigEtcdKVS 		`json:"config" yaml:"config"`
		//Cluster 		map[string]*EtcdClientPool 		`json:"-" yaml:"-"`
		//pool 			*EtcdClientPool 				`json:"-" yaml:"-"`
		Version float64 `json:"version,omitempty" yaml:"version,omitempty"`
		//Hosts			[]string 						`json:"hosts,omitempty" yaml:"hosts,omitempty"`
		TimeOut float64 `default:"2.0" json:"timeout,omitempty" yaml:"timeout,omitempty"`
	} `json:"ectd,omitempty" yaml:"ectd,omitempty"`

	// https://github.com/rchardzhu/searchui/blob/master/models/searchengine.go
	Bleve struct {
		Ok  bool        `default:"false" json:"status" yaml:"status"`
		Cli bleve.Index `json:"-" yaml:"-"`
		//Config 		map[string]ConfigBleveIDX 		`json:"config" yaml:"config"`
		//Cluster 		map[string]bleve.Index 			`json:"-" yaml:"-"`
		Version float64 `json:"version,omitempty" yaml:"version,omitempty"`
	} `json:"bleve,omitempty" yaml:"bleve,omitempty"`

	// https://github.com/rchardzhu/searchui/blob/master/models/elastic_search.go
	// https://github.com/rchardzhu/searchui/blob/master/models/bleve_search.go
	Elastic struct {
		Ok  bool            `default:"false" json:"status" yaml:"status"`
		Cli *elastic.Client `json:"-" yaml:"-"`
		//Cluster 		map[string]elastic.Client 		`json:"-" yaml:"-"`
		//Config 		map[string]ConfigElasticIDX 	`json:"config" yaml:"config"`
		Version float64 `json:"version,omitempty" yaml:"version,omitempty"`
		//Host			string 							`json:"host,omitempty" yaml:"host,omitempty"`
		//Port			string 							`json:"port,omitempty" yaml:"port,omitempty"`
		//pool 			*ElasticClientPool 				`json:"-" yaml:"-"`
		TimeOut float64 `default:"2.0" json:"timeout,omitempty" yaml:"timeout,omitempty"`
	} `json:"elastic,omitempty" yaml:"elastic,omitempty"`

	// github.com/slvmnd/gosphinx
	// https://github.com/yunge/sphinx/blob/master/sphinx.go
	SphinxSearch struct {
		Ok bool `default:"false" json:"status" yaml:"status"`
		//Cli 			*sphinx.Client 					`json:"-" yaml:"-"`
		//Cluster 		map[string]*sphinx.Client 		`json:"-" yaml:"-"`
		//Config 		map[string]ConfigSphinxIDX 		`json:"config" yaml:"config"`
		Version float64 `json:"version,omitempty" yaml:"version,omitempty"`
		//Host			string 							`default:"localhost" json:"host,omitempty" yaml:"host,omitempty"`
		//Port			string 							`default:"9312" json:"port,omitempty" yaml:"port,omitempty"`
		//SqlPort		string 							`default:"9306" json:"port,omitempty" yaml:"port,omitempty"`
		TimeOut int `default:"5000" json:"timeout,omitempty" yaml:"timeout,omitempty"`
	} `json:"sphinx,omitempty" yaml:"sphinx,omitempty"`

	Neo4j struct {
		Ok  bool             `default:"false" json:"status" yaml:"status"`
		Cli *neoism.Database `json:"-" yaml:"-"`
		//Cluster 		map[string]*neoism.Database 	`json:"-" yaml:"-"`
		//Config 		map[string]ConfigNeo4jGraphDB 	`json:"config" yaml:"config"`
		Version float64 `json:"version,omitempty" yaml:"version,omitempty"`
		//Host			string 							`json:"host,omitempty" yaml:"host,omitempty"`
		//Port			string 							`json:"port,omitempty" yaml:"port,omitempty"`
		TimeOut float64 `default:"2.0" json:"timeout,omitempty" yaml:"timeout,omitempty"`
	} `json:"neo4j,omitempty" yaml:"neo4j,omitempty"`

	Cayley struct {
		Ok  bool           `default:"false" json:"status" yaml:"status"`
		Cli *cayley.Handle `json:"-" yaml:"-"`
		//Cluster 		map[string]*cayley.Handle 		`json:"-" yaml:"-"`
		//Config 		map[string]ConfigCayleyGraphDB 	`json:"config" yaml:"config"`
		Version float64 `json:"version,omitempty" yaml:"version,omitempty"`
	} `json:"cayley,omitempty" yaml:"cayley,omitempty"`

	BigCache struct {
		Ok  bool               `default:"false" json:"status" yaml:"status"`
		Cli *bigcache.BigCache `json:"-" yaml:"-"`
		//Mode 			map[string]*bigcache.BigCache 	`json:"-" yaml:"-"`
		//Config 		map[string]ConfigBigCacheKVS 	`json:"config" yaml:"config"`
		Version float64 `json:"version,omitempty" yaml:"version,omitempty"`
	} `json:"bigcache,omitempty" yaml:"bigcache,omitempty"`
}

var (
	RootConf    *config.Config   //
	RootDrivers *DatabaseDrivers //
	//RootDrivers 		*DatabaseDrivers 			//
	AdminUI struct {
		Ok  bool         `default:"false" json:"-" yaml:"-"`
		Res *admin.Admin `default:"false" json:"-" yaml:"-"`
	}
	log                 = logrus.New()
	tagg                = taggraph.NewTagGaph()
	keyValMap           map[string]string
	defaultBoltFilePath = "./shared/data/limo/limo.boltdb"
	defaultGormFilePath = "./shared/data/limo/limo.db"
	defaultGormAdapter  = "sqlite3"
)

func init() {
	log.Out = os.Stdout                      // logs
	formatter := new(prefixed.TextFormatter) // logs
	log.Formatter = formatter                // logs
	log.Level = logrus.DebugLevel            // logs
	RootDrivers = &DatabaseDrivers{Initialized: false}
	/*
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
	*/
}

func GetDrivers() *DatabaseDrivers {
	pp.Print(RootDrivers)
	return RootDrivers
}

//func New(cfg *config.Config, verbose bool, debug bool) (db *DatabaseDrivers, err error) {
func New(sqlEngine string, kvEngines map[string]bool, graphEngines map[string]bool) *DatabaseDrivers {

	RootDrivers.Initialized = true

	// init GORM drivers
	var gormClientStatus bool
	gormClient, err := InitGorm(defaultGormFilePath, defaultGormAdapter)
	if err == nil {
		gormClientStatus = true
	} else {
		log.WithError(err).WithFields(
			logrus.Fields{"prefix": "dbs-gorm",
				"group":                   "dbs-client-sql",
				"src.file":                "model/data-core.go",
				"method.name":             "New(...)",
				"method.prev":             "db.initGorm(...)",
				"db.type":                 "sql",
				"db.driver":               "gorm",
				"struct.dbs":              RootDrivers,
				"var.defaultGormAdapter":  defaultGormAdapter,
				"var.defaultGormFilePath": defaultGormFilePath,
				"var.adpater":             sqlEngine,
			}).Error("error while trying to init 'Gorm' database driver.")
		//return RootDrivers
	}
	RootDrivers.Gorm.Ok = gormClientStatus
	RootDrivers.Gorm.Cli = gormClient
	log.WithFields(
		logrus.Fields{"prefix": "dbs-gorm",
			"group":                   "data-client-sql",
			"src.file":                "model/data-core.go",
			"method.name":             "New(...)",
			"db.type":                 "sql",
			"db.driver":               "gorm",
			"var.sqlEngine":           sqlEngine,
			"var.defaultGormAdapter":  defaultGormAdapter,
			"var.defaultGormFilePath": defaultGormFilePath,
			"struct.dbs":              RootDrivers,
		}).Info("status 'Gorm' drivers")

	if err := AutoLoadGorm(RootDrivers.Gorm.Cli, true, true, true, Tables...); err != nil {
		log.WithError(err).WithFields(
			logrus.Fields{"src.file": "model/data-core.go",
				"method.name": "New(...)",
				"db.type":     "sql",
				"db.driver":   "gorm",
				"db.adpater":  sqlEngine,
				"method.prev": "db.autoLoadGorm(...)",
				"prefix":      "dbs-new",
			}).Error("error while trying to init 'Gorm' database driver.")
		// return err
	}

	// init BoltDB drivers
	var boltClientStatus bool
	boltClient, err := InitBoltDB(defaultBoltFilePath)
	if err == nil {
		boltClientStatus = true
	} else {
		log.WithError(err).WithFields(
			logrus.Fields{"prefix": "dbs-boltdb",
				"src.file":      "model/data-core.go",
				"method.name":   "New(...)",
				"db.type":       "kvs",
				"db.driver":     "bolt",
				"db.adapter":    "boltdb",
				"method.prev":   "db.initBoltDB(...)",
				"var.kvEngines": kvEngines,
				"struct.dbs":    RootDrivers,
			}).Error("error while trying to init 'BoltDB' database driver")
	}
	RootDrivers.Bolt.Ok = boltClientStatus
	RootDrivers.Bolt.Cli = boltClient
	log.WithFields(
		logrus.Fields{"prefix": "dbs-boltdb",
			"src.file":                 "model/data-core.go",
			"method.name":              "New(...)",
			"var.kvEngines":            kvEngines,
			"var.RootDrivers.Bolt.Ok":  RootDrivers.Bolt.Ok,
			"var.RootDrivers.Bolt.Cli": RootDrivers.Bolt.Cli,
			"var.defaultBoltFilePath":  defaultBoltFilePath,
		}).Info("status 'BoltDB' drivers")

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
	return RootDrivers

}

// https://github.com/qor/qor-example/blob/master/db/db.go
// InitDB initializes the database at the specified path
func InitGorm(filepath string, adapter string) (*gorm.DB, error) {
	//func InitDB(filepath string, adapter string, verbose bool) (*gorm.DB, error) {
	gormDB, err := gorm.Open(adapter, filepath) // Get more config options to setup the SQL database
	if err != nil {
		log.WithError(err).WithFields(
			logrus.Fields{"prefix": "dbs-init",
				"db.adapter":  adapter,
				"src.file":    "model/data-core.go",
				"method.name": "InitGorm(...)",
				"method.prev": "gorm.Open(...)",
			}).Warn("error while init the database with gorm.")
		return gormDB, err
	}
	gormDB.LogMode(false) // cfg.App.DebugMode
	RootDrivers.Gorm.Ok = true
	RootDrivers.Gorm.Cli = gormDB
	return gormDB, nil
}

func (dbs *DatabaseDrivers) Init(sqlEngine string, kvEngines map[string]bool, graphEngines map[string]bool) error {
	var err error
	err = dbs.initGorm(defaultGormFilePath, defaultGormAdapter)
	if err != nil {
		log.WithError(err).WithFields(
			logrus.Fields{"prefix": "dbs-init",
				"task":                    "dbs-init-sql-gorm",
				"src.file":                "model/data-core.go",
				"var.defaultGormFilePath": defaultGormFilePath,
				"var.defaultGormAdapter":  defaultGormAdapter,
				"ctx.method.name":         "(dbs *RootDrivers) Init(...)",
				"ctx.method.last":         "dbs.initGorm(...)",
			}).Error("error while init the database with gorm.")
		return err
	}
	err = dbs.initBoltDB(defaultBoltFilePath)
	if err != nil {
		log.WithError(err).WithFields(
			logrus.Fields{"prefix": "dbs-init",
				"task":                    "dbs-init-kvs-boltdb",
				"src.file":                "model/data-core.go",
				"var.defaultBoltFilePath": defaultBoltFilePath,
				"ctx.method.name":         "(dbs *RootDrivers) Init(...)",
				"ctx.method.last":         "dbs.initBoltDB(...)",
			}).Error("error while init the database with gorm.")
		return err
	}
	return nil
}

func (dbs *DatabaseDrivers) GetDrivers() *DatabaseDrivers {
	//pp.Print(dbs)
	return dbs
}

func (dbs *DatabaseDrivers) initGorm(filepath string, adapter string) error {
	gormDB, err := gorm.Open(adapter, filepath) // Get more config options to setup the SQL database
	if err != nil {
		log.WithError(err).WithFields(
			logrus.Fields{"prefix": "dbs-init",
				"db.adapter":  adapter,
				"src.file":    "model/data-core.go",
				"method.name": "(dbs *RootDrivers) initGorm(...)",
				"method.prev": "gorm.Open(...)",
			}).Warn("error while init the database with gorm.")
		return err
	}
	gormDB.LogMode(false) // cfg.App.DebugMode
	dbs.Gorm.Ok = true
	dbs.Gorm.Cli = gormDB
	if err := dbs.autoLoadGorm(true, true, false, Tables...); err != nil {
		log.WithError(err).WithFields(
			logrus.Fields{"file": "model/data-core.go",
				"method.name": "InitGorm",
				"adapter":     adapter,
				"action":      "AutoloadDB",
			}).Warn("error while trying to auto-load all program the tables")
		return err
	}
	return nil
}

func (dbs *DatabaseDrivers) initBoltDB(filePath string) error {
	boltDB, err := bolt.Open(filePath, 0600, &bolt.Options{Timeout: time.Second})
	if err != nil {
		log.WithError(err).WithFields(
			logrus.Fields{"prefix": "db-boltdb",
				"method.name":       "(dbs *RootDrivers) initGorm(...)",
				"method.prev":       "bolt.Open(...)",
				"db.adapter":        "boltdb",
				"src.file":          "model/data-core.go",
				"var.bolt.filepath": filePath,
				"var.bolt.options":  &bolt.Options{Timeout: time.Second},
			}).Warn("error while init the database with boltDB.")
		return err
	}
	dbs.Bolt.Ok = true
	dbs.Bolt.Cli = boltDB

	//if dbs.Bolt.Ok {
	//	dbs.autoLoadBolt()
	//}

	return nil
}

func (dbs *DatabaseDrivers) autoLoadGorm(isAutoMigrate bool, isTruncate bool, isAdminUI bool, tables ...interface{}) error {
	if err := AutoLoadGorm(dbs.Gorm.Cli, isAutoMigrate, isTruncate, isAdminUI, Tables...); err != nil {
		log.WithError(err).WithFields(
			logrus.Fields{"src.file": "model/data-core.go",
				"method.name": "New(...)",
				"db.type":     "sql",
				"db.driver":   "gorm",
				//"db.adpater": 		sqlEngine,
				"method.prev": "db.autoLoadGorm(...)",
				"prefix":      "dbs-new",
			}).Error("error while trying to init 'Gorm' database driver.")
		return err
	}
	return nil
}

//func (dbs *DatabaseDrivers) autoLoadBolt(isAutoMigrate bool, isTruncate bool, isAdminUI bool, tables ...interface{}) error {
//}

func AutoLoadGorm(db *gorm.DB, isAutoMigrate bool, isTruncate bool, isAdminUI bool, tables ...interface{}) error {
	AdminUI.Ok = isAdminUI
	if isAdminUI {
		AdminUI.Res = admin.New(&qor.Config{DB: db})
	}
	for _, table := range tables {
		if isTruncate {
			if err := db.DropTableIfExists(table).Error; err != nil {
				log.WithError(err).WithFields(
					logrus.Fields{"src.file": "model/data-core.go",
						"prefix":            "db-gorm",
						"method.name":       "(db *DatabaseDrivers) autoLoadGorm(...)",
						"method.prev":       "db.gormCli.DropTableIfExists(...)",
						"var.db.isTruncate": isTruncate,
						"var.db.table":      table,
					}).Warn("error while trying to drop an existing SQL table")
				return err
			}
		}
		if isAutoMigrate {
			if err := db.AutoMigrate(table).Error; err != nil {
				log.WithError(err).WithFields(
					logrus.Fields{"src.file": "model/data-core.go",
						"prefix":               "db-gorm",
						"method.name":          "(db *DatabaseDrivers) autoLoadGorm(...)",
						"method.prev":          "db.gormCli.AutoMigrate(...)",
						"var.db.isAutoMigrate": isAutoMigrate,
						"var.db.table":         table,
					}).Warn("error while trying to auto-migrate db table")
				return err
			}
		}
		if isAdminUI {
			AdminUI.Res.AddResource(table)
			log.WithFields(
				logrus.Fields{"src.file": "model/data-core.go",
					"method.name":        "(db *DatabaseDrivers) autoLoadGorm(...)",
					"method.prev":        "adminUI.AddResource(...)",
					"var.adminui.status": isAdminUI,
					"var.adminui.table":  table,
				}).Info("adding admin UI resource for the table")
		}
	}
	if isAdminUI {
		if len(AdminUI.Res.GetResources()) > 0 {
			for _, resource := range AdminUI.Res.GetResources() {
				log.WithFields(
					logrus.Fields{"src.file": "model/data-core.go",
						"method.name":          "(db *DatabaseDrivers) autoLoadGorm(...)",
						"method.prev":          "adminUI.GetResources()",
						"prefix":               "webui-admin",
						"var.adminui.resource": resource,
					}).Info("detected new admin UI resource")
			}
		}
	}
	return nil
}

func InitBoltDB(filePath string) (*bolt.DB, error) {
	// Get more config options to setup the bucket or the queue of tasks
	boltDB, err := bolt.Open(filePath, 0600, &bolt.Options{Timeout: time.Second})
	if err != nil {
		log.WithError(err).WithFields(
			logrus.Fields{"prefix": "db-boltdb",
				"method.name":       "InitBoltDB(...)",
				"method.prev":       "bolt.Open(...)",
				"db.adapter":        "boltdb",
				"src.file":          "model/data-core.go",
				"var.bolt.filepath": filePath,
				"var.bolt.options":  &bolt.Options{Timeout: time.Second},
			}).Warn("error while init the database with boltDB.")
		return boltDB, err
	}
	RootDrivers.Bolt.Ok = true
	RootDrivers.Bolt.Cli = boltDB
	return boltDB, nil
}

func InitEtcd(hosts []string, timeOut time.Duration) (*etcd.Client, error) {
	ectdConfig := etcd.Config{
		Endpoints:   hosts,
		DialTimeout: timeOut,
	}
	ectdClient, err := etcd.New(ectdConfig)
	if err != nil {
		log.WithError(err).WithFields(
			logrus.Fields{"prefix": "kvs-etcd",
				"method.name":  "InitEtcd(...)",
				"method.prev":  "etcd.New(...)",
				"db.adapter":   "etcd",
				"var.etcd.cfg": ectdConfig,
				"src.file":     "model/data-core.go",
			}).Warn("error while init the client connection with Etcd Key/Value store.")
		return ectdClient, err
	}
	defer ectdClient.Close()
	etcdKvc := etcd.NewKV(ectdClient)
	_, err = etcdKvc.Get(context.Background(), "/foo", nil)
	if err != nil {
		log.WithError(err).WithFields(
			logrus.Fields{"method.name": "InitEtcd(...)",
				"db.adapter":   "etcd",
				"prefix":       "kvs-etcd",
				"method.prev":  "etcdClient.Get(...)",
				"var.etcd.cfg": ectdConfig,
				"src.file":     "model/data-core.go",
			}).Warn("error while init the client connection with Etcd Key/Value store.")
		return ectdClient, err
	}
	RootDrivers.Etcd.Ok = true
	RootDrivers.Etcd.Cli = ectdClient
	RootDrivers.Etcd.Kvc = etcdKvc
	return RootDrivers.Etcd.Cli, nil
}

// https://github.com/skyrunner2012/xormplus/blob/master/xorm/dataset.go
// NewDataset creates a new Dataset.
func NewDataset(headers []string) *tablib.Dataset {
	return tablib.NewDataset(headers)
}

// NewStarDump(ds)
func NewDump(content []byte, dumpPrefixPath string, dumpType string, dataFormat []string) error {
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
		filePath := path.Join(dumpPrefixPath, dumpType+"."+t) // fmt.Sprintf("%s/%s", dumpPrefixPath, "repository.yaml") // will create a function
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
				return errors.New("Error while converting data to " + df + " format")
			}
			json.WriteTo(file)
			// log.WithFields(logrus.Fields{"method": "NewStarDump", "call": "WriteTo"}).Debugf("%#v Write to %#v",  df, filePath)
		case "yaml":
			yaml, err := ds.YAML()
			if err != nil {
				// panic(err)
				return errors.New("Error while converting data to " + df + " format")
			}
			yaml.WriteTo(file)
			// log.WithFields(logrus.Fields{"method": "NewStarDump", "call": "WriteTo"}).Debugf("%#v Write to %#v",  df, filePath)
		case "csv":
			csv, err := ds.CSV()
			if err != nil {
				// panic(err)
				return errors.New("Error while converting data to " + df + " format")
			}
			csv.WriteTo(file)
			// log.WithFields(logrus.Fields{"method": "NewStarDump", "call": "WriteTo"}).Debugf("%#v Write to %#v",  df, filePath)
		case "xml":
			xml, err := ds.XML()
			if err != nil {
				// panic(err)
				return errors.New("Error while converting data to " + df + " format")
			}
			xml.WriteTo(file)
			// log.WithFields(logrus.Fields{"method": "NewStarDump", "call": "WriteTo"}).Debugf("%#v Write to %#v",  df, filePath)
		case "markdown":
			ascii := ds.Tabular("markdown")
			if ascii == nil {
				// panic(err)
				return errors.New("Error while converting data to " + df + " format")
			}
			ascii.WriteTo(file)
			// log.WithFields(logrus.Fields{"method": "NewStarDump", "call": "WriteTo"}).Debugf("%#v Write to %#v",  df, filePath)
		default:
			return errors.New("Unsupported data format: " + df)
		}
		file.Close()
	}
	return nil
}

func TimeToMicroseconds(t time.Time) int64 {
	return t.Unix()*int64(time.Second/time.Microsecond) + int64(t.Nanosecond())/int64(time.Microsecond)
}
