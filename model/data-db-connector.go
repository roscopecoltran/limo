package model

// curl -s https://api.github.com/repos/chimeracoder/gojson | gojson -name=Repository

import (
    "errors"
	"time"
	"path"
	"os"
	// "path/filepath"


	/*
		Database/stores packages
	*/
	// gorm
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	// etcd
	etcd "github.com/coreos/etcd/client"

	// mongdb
	// "gopkg.in/mgo.v2"
	// "gopkg.in/mgo.v2/bson"

	// boltdb
	"github.com/boltdb/bolt"

	// graph

	// Neo4j
	"github.com/jmcvetta/neoism"

	// cayleyGraphDB
	"github.com/cayleygraph/cayley"
	"github.com/cayleygraph/cayley/graph"
	_ "github.com/cayleygraph/cayley/graph/bolt"
	"github.com/cayleygraph/cayley/quad"

	// data transform
	tablib "github.com/agrison/go-tablib"
	//"github.com/google/go-github/github"
	//"github.com/xanzy/go-gitlab"
	// jsoniter "github.com/json-iterator/go"
	// "github.com/davecgh/go-spew/spew"
	// beego
	// "github.com/astaxie/beego"
	// qor
	"golang.org/x/net/context"
    "github.com/qor/qor"
    "github.com/qor/admin"
	// search
	"github.com/blevesearch/bleve"
	// graphs
	"github.com/ckaznocha/taggraph"
	// logs
	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)


// datasets - formats
var validDataOutput 	= []string{"md","csv","yaml","json","xlsx","xml","tsv","mysql","postgres","html","ascii"}
var availableLocales 	= []string{"en-US", "fr-FR", "pl-PL"}
// var serviceConfig config.Config
// var cfg *config.Config

type EnhancedTime 		time.Time

type Databases struct {
	Datastore 			map[string]*bolt.DB
	Database   			map[string]*gorm.DB
	SearchIdx 			map[string]*bleve.Index
	KvIdx 				map[string]etcd.KeysAPI
}

var (
	Tables       	= []interface{}{
		&Service{}, 	&Category{}, 																				// service + registry organization
		&Star{}, 		&Readme{}, 		&WikiPage{}, 	&User{},													// vcs content indexation
		&Tag{}, 		&Topic{}, 		&Tree{}, 		&Language{}, 	&LanguageDetected{}, 	&LanguageType{}, 	// vcs repository classification
	}
)

var (
	adminUI 			*admin.Admin
	db 					DatabaseDrivers
	log 				= logrus.New()
	tagg 				= taggraph.NewTagGaph()
)

type DatabaseDrivers struct {
	boltCli  			*bolt.DB
	gormCli 			*gorm.DB
	bleveIdx 			*bleve.Index
	etcdCli  			etcd.KeysAPI
	neo4jCli 			*neoism.Database
	//dynamodbClient 	*dynamodb.DynamoDB
}

//	globalSetting := make(map[string]string)

// https://github.com/thesyncim/365a/blob/master/server/app.go
// https://github.com/emotionaldots/arbitrage/blob/master/cmd/arbitrage-db/main.go

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

	// db = &DatabaseDrivers{}
	InitDatabases() 		// databases

}

// New return a new db with config input
//func NewDatabases(conf config.Config) (db *Databases, err error) {

//}

func GetDatabases() (db *DatabaseDrivers, err error) {
	return db, nil
}

