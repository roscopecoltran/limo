package model

import (
	"encoding/binary"
	"encoding/json"
	"github.com/boltdb/bolt"
	// "github.com/roscopecoltran/sniperkit-limo/model"
)

var importsBucket = []byte("imports")

type mapping map[string]map[string]int

type BucketRepository struct {
	databaseDrivers *DatabaseDrivers
}

func NewDataRepository(databaseDrivers *DatabaseDrivers) *BucketRepository {
	return &BucketRepository{
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
		return 0, err
	}

	return m[source][ref], nil
}

func itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}

