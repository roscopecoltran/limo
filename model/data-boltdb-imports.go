package model

import (
	"encoding/binary"
	"encoding/json"
	"github.com/boltdb/bolt"
	"github.com/sirupsen/logrus"
)

var importsBucket = []byte("imports")

type mapping map[string]map[string]int

type BoltdbBucketRepository struct {
	databaseDrivers *DatabaseDrivers
}

func NewDataRepository(databaseDrivers *DatabaseDrivers) *BoltdbBucketRepository {
	return &BoltdbBucketRepository{
		databaseDrivers: databaseDrivers,
	}
}

func (r *DatabaseDrivers) Get(userID int, source, ref string) (int, error) {
	var m mapping
	err := r.boltCli.View(func(tx *bolt.Tx) error {
		boltCli := tx.Bucket(importsBucket)
		data := boltCli.Get(itob(userID))
		if data == nil {
			return nil
		}
		return json.Unmarshal(data, &m)
	})
	if err != nil {
		log.WithError(err).WithFields(
			logrus.Fields{	"prefix": 			"db-boltdb",
							"src.file": 		"model/data-boltdb-imports.go", 
							"method.name": 		"(r *DatabaseDrivers) Get(...)", 							
							"method.prev": 		"r.boltCli.View(...)",
							"var.userID": 		userID,
							"var.source": 		source,
							"var.ref": 			ref,
							}).Warn("error occured while trying to get some content from boltdb bucket")
		return 0, err
	}

	return m[source][ref], nil
}

func itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}

