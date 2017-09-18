package gorm

import (
	"github.com/roscopecoltran/sniperkit-limo/config" 							// app-config
	"github.com/roscopecoltran/sniperkit-limo/model" 							// app-model
	"github.com/jinzhu/gorm" 													// db-sql-gorm
	_ "github.com/jinzhu/gorm/dialects/sqlite" 									// db-sql-gorm-sqlite3
	_ "github.com/jinzhu/gorm/dialects/mysql" 									// db-sql-gorm-mysql
	_ "github.com/jinzhu/gorm/dialects/postgres" 								// db-sql-gorm-postgres
	_ "github.com/jinzhu/gorm/dialects/postgres" 								// db-sql-gorm-postgres
	"github.com/sirupsen/logrus"												// logs-logrus
	prefixed "github.com/x-cray/logrus-prefixed-formatter" 						// logs-logrus
)

const PKG_GORM_LABEL_CLUSTER 		= 		"dbs"
const PKG_GORM_LABEL_GROUP 			= 		"dbs-rdb"
const PKG_GORM_LABEL_PREFIX 		= 		"dbs-rdb-gorm"
const PKG_GORM_LABEL_DRIVER 		= 		"gorm"

var (
	dbs 							*model.RootDrivers
)

var (
	log 							= logrus.New()
	defaultGormFilePath 			= 		"./shared/data/limo/limo.db"
	defaultGormAdapter 				= 		"sqlite3"
)

func init() {
	log.Out 						= 		os.Stdout 							// logs-logrus
	formatter 						:= 		new(prefixed.TextFormatter) 		// logs-logrus
	log.Formatter 					= 		formatter 							// logs-logrus
	log.Level 						= 		logrus.DebugLevel 					// logs-logrus
}

// Gorm Resources
type GormRes 	struct {
	Ok 				bool  			`default:"false" json:"-" yaml:"-"`
	Cli 			*gorm.DB 		`json:"-" yaml:"-"`
}

func New(filepath string, adapter string) (*gorm.DB, error) {
	client, err := InitGorm(defaultGormFilePath, defaultGormAdapter)
	if err != nil {
		log.WithError(err).WithFields(
			logrus.Fields{	"prefix": 							PKG_GORM_LABEL_PREFIX,
							"group": 							PKG_GORM_LABEL_GROUP,
							"src.file": 						"model/drivers/rdb/gorm/client.go", 
							"method.name": 						"New(...)", 
							"method.prev": 						"db.initGorm(...)",
							"dbs.rdb.adapter": 					adapter, 
							"db.type": 							"sql", 
							"db.driver.name": 					"gorm", 
							"struct.dbs": 						dbs,
							"var.filePath": 					filePath, 
							"var.defaultGormFilePath": 			defaultGormFilePath, 
							"var.adpater": 						sqlEngine,
							}).Error("error while trying to init 'Gorm' database driver.")
		return nil, err
	}
	log.WithFields(logrus.Fields{	
		"prefix": 							PKG_GORM_LABEL_PREFIX,
		"group": 							PKG_GORM_LABEL_GROUP,
		"src.file": 						"model/drivers/rdb/gorm/client.go", 
		"method.name": 						"New(...)", 
		"method.prev": 						"db.initGorm(...)",
		"dbs.rdb.adapter": 					adapter, 
		"db.type": 							"sql", 
		"db.driver.name": 					"gorm", 
		"struct.dbs": 						dbs,
		"var.filePath": 					filePath, 
		"var.defaultGormFilePath": 			defaultGormFilePath, 
		"var.adpater": 						sqlEngine,
		}).Debug("status 'Gorm' drivers")
	return client, nil
}

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
