package gorm

import (
	//"github.com/jinzhu/gorm" 													// db-sql-gorm
	//_ "github.com/jinzhu/gorm/dialects/sqlite" 									// db-sql-gorm-sqlite3
	//_ "github.com/jinzhu/gorm/dialects/mysql" 									// db-sql-gorm-mysql
	//_ "github.com/jinzhu/gorm/dialects/postgres" 								// db-sql-gorm-postgres
	//_ "github.com/jinzhu/gorm/dialects/postgres" 								// db-sql-gorm-postgres
)

func init() {
}

func New() error {
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
}

func AutoLoad(db *gorm.DB, isAutoMigrate bool, isTruncate bool, isAdminUI bool, tables ...interface{}) error {
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