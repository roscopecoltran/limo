package model

import (
	"bytes"
	"encoding/gob"
	"errors"
	"io/ioutil"
	"strconv"
	"sync"
	"github.com/boltdb/bolt"
	"github.com/sirupsen/logrus"									// logs-logrus
)

// refs. 
//  - https://raw.githubusercontent.com/hfurubotten/autograder/master/database/database_test.go
//  - https://github.com/ssut/pocketnpm/blob/master/db/bolt_backend.go
//  - https://github.com/ssut/pocketnpm/blob/master/db/bolt_backend.go#L42-L55

var boltDB *bolt.DB
var registeredBucketNames = make([]string, 0)

/*
// InitBolt init bolt
func InitBolt() error {
	var err error
	// init Bolt DB
	Bolt, err = bolt.Open(Cfg.GetBoltFile(), 0600, nil)
	if err != nil {
		return err
	}
	// create buckets if not exists
	return Bolt.Update(func(tx *bolt.Tx) error {
		if _, err = tx.CreateBucketIfNotExists([]byte("koip")); err != nil {
			return err
		}
		return nil
	})
}
*/

// Start will start up the database. If the database does not already exist, a new one will be created.
func Start(dbloc string) (err error) {
	boltDB, err = bolt.Open(dbloc, 0666, nil)
	if err != nil {
		return err
	}
	return boltDB.Update(func(tx *bolt.Tx) (err error) {
		// Create a buckets.
		for _, bucket := range registeredBucketNames {
			if _, err = tx.CreateBucketIfNotExists([]byte(bucket)); err != nil {
				return err
			}
		}
		return nil
	})
}

// Store will put a new value in the assigned bucket(e.g. table) with given key.
//
// Key variable can be integer or string type.
// Value variable can be any type.
func StoreBolt(bucket string, key string, value interface{}) (err error) {
	return boltDB.Update(func(tx *bolt.Tx) (err error) {
		// open the bucket
		b := tx.Bucket([]byte(bucket))

		// Checks if the bucket was opened, and creates a new one if not existing. Returns error on any other situation.
		if b == nil {
			// Create a bucket.
			b, err = tx.CreateBucket([]byte(bucket))
			if err != nil {
				log.WithError(err).WithFields(
					logrus.Fields{	"prefix": 						"dbs-new",
									"src.file": 					"model/data-boltdb.go", 
									"db.adapter": 					"bolt",
									"db.driver": 					"boltdb", 
									"db.type": 						"kvs",
									"method.name": 					"StoreBolt(...)", 
									"method.prev": 					"tx.CreateBucket(...)",
									"var.bucket": 					bucket,
									"var.key": 						key,
									"var.value": 					value,
									}).Error("Error while creating the BoltDB bucket.")
				return err
			}
			if b == nil {
				log.WithFields(
					logrus.Fields{	"prefix": 						"dbs-new",
									"src.file": 					"model/data-boltdb.go", 
									"db.adapter": 					"bolt",
									"db.driver": 					"boltdb", 
									"db.type": 						"kvs",
									"method.name": 					"StoreBolt(...)", 
									"method.prev": 					"tx.CreateBucket(...)",
									"var.bucket": 					bucket,
									"var.key": 						key,
									"var.value": 					value,
									}).Error("Couldn't create bucket.")
				return errors.New("Couldn't create bucket.")
			}
		}

		defer UnlockBolt(bucket, key)

		buf := &bytes.Buffer{}
		encoder := gob.NewEncoder(buf)

		if err = encoder.Encode(value); err != nil {
			return
		}

		data, err := ioutil.ReadAll(buf)
		if err != nil {
			return err
		}

		return b.Put([]byte(key), data)
	})
}

// Get will get a value for the given key in a bucket(e.g. table).
func GetBolt(bucket string, key string, val interface{}, readonly bool) (err error) {
	return boltDB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return errors.New("Trying to access a nonexisting bucket.")
		}

		if !readonly {
			LockBolt(bucket, key)
		}

		data := b.Get([]byte(key))
		if data == nil {
			return errors.New("No data in database.")
		}

		buf := &bytes.Buffer{}
		decoder := gob.NewDecoder(buf)

		n, _ := buf.Write(data)

		if n != len(data) {
			return errors.New("Couldn't write all data to buffer while getting data from database. " + strconv.Itoa(n) + " != " + strconv.Itoa(len(data)))
		}

		return decoder.Decode(val)
	})
}

// Has will check if the key is pressent in the database.
func HasBolt(bucket, key string) bool {
	found := false

	err := boltDB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return errors.New("Unknown bucket")
		}
		c := b.Cursor()

		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			if key == string(k) {
				found = true
				break
			}
		}

		return nil
	})

	if err != nil {
		return false
	}

	return found
}

// Remove will delete a key in specified bucket.
func RemoveBolt(bucket, key string) (err error) {
	return boltDB.Update(func(tx *bolt.Tx) error {
		return tx.Bucket([]byte(bucket)).Delete([]byte(key))
	})
}

// RegisterBucket Will store all bucket names reserved by other packages. When
// the database is started these bucket names will be made sure exists in the boltDB.
func RegisterBoltBucket(bucket string) (err error) {
	registeredBucketNames = append(registeredBucketNames, bucket)
	if boltDB == nil {
		return
	}
	return boltDB.Update(func(tx *bolt.Tx) (err error) {
		_, err = tx.CreateBucketIfNotExists([]byte(bucket)) // Create a bucket.
		return err
	})
}

// GetPureDB Returns the pure connection to the database. Can be used with more
// advanced DB interaction.
func GetPureDB() *bolt.DB {
	if boltDB == nil {
		panic("Trying to obtain uninitalized database")
	}
	return boltDB
}

// Close will shut down the database in a safe mather.
func CloseBolt() (err error) {
	return boltDB.Close()
}

var writerslock sync.Mutex
var writerkeys = make(map[string]map[string]valueLocker)

type valueLocker struct {
	sync.Mutex
	islocked bool
}

// Lock will lock a specified key in a bucket for further use.
func LockBolt(bucket string, key string) {
	writerslock.Lock()
	defer writerslock.Unlock()

	if _, ok := writerkeys[bucket]; !ok {
		writerkeys[bucket] = make(map[string]valueLocker)
	}

	if _, ok := writerkeys[bucket][key]; !ok {
		writerkeys[bucket][key] = valueLocker{}
	}

	wkl := writerkeys[bucket][key]
	wkl.Lock()
	wkl.islocked = true
	writerkeys[bucket][string(key)] = wkl
}

// Unlock will unlock a specified key in a bucket and make it usable for other
// tasks running.
func UnlockBolt(bucket string, key string) {
	writerslock.Lock()
	defer writerslock.Unlock()

	if _, ok := writerkeys[bucket]; !ok {
		return
	}

	if _, ok := writerkeys[bucket][key]; !ok {
		return
	}

	wkl := writerkeys[bucket][key]
	if wkl.islocked {
		wkl.Unlock()
		wkl.islocked = false
	}
	writerkeys[bucket][key] = wkl
}