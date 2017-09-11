package model

import (

	// golang
    // "errors"
	// "time"

	// limo
	// "github.com/roscopecoltran/sniperkit-limo/config"

	// gorm
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	// mongdb
	// "gopkg.in/mgo.v2"

	// boltdb
	"github.com/boltdb/bolt"

	// claey + Boltbd
	//"github.com/cayleygraph/cayley"
	//"github.com/cayleygraph/cayley/graph"
	//_ "github.com/cayleygraph/cayley/graph/bolt"
	//"github.com/cayleygraph/cayley/quad"

	// beego
	// "github.com/astaxie/beego"

	// qor
    // "github.com/qor/qor"
    // "github.com/qor/admin"

	// logs
	"github.com/sirupsen/logrus"

)

// https://github.com/thesyncim/365a/blob/master/server/app.go
// https://github.com/emotionaldots/arbitrage/blob/master/cmd/arbitrage-db/main.go

// ErrNotGorm is used in case when the database type in the config file isn't a Gorm type of database
// var ErrNotGorm 		= errors.New("Not a Gorm database")
// var ErrNotMongoDB 	= errors.New("Not a MongoDB database")
// var ErrNotBoltDB 	= errors.New("Not a BoltDB key/value store")
// var ErrNotEtcd 		= errors.New("Not an Etcd key/value store")

// var db *dynamodb.DynamoDB

// var serviceConfig config.Config
// var cfg *config.Config

//type Databases struct {
//	Databases map[string]*bolt.DB
//	Indexes   map[string]*gorm.DB
//}

// https://github.com/qor/qor-example/blob/master/db/db.go

// InitDB initializes the database at the specified path
//func SyncDB(filepath string, verbose bool) (*gorm.DB, error) {
func SyncDB(db *gorm.DB, star *Star, service *Service) (bool, error) {
	// Get existing by remote ID and service ID
	var existing Star
	if db.Where("remote_id = ? AND service_id = ?", star.RemoteID, service.ID).First(&existing).RecordNotFound() {
		star.ServiceID = service.ID
		err := db.Create(star).Error
		if err != nil {
			log.WithError(err).WithFields(logrus.Fields{"config": "SyncDB", "star.RemoteID": star.RemoteID, "service.ID": service.ID}).Warnf("error while syncing the sql database service.")
		}
		return err == nil, err
	}
	star.ID = existing.ID
	star.ServiceID = service.ID
	star.CreatedAt = existing.CreatedAt
	return false, db.Save(star).Error
}

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

func SyncBuckets(bucket *bolt.DB, star *Star, service *Service) (bool, error) {
	return false, nil
}