func InitDatabases() (db *DatabaseDrivers, err error) {
	if err := db.initGorm("sqlite3", "./shared/data/limo/limo.db", true); err != nil {
		log.WithError(err).WithFields(
			logrus.Fields{	"file": "data-db-connector.go", 
							"method_name": "Init", 
							"driver": "gorm", 
							"action": "InitGorm"}).Error("error while trying to init 'Gorm' database driver.")
		return nil, err
	}

	if err := db.autoloadGorm(true, true, true, Tables...); err != nil {
		log.WithError(err).WithFields(
			logrus.Fields{	"file": "data-db-connector.go", 
							"method_name": "Init", 
							"driver": "gorm", 
							"action": "InitGorm"}).Error("error while trying to init 'Gorm' database driver.")
		return nil, err
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
			logrus.Fields{	"file": "data-db-connector.go", 
							"method_name": "Init", 
							"driver": "boltdb", 
							"action": "InitBoltDB"}).Error("error while trying to init 'BoltDB' database driver")
		return nil, err
	}

	etcdDefaultHost 	:= []string{"http://127.0.0.1:2379"}
	etcdDefaultTimeout 	:= 1 * time.Second
	if err := db.initEtcd(etcdDefaultHost, etcdDefaultTimeout, true); err != nil {
		log.WithError(err).WithFields(
			logrus.Fields{	"file": "data-db-connector.go", 
							"method_name": "Init", 
							"driver": "etcd", 
							"action": "AutoloadDB"}).Error("error while trying to auto-load all program the tables")
		return nil, err
	}
	return db, nil
}

// https://github.com/qor/qor-example/blob/master/db/db.go
// InitDB initializes the database at the specified path
func InitDB(filepath string, adapter string, verbose bool) (*gorm.DB, error) {
//func InitDB(filepath string, adapter string, verbose bool) (*gorm.DB, error) {
	gormDB, err := gorm.Open(adapter, filepath) 	// Get more config options to setup the SQL database
	if err != nil {
		log.WithError(err).WithFields(
			logrus.Fields{	"file": "data-db-connector.go", 
							"method_name": "InitDB", 
							"adapter": adapter}).Warn("error while init the database with gorm.")
		return nil, err
	}
	gormDB.LogMode(verbose) 	// cfg.App.DebugMode
	/*
	// isAutoMigrate, isTruncate, isAdminUIResource
	if err := gormDB.AutoloadDB(true, true, false, Tables...); err != nil {
		log.WithError(err).WithFields(
			logrus.Fields{	"file": "data-db-connector.go", 
							"method_name": "InitDB", 
							"adapter": adapter, 
							"action": "AutoloadDB"}).Warn("error while trying to auto-load all program the tables")
	}
	*/
	return gormDB, nil
}

func InitBoltDB(filepath string) (*bolt.DB, error) {
	// Get more config options to setup the bucket or the queue of tasks
	boltDB, err := bolt.Open(filepath, 0600, &bolt.Options{Timeout: time.Second})
	if err != nil {
		log.WithError(err).WithFields(
			logrus.Fields{	"db": "InitDB", 
							"engine": "boltdb", 
							"filepath": filepath, 
							"bolt.Options":  &bolt.Options{Timeout: time.Second}}).Warnf("error while init the database with boltDB.")
		return nil, err
	}
	return boltDB, err
}

func InitEtcd(hosts []string, timeOut time.Duration, verbose bool) error {
	cfg := etcd.Config{
		Endpoints:               hosts,
		Transport:               etcd.DefaultTransport,
		HeaderTimeoutPerRequest: timeOut,
	}
	cli, err := etcd.New(cfg)
	if err != nil {
		return err
	}
	etcdClient := etcd.NewKeysAPI(cli)
	_, err = etcdClient.Get(context.Background(), "/foo", nil)
	if err != nil && err.Error() == etcd.ErrClusterUnavailable.Error() {
		return err
	}
	return nil
}

func TimeToMicroseconds(t time.Time) int64 {
	return t.Unix()*int64(time.Second/time.Microsecond) + int64(t.Nanosecond())/int64(time.Microsecond)
}

/*
type DatabaseDrivers struct {
	boltCli  			*bolt.DB
	gormCli 			*gorm.DB
	etcdCli  			etcd.KeysAPI
	bleveIdx 			bleve.Index
	//dynamodbClient 	*dynamodb.DynamoDB
}
*/

//func initGorm(db *gorm.DB) {
func (db *DatabaseDrivers) initGorm(filepath string, adapter string, verbose bool) (error) {
//(db *DatabaseDrivers) func initGorm(filepath string, adapter string, verbose bool) (error) {
	gormDB, err := gorm.Open(adapter, filepath) 	// Get more config options to setup the SQL database
	if err != nil {
		log.WithError(err).WithFields(
			logrus.Fields{	"file": "data-db-connector.go", 
							"method_name": "InstantiateGorm", 
							"adapter": adapter}).Warn("error while init the database with gorm.")
		return err
	}
	gormDB.LogMode(verbose) 	// cfg.App.DebugMode
	// isAutoMigrate, isTruncate, isAdminUIResource
	if err := db.autoloadGorm(true, true, false, Tables...); err != nil {
		log.WithError(err).WithFields(
			logrus.Fields{	"file": "data-db-connector.go", 
							"method_name": "InstantiateGorm", 
							"adapter": adapter, 
							"action": "AutoloadDB"}).Warn("error while trying to auto-load all program the tables")
		return err
	}
	db.gormCli = gormDB
	return nil
}

func (db *DatabaseDrivers) autoloadGorm(isAutoMigrate bool, isTruncate bool, isAdminUI bool, tables ...interface{}) (error) {
//(db *DatabaseDrivers) func autoloadGorm(isAutoMigrate bool, isTruncate bool, isAdminUIResource bool, tables ...interface{}) (error) {
	if isAdminUI {
		adminUI 	= 	admin.New(&qor.Config{DB: db.gormCli})
	}
	for _, table := range tables {
		if isTruncate {
			if err := db.gormCli.DropTableIfExists(table).Error; err != nil {
				log.WithError(err).WithFields(
					logrus.Fields{	"file": "data-db-connector.go", 
									"method_name": "AutoloadDB", 
									"action": "DropTableIfExists"}).Warnf("error while trying to drop the existing table: %#s.", table)
				return err
			}
		}
		if isAutoMigrate {
			if err := db.gormCli.AutoMigrate(table).Error; err != nil {
				log.WithError(err).WithFields(
					logrus.Fields{	"file": "data-db-connector.go", 
									"method_name": "AutoloadDB", 
									"action": "AutoMigrate"}).Warnf("error while trying to auto-migrate the table: %#s.", table)
				return err
			}
		}
		if isAdminUI {
			log.WithFields(
				logrus.Fields{	"file": "data-db-connector.go", 
								"method_name": "AutoloadDB", 
								"action": "AddResource"}).Infof("adding admin UI resource for the table: %#s.", table)
			adminUI.AddResource(table)
		}
	}
	if isAdminUI {
		if len(adminUI.GetResources()) > 0 {
			for _, resource := range adminUI.GetResources() {	
				log.WithFields(
					logrus.Fields{	"file": "data-db-connector.go", 
									"method_name": "AutoloadDB", 
									"resource": resource}).Infof("detected new admin UI resource: %#s.", resource)
			}		
		}
	}
	return nil
}

func (db *DatabaseDrivers) initBoltDB(filePath string, filePermissions os.FileMode, options *bolt.Options, verbose bool) (error) {
//(db *DatabaseDrivers) func initBoltDB(filePath string, filePermissions string, options *bolt.Options, verbose bool) (error) {
	boltDB, err := bolt.Open(filePath, filePermissions, options)
	if err != nil {
		log.WithError(err).WithFields(
			logrus.Fields{	"db": "InstantiateBoltDB", 
							"engine": "boltdb", 
							"filepath": filePath, 
							"bolt.Options":  options}).Warnf("error while init the database with boltDB.")
		return err
	}
	db.boltCli = boltDB
	return nil
}

func (db *DatabaseDrivers) initEtcd(hosts []string, timeout time.Duration, verbose bool) error {
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
							"cfg": cfg}).Warnf("error while init the database with boltDB.")
		return err
	}
	etcdClient := etcd.NewKeysAPI(cli)
	_, err = etcdClient.Get(context.Background(), "/sniperkit", nil)
	if err != nil && err.Error() == etcd.ErrClusterUnavailable.Error() {
		log.WithError(err).WithFields(
			logrus.Fields{	"db": "InstantiateEtcd", 
							"engine": "etcd", 
							"cfg": cfg}).Warnf("error while init the database with boltDB.")
		return err
	}
	db.etcdCli = etcdClient
	return nil
}

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

