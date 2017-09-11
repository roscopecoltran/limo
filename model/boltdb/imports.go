package boltdb

import (
	"encoding/binary"
	"encoding/json"
	"github.com/boltdb/bolt"
	// "github.com/roscopecoltran/sniperkit-limo/model"
)

var importsBucket = []byte("imports")

type mapping map[string]map[string]int

type BucketRepository struct {
	driver *Driver
}

func NewDataRepository(driver *Driver) *BucketRepository {
	return &BucketRepository{
		driver: driver,
	}
}

func (r *BucketRepository) Get(userID int, source, ref string) (int, error) {
	var m mapping
	err := r.driver.bucket.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(importsBucket)

		data := bucket.Get(itob(userID))
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

