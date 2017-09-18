package boltdb

import (
	//"os" 																		// go-core
	"github.com/roscopecoltran/sniperkit-limo/config" 							// app-config
	"github.com/roscopecoltran/sniperkit-limo/model" 							// app-model
	"github.com/boltdb/bolt" 													// dbs-kvs-boltdb
	"github.com/sirupsen/logrus"												// logs-logrus
	prefixed "github.com/x-cray/logrus-prefixed-formatter" 						// logs-logrus
)

const PKG_BOLTDB_LABEL_CLUSTER 		= 		"dbs"
const PKG_BOLTDB_LABEL_GROUP 		= 		"dbs-kvs"
const PKG_BOLTDB_LABEL_PREFIX 		= 		"dbs-kvs-boltdb"
const PKG_BOLTDB_LABEL_DRIVER 		= 		"boltdb"

var (
	dbs 							*model.RootDrivers
	log 							= 		logrus.New()
	defaultFilePermissions 			= 		0600
	defaultOptions 					= 		&bolt.Options{Timeout: time.Second}
)

func init() {
	log.Out 						= 		os.Stdout 							// logs-logrus
	formatter 						:= 		new(prefixed.TextFormatter) 		// logs-logrus
	log.Formatter 					= 		formatter 							// logs-logrus
	log.Level 						= 		logrus.DebugLevel 					// logs-logrus
}

// Bolt Resources
type BoltRes 	struct {
	Ok 				bool  			`default:"false" json:"-" yaml:"-"`
	Cli 			*bolt.DB 		`json:"-" yaml:"-"`
}

func Init(filePath string) (*bolt.DB, error) {
	client, err := bolt.Open(filePath, defaultFilePermissions, defaultOptions)
	if err != nil {
		log.WithError(err).WithFields(
			logrus.Fields{	"prefix": 							PKG_BOLTDB_LABEL_PREFIX,
							"group": 							PKG_BOLTDB_LABEL_GROUP,
							"src.file": 						"model/drivers/kvs/boltdb/client.go", 
							"method.name": 						"Init(...)", 
							"method.prev": 						"bolt.Open(...)",
							"dbs.adapter": 						PKG_BOLTDB_LABEL_DRIVER, 
							"dbs.driver.name": 					"bolt", 
							"var.defaultFilePermissions": 		defaultFilePermissions, 
							"var.bolt.filepath": 				filePath, 
							"var.bolt.options": 				defaultOptions,
							}).Error("error while init the database with boltDB.")
		return client, err
	}
	return client, nil
}

