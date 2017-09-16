package boltdb

import (
    "errors"
	"time"
	// "github.com/roscopecoltran/sniperkit-limo/config"
	// "github.com/roscopecoltran/sniperkit-limo/model"
	"github.com/boltdb/bolt"
 	// "github.com/turtlemonvh/blanket-api"
 	// "github.com/turtlemonvh/blanket/worker"
 	// "github.com/alaska/boltqueue"
 	// "github.com/oliveagle/boltq"
    // "github.com/qor/qor"
    // "github.com/qor/admin"
	log "github.com/sirupsen/logrus"
)

type Driver struct {
	bucket *bolt.DB
}

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