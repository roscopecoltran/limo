package model

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"time"
	"github.com/sirupsen/logrus"
)

// https://github.com/toomore/lazyflickrgo/blob/master/simplecache/simplecache.go#L65-L72

// DataCacheRamDisk struct
type DataCacheRamDisk struct {
	Folder   		string					`json:"folder,omitempty" yaml:"folder,omitempty"`
	Dir      		string					`json:"dir,omitempty" yaml:"dir,omitempty"`	
	FullPath 		string 					`json:"full_path,omitempty" yaml:"full_path,omitempty"`	
	Expired  		time.Duration			`default:"3600" json:"debug,omitempty" yaml:"debug,omitempty"`	
}

// NewSimpleCache new a DataCacheRamDisk
func NewSniperkitCache(dir string, folder string, expired time.Duration) *DataCacheRamDisk {
	if dir == "" {
		dir = getOSRamdiskPath()
	}
	fullPath := filepath.Join(dir, folder)
	if err := os.Mkdir(fullPath, 0700); os.IsNotExist(err) {
		fullPath = filepath.Join(os.TempDir(), folder)
		os.Mkdir(fullPath, 0700)
		log.WithError(err).WithFields(
			logrus.Fields{
				"prefix": 				"data-cache",
				"src.file": 			"models/data-cache-ramdisk.go",
				"method.name": 			"NewSimpleCache(...)", 
				"method.prev": 			"os.Mkdir(...)", 
				"var.dir": 				dir,
				"var.fullPath": 		fullPath,
			}).Warn("creating a sniperkit cache directory...")
	}

	log.WithFields(
		logrus.Fields{
			"prefix": 				"data-cache",
			"src.file": 			"models/data-cache-ramdisk.go",
			"method.name": 			"NewSimpleCache(...)", 
			"method.prev": 			"filepath.Join(...)", 
			"var.dir": 				dir,
			"var.fullPath": 		fullPath,
		}).Info("initialize the sniperkit cache directory...")

	return &DataCacheRamDisk{
		Dir:      dir,
		Folder:   folder,
		Expired:  expired,
		FullPath: fullPath,
	}
}

// Get get cache
func (s *DataCacheRamDisk) Get(name string) ([]byte, error) {
	var err error
	if file, err := os.Open(filepath.Join(s.FullPath, name)); err == nil {
		defer file.Close()
		if stat, _ := file.Stat(); time.Now().Sub(stat.ModTime()) > s.Expired {
			log.WithFields(
				logrus.Fields{
					"prefix": 				"data-cache",
					"src.file": 			"models/data-cache-ramdisk.go",
					"method.name": 			"(s *DataCacheRamDisk) Get(...)", 
					"method.prev": 			"file.Stat(...)", 
					"var.name": 			name,
					"var.stat": 			stat,
					"var.s.Expired": 		s.Expired,
				}).Warn("Cache expired.")
			return nil, errors.New("Cache expired.")
		}
		return ioutil.ReadAll(file)
	}
	log.WithError(err).WithFields(
		logrus.Fields{
			"prefix": 				"data-cache",
			"src.file": 			"models/data-cache-ramdisk.go",
			"method.name": 			"(s *DataCacheRamDisk) Get(...)", 
			"method.prev": 			"os.Open(...)", 
			"var.name": 			name,
		}).Error("could not get the sniperkit cache directory...")
	return nil, err
}

// Set data
func (s *DataCacheRamDisk) Set(name string, data []byte) error {
	var err error
	if file, err := os.Create(filepath.Join(s.FullPath, name)); err == nil {
		defer file.Close()
		file.Write(data)
	}
	log.WithError(err).WithFields(
		logrus.Fields{
			"prefix": 				"data-cache",
			"src.file": 			"models/data-cache-ramdisk.go",
			"method.name": 			"(s *DataCacheRamDisk) Set(...)", 
			"method.prev": 			"os.Create(...)", 
			"var.name": 			name,
			"var.s.FullPath": 		s.FullPath,
		}).Error("could not set the sniperkit cache directory...")
	return err
}

func getOSRamdiskPath() string {

	ramDiskPath 	:= runtime.GOOS
	log.WithFields(
		logrus.Fields{
			"prefix": 				"data-cache",
			"src.file": 			"models/data-cache-ramdisk.go",
			"method.name": 			"getOSRamdiskPath(...)", 
			"method.prev": 			"runtime.GOOS", 
			"var.ramDiskPath": 		ramDiskPath,
		}).Info("detecting the operating system Ramdisk path...")

	switch ramDiskPath {
		//case "darwin":
		//	return "/run/shm/"
		case "linux":
			return "/run/shm/"
		default:
			return os.TempDir()
	}

}

