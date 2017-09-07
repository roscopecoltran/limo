package config

import (
	"fmt"
	"io/ioutil"
	// "log"
	"os"
	"path"
	"github.com/cep21/xdgbasedir"
	//"github.com/jinzhu/configor"
	log "github.com/sirupsen/logrus"
	// prefixed "github.com/x-cray/logrus-prefixed-formatter"
	"gopkg.in/yaml.v2"
)

// 
// HeaderValue is the value of the custom Limo header
var HeaderValue = fmt.Sprintf("Version %s", Version)

// UserAgent is the value of the user agent header sent to the backends
var UserAgent = fmt.Sprintf("%s Version %s", ProgramName, Version)

// configuration path
var configDirectoryPath string

// ServiceConfig contains configuration information for a service
type ServiceConfig struct {
	Token string
	User  string
}

// OutputConfig sontains configuration information for an output
type OutputConfig struct {
	SpinnerIndex    int 	`yaml:"spinnerIndex"`
	SpinnerInterval int 	`yaml:"spinnerInterval"`
	SpinnerColor    string 	`yaml:"spinnerColor"`
}

// GetService returns the configuration information for a service
func (config *Config) GetService(name string) *ServiceConfig {
	if config.Services == nil {
		config.Services = make(map[string]*ServiceConfig)
	}
	service := config.Services[name]
	if service == nil {
		service = &ServiceConfig{}
		config.Services[name] = service
	}
	return service
}

// GetOutput returns the configuration information for an output
func (config *Config) GetOutput(name string) *OutputConfig {
	if config.Outputs == nil {
		config.Outputs = make(map[string]*OutputConfig)
	}
	output := config.Outputs[name]
	if output == nil {
		output = &OutputConfig{}
		config.Outputs[name] = output
	}
	return output
}

// ReadConfig reads the configuration information
func ReadConfig() (*Config, error) {
	file := configFilePath()
	var config Config
	if _, err := os.Stat(file); err == nil {
		// Read and unmarshal file only if it exists
		f, err := ioutil.ReadFile(file)
		if err != nil {
			return nil, err
		}
		err = yaml.Unmarshal(f, &config)
		if err != nil {
			return nil, err
		}
	}
	// Set default database path
	if config.DatabasePath == "" {
		config.DatabasePath = path.Join(configDirectoryPath, fmt.Sprintf("%s.db", ProgramName))
	}
	log.WithFields(log.Fields{"config": "ReadConfig"}).Infof("config.DatabasePath: %#v", config.DatabasePath)
	// Set default datastore path
	if config.DatastorePath == "" {
		config.DatastorePath = path.Join(configDirectoryPath, fmt.Sprintf("%s.boltdb", ProgramName))
	}
	log.WithFields(log.Fields{"config": "ReadConfig"}).Infof("config.DatastorePath: %#v", config.DatastorePath)

	// Set default search index path
	if config.IndexPath == "" {
		config.IndexPath = path.Join(configDirectoryPath, fmt.Sprintf("%s.idx", ProgramName))
	}
	log.WithFields(log.Fields{"config": "ReadConfig"}).Infof("config.IndexPath: %#v", config.IndexPath)
	return &config, nil
}

// WriteConfig writes the configuration information
func (config *Config) WriteConfig() error {
	err := os.MkdirAll(configDirectoryPath, 0700)
	if err != nil {
		return err
	}
	log.WithFields(log.Fields{"config": "WriteConfig"}).Infof("configDirectoryPath: %#v", configDirectoryPath)
	data, err := yaml.Marshal(config)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(configFilePath(), data, 0600)
}

func configFilePath() string {
	configFilePath := path.Join(configDirectoryPath, fmt.Sprintf("%s.yaml", ProgramName))
	log.WithFields(log.Fields{"config": "configFilePath"}).Infof("configDirectoryPath: %#v", configDirectoryPath)
	log.WithFields(log.Fields{"config": "configFilePath"}).Infof("configFilePath: %#v", configFilePath)
	return configFilePath
}

func init() {
	// Log as JSON instead of the default ASCII formatter.
	// log.SetFormatter(&log.JSONFormatter{})
	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)
	// Only log the warning severity or above.
	log.SetLevel(log.InfoLevel)
	baseDir, err := xdgbasedir.ConfigHomeDirectory()
	if err != nil {
		log.WithFields(log.Fields{"config": "init"}).Fatal("Can't find XDG BaseDirectory")
		// log.Fatal("Can't find XDG BaseDirectory")
	} else {
		configDirectoryPath = path.Join(baseDir, ProgramName)
	}
	log.WithFields(log.Fields{"config": "init"}).Infof("baseDir: %#v", baseDir)
	log.WithFields(log.Fields{"config": "init"}).Infof("configDirectoryPath: %#v", configDirectoryPath)
}
