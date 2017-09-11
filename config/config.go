package config

import (
	"fmt"
	"io/ioutil"
	"time"
	"strings"
	"os"
	"path"
	"flag"
	"github.com/spf13/viper"
	"path/filepath"
	"github.com/cep21/xdgbasedir"
	//"github.com/jinzhu/configor"
	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
	"gopkg.in/yaml.v2"

)

var cfg 	*Config 				// config
var	log 	= logrus.New() 			// logs

var configDirectoryPath string 													// configuration path
var flagConfigPath 		= flag.String("config", "", "Path to look for a config file. (directory)")

func init() {
	// Log as JSON instead of the default ASCII formatter.
	// log.SetFormatter(&log.JSONFormatter{})
	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.Out = os.Stdout

	// Only log the warning severity or above.
	log.Level = logrus.InfoLevel

	formatter := new(prefixed.TextFormatter)
	log.Formatter = formatter

	baseDir, err := xdgbasedir.ConfigHomeDirectory()
	if err != nil {
		log.WithFields(logrus.Fields{"config": "init"}).Fatal("Can't find XDG BaseDirectory")
		// log.Fatal("Can't find XDG BaseDirectory")
	}

	configDirectoryPath = path.Join("shared", "conf.d", ProgramName)

	//}
	log.WithFields(logrus.Fields{"config": "init"}).Infof("baseDir: %#v", baseDir)
	log.WithFields(logrus.Fields{"config": "init"}).Infof("configDirectoryPath: %#v", configDirectoryPath)
}

type ServiceConfig struct {			// ServiceConfig contains configuration information for a service
	Token string
	User  string
}

type Config struct {

	DatabasePath 	string                    	`default:"./shared/data/limo/gorm" json:"database_path,omitempty" yaml:"database_path,omitempty"`
	DatastorePath 	string                    	`default:"./shared/data/limo/boltdb" json:"datastore_path,omitempty" yaml:"datastore_path,omitempty"`
	IndexPath    	string                    	`default:"./shared/data/limo/bleve" json:"index_path,omitempty" yaml:"index_path,omitempty"`
	Services     	map[string]*ServiceConfig 	`json:"services,omitempty" yaml:"services,omitempty"`
	Outputs     	map[string]*OutputConfig  	`json:"outputs,omitempty" yaml:"outputs,omitempty"`

	Debug   		bool						// run sniperkit-gateway api in debug mode

	App struct {
		Name 		string 						`default:"Sniperkit App" json:"name,omitempty" yaml:"name,omitempty" `
		Debug 		bool 						`default:"false" json:"debu,omitemptyg" yaml:"debug,omitempty"`
		Version 	string 						`default:"dev" json:"version,omitempty" yaml:"version,omitempty"`
		Config struct {
			FilePath string 					`default:"./shared/conf.d/limo/limo.yaml" json:"file_path,omitempty" yaml:"file_path,omitempty"`
			Paths 	[]string    				`json:"paths,omitempty" yaml:"paths,omitempty"`
			Formats []string					`json:"formats,omitempty" yaml:"formats,omitemptys"`
			Write 	bool 						`default:"true" json:"write,omitempty" yaml:"write,omitempty"`
			Print 	bool 						`default:"false" json:"print,omitempty" yaml:"print,omitempty"`
		}
		Performances struct {
			Parallelism int  					`default:"20" json:"parallelism,omitempty" yaml:"parallelism,omitempty"`
		}
	}

	Aggregate struct {
		Port 		int 						`default:"8000" json:"port,omitempty" yaml:"port,omitempty"`
		Debug 		bool 						`default:"false" json:"debug,omitempty" yaml:"debug,omitempty"`
		ListenAddr 	string 						`default:":8000" json:"listen_addr,omitempty" yaml:"listen_addr,omitempty"`
		Config struct {
			FilePath string 					`default:"./shared/conf.d/krakend/configuration.json" json:"file_path,omitempty" yaml:"file_path,omitempty"`
		}
	}

	Auth AuthConfig

	WebUI struct {
		Default 	WebUIConfig
		FrontEnd 	WebUIConfig
		BackEnd 	WebUIConfig
	} 

	// Gateway 		ServiceConfig 				`mapstructure:"service" json:"gateway,omitempty" yaml:"gateway,omitempty"`

	// SearchEngine 	SearchEnginesConfig

	Database struct {

		Adapter  	string 						`env:"DBAdapter" default:"sqlite3" json:"adapter,omitempty" yaml:"adapter,omitempty"` 	// Options: mysql, postgres, sqlite3, mongodb
		Name     	string 						`env:"DBName" default:"sniperkit-limo" json:"name,omitempty" yaml:"name,omitempty"`	// Name of the database (which is a filepath for sqlite3
		Host     	string 						`env:"DBHost" default:"localhost" json:"host,omitempty" yaml:"host,omitempty"`		// For mysql, postgres, mongodb
		User     	string 						`env:"DBUser" json:"user,omitempty" yaml:"user,omitempty"`							// For mysql, postgres, mongodb
		Port     	string 						`env:"DBPort" default:"3306" json:"port,omitempty" yaml:"port,omitempty"`			// For mysql, postgres, mongodb
		Password 	string 						`env:"DBPassword" json:"password,omitempty" yaml:"password,omitempty"`						// For mysql, postgres, mongodb
		SSLMode		bool 						`default:"false" json:"ssl_mode,omitempty" yaml:"ssl_mode,omitempty"`						// For mysql, postgres

		Charset 	string 						`default:"utf8" json:"charset,omitempty" yaml:"charset,omitempty"`						// For mysql only
		ParseTime 	bool 						`default:"true" json:"parse_time,omitempty" yaml:"parse_time,omitempty"`						// For mysql only
		Local 		string 						`default:"Local" json:"local,omitempty" yaml:"local,omitempty"`						// For mysql only

		Mode 		string 						`default:"strong" json:"mode,omitempty" yaml:"mode,omitempty"`						// MongoDB only

		CreatedAt 	time.Time 					`json:"created_at,omitempty" yaml:"created_at,omitempty"`
		UpdatedAt 	time.Time 					`json:"updated_at,omitempty" yaml:"updated_at,omitempty"`

		SQLite 		SQLiteConfig 				`json:"sqlite_params,omitempty" yaml:"sqlite_params,omitempty"`
		BoltDB 		BoltDBConfig 				`json:"boltdb_params,omitempty" yaml:"boltdb_params,omitempty"`
		Graph 		GraphConfig 				`json:"graphdb_params,omitempty" yaml:"graphdb_params,omitempty"`

	}

	Dirs 				DirectoriesConfig 		`json:"dirs,omitempty" yaml:"dirs,omitempty"`
	Files 				FilesConfig 			`json:"files,omitempty" yaml:"files,omitempty"`
	Environment			EnvConfig    			`json:"env,omitempty" yaml:"env,omitempty"`
	Notifications 		NotificationConfig 		`json:"notifications,omitempty" yaml:"notifications,omitempty"`
	Logging     		LogConfig 				`json:"logging,omitempty" yaml:"logging,omitempty"`

	// Docker 				DockerConfig 		`json:"containers,omitempty" yaml:"containers,omitempty"` 	// Containers

}


// RegistryConfig represents information about the Registry of plugins, services and outputs
type RegistryConfig struct {
	Services     map[string]*ServiceConfig 		`json:"services,omitempty" yaml:"services,omitempty"`
	Outputs      map[string]*OutputConfig  		`json:"outputs,omitempty" yaml:"outputs,omitempty"`
	// Modules      map[string]*ModuleConfig  	`json:"modules" yaml:"modules"`
	// Plugins      map[string]*PluginConfig  	`json:"plugins" yaml:"plugins"`
	// Engines      map[string]*EngineConfig  	`json:"engines" yaml:"engines"`
	// Providers    map[string]*ProviderConfig  `json:"providers" yaml:"providers"`
	// Auths    	map[string]*AuthConfig  	`json:"auths" yaml:"auths"`
	// Patterns    	map[string]*PatternConfig  	`json:"patterns" yaml:"patterns"`
	// Keywords    	map[string]*KeywordConfig  	`json:"keywords" yaml:"keywords"`
	// Topics    	map[string]*TopicConfig  	`json:"topics" yaml:"topics"`
	// Analyzers    map[string]*AnalyzerConfig  `json:"analyzers" yaml:"analyzers"`
}

// OutputConfig sontains configuration information for an output
type OutputConfig struct {
	SpinnerIndex    int 						`json:"spinner_index,omitempty" yaml:"spinner_index,omitempty"`
	SpinnerInterval int 						`json:"spinner_interval,omitempty" yaml:"spinner_interval,omitempty"`
	SpinnerColor    string 						`json:"spinner_color,omitempty" yaml:"spinner_color,omitempty"`
}

// AppProfileConfig sontains configuration information for an app
type AppProfileConfig struct {
	Name 			string 			`default:"app name"`				// app name
	Version 		int 			`mapstructure:"version"`			// version code of the configuration
	Settings 		SettingsConfig										// settings	
	Port     		int 			`env:"AppPort" default:"8080"` 		// port	
	Host     		string 			`env:"AppHost" default:"0.0.0.0"` 	// host
	ListenAddr 		string 			`env:"AppListen" default:":8080"` 	// listen address
	// Contacts
	Contacts 	[]struct {												
		Name  		string
		Email 		string 			`required:"true"`
	}	
}

// SettingsConfig sontains configuration information for sources of settings variables
type SettingsConfig struct {
	File 	SettingsFileConfig
	Env 	EnvConfig
}

// SettingsFileConfig sontains configuration information for a settings file
type SettingsFileConfig struct {
	Name 			string 				`default:"settings_default.yml" json:"filename,omitempty" yaml:"filename,omitempty"`
	PrefixPath 		string 				`default:"./conf.d" json:"prefix_path,omitempty" yaml:"prefix_path,omitempty"`
	Path 			string 				`default:"./conf.d/settings_default.yml" default:"10" json:"filepath,omitempty" yaml:"filepath,omitempty"`
}

// EnvConfig sontains configuration information for environement variables to load during the runtime
type EnvConfig struct {
	Active   			bool			`default:true"`
	Variables  			[]string 		`long:"variables" description:"load env var(s) VAR" value-name:"VAR" json:"-" yaml:"-"`
	Files 				[]string 		`long:"files" description:"load env file(s) FILE" value-name:"FILE" json:"-" yaml:"-"`
}

// TimeoutConfig sontains configuration information for timeouts for most of the services/actions of the app
type TimeoutConfig struct {
	Read  				time.Duration 	`default:"10" json:"read_timeout" yaml:"read_timeout"`
	Write 				time.Duration 	`default:"10" json:"write_timeout" yaml:"write_timeout"`
}

// SecurityConfig sontains configuration information for outbound and inbound network traffic
type SecurityConfig struct {
	Hosts struct {
	    Allowed			[]string 		`json:"allowed,omitempty" yaml:"allowed,omitempty"` // if the user has more than one network interface
	}
	Outgoing 			OutgoingConfig
}

// ApiConfig sontains configuration information for the api networking parameters
type ServerConfig struct {				
	Active   			bool			`default:true"`
	Host 				string 			`default:"0.0.0.0" json:"host,omitempty" yaml:"host,omitempty"` // port to bind the frontend (results) service			
	Port 				int 			`default:"8000" json:"port,omitempty" yaml:"port,omitempty"` // port to bind the frontend (results) service			
	ListenAddr   		string        	`default:":8000" json:"listen_addr" yaml:"listen_addr"`
	Timeout 			TimeoutConfig
}

type SentryConfig struct {
	Active 				bool 			`default:"false"`
	DSN 				string 			`required:"false" yaml:"sentry_dsn"`
}

type OutgoingConfig struct {				 		// communication with search engines
    RequestTimeout 		time.Duration  	`default:"2.0" json:"request_timeout" yaml:"request_timeout"` // seconds
    UserAgentSuffix 	string 	 		`json:"useragent_suffix,omitempty" yaml:"useragent_suffix,omitempty"` // suffix of searx_useragent, could contain informations like an email address to the administrator
    PoolConnections 	int 	 		`default:"100" json:"pool_connections" yaml:"pool_connections"` // Number of different hosts
    PoolMaxsize 		int 	 		`default:"10" json:"pool_maxsize" yaml:"pool_maxsize"` // Number of simultaneous requests by host
    Proxies				[]string 		// SOCKS proxies are also supported: see http://docs.python-requests.org/en/master/user/advanced/#socks
    SourceIps			[]string 		`yaml:"source_ips"` // if the user has more than one network interface
}

type AuthConfig struct {
	Active 				bool 			`default:"false"`
	Required 			bool 			`default:"false"`
}

// WebUIConfig sontains configuration information for a webui service if not routed by the API service 
type WebUIConfig struct {
	Active 				bool 			`default:"false"`
	ListenAddr 			string 			`default:":8090" json:"port" yaml:"port"`			
	Host 				string 			`default:"0.0.0.0" json:"host,omitempty" yaml:"host,omitempty"`			
	Port 				int 			`default:"8090" json:"port,omitempty" yaml:"port,omitempty"`			
	Debug   			bool			// run sniperkit-gateway backend in debug mode
    StaticPath 			string 			`json:"static_path,omitempty" yaml:"static_path,omitempty"` // Custom static path - leave it blank if you didn't change
    TemplatesPath 		string 			`json:"templates_path,omitempty" yaml:"templates_path,omitempty"` // Custom templates path - leave it blank if you didn't change
    DefaultTheme 		string 			`default:":sniperkit" json:"default_theme,omitempty" yaml:"default_theme,omitempty"` // ui theme
    DefaultLocale 		string 			`json:"default_locale,omitempty" yaml:"default_locale,omitempty"` // Default interface locale - leave blank to detect from browser information or use codes from the 'locales' config section			
    HTTPS 				ServerConfig
    Security 			SecurityConfig
}

// Database implements database connection configuration details
type DBConnectionConfig struct {
	Adapter  	string
	Type 		string 					// Options: mysql, postgres, sqlite3, mongodb
	DB 			string 					// Name of the database (which is a filepath for sqlite3
	Host 		string 					// For mysql, postgres, mongodb
	User 		string 					// For mysql, postgres, mongodb
	Password 	string 					// For mysql, postgres, mongodb
	SSLMode 	bool 					// For postgres only
	Charset 	string 					// For mysql only
	ParseTime 	bool 					// For mysql only
	Local 		string 					// For mysql only
	Mode 		string 					// MongoDB only
	Settings 	DBSettingsConfig

}

// DBSettingsConfig represents information about a database service settings (cache, logs,..).
type DBSettingsConfig struct {
	Cache struct {
		MaxCacheSize  	int         	`json:"max_cache_size" yaml:"max_cache_size"`
		CacheLifetime 	int         	`json:"cache_lifetime" yaml:"cache_lifetime"`
	}
	Logging LogConfig
}

type GraphConfig struct {
	Active 				bool 			`default:"false"`	
}

type BoltDBConfig struct {
	FilePath    	string 				`env:"DBFilePath" default:"./shared/data/sniperkit.boltdb" json:"file_path,omitempty" yaml:"file_path,omitempty"`
	Buckets 		[]string 			`json:"buckets,omitempty" yaml:"buckets,omitempty"`
}

type SQLiteConfig struct {
	Version 		int 				`mapstructure:"version"`							// version code of the configuration
	FilePath    	string 				`env:"DBFilePath" default:"./shared/data/sniperkit.sqlite" json:"file_path,omitempty" yaml:"file_path,omitempty"` // local filepath	
}

// config.Files.Extensions.Allowed
type FilesConfig struct {
	Active 			bool 				`default:"true"`
	//Extensions struct {
	//	Allowed 	[]string 	`default:"go|py|md|cpp|h|php|java|cmake|txt|ini|yml|yaml|toml|ini|conf|log|js|html|htm|jx|jsx" json:"allowed,omitempty" yaml:"allowed,omitempty"`
	//	Blocked 	[]string 	`default:"epub|mobi|mp3|flac|mkv|avi|log" json:"forbidden,omitempty" yaml:"forbidden,omitempty"`
	//} 
}

type NotificationConfig struct {
	Email struct {
		SMTP   SMTPConfig
	}
}

type SMTPConfig struct {
	Host     string
	Port     string
	User     string
	Password string
}

type DirectoriesConfig struct {
	Shared 						string 			`default:"shared" json:"shared_dir,omitempty" yaml:"shared_dir,omitempty"`
	Conf 						string 			`default:"conf.d" json:"conf_dir,omitempty" yaml:"conf_dir,omitempty"`
	Data 						string 			`default:"data" json:"conf_dir,omitempty" yaml:"data_dir,omitempty"`
	Load 						string 			`default:"load" json:"conf_dir,omitempty" yaml:"load_dir,omitempty"`
	Logs 						string 			`default:"logs" json:"conf_dir,omitempty" yaml:"logs_dir,omitempty"`
	Certs 						string 			`default:"certs" json:"certs_dir,omitempty" yaml:"certs_dir,omitempty"`
	Debug 						string 			`default:"debug" json:"debug_dir,omitempty" yaml:"debug_dir,omitempty"`
}

func GetTmpDir() (string) {
	fmt.Println("OS default temp is " + osTmpDir)
	return osTmpDir
}

func GetConfig() *Config {
	return cfg
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
	log.WithFields(logrus.Fields{"config": "ReadConfig"}).Infof("config.DatabasePath: %#v", config.DatabasePath)
	// Set default datastore path
	if config.DatastorePath == "" {
		config.DatastorePath = path.Join(configDirectoryPath, fmt.Sprintf("%s.boltdb", ProgramName))
	}
	log.WithFields(logrus.Fields{"config": "ReadConfig"}).Infof("config.DatastorePath: %#v", config.DatastorePath)

	// Set default search index path
	if config.IndexPath == "" {
		config.IndexPath = path.Join(configDirectoryPath, fmt.Sprintf("%s.idx", ProgramName))
	}
	log.WithFields(logrus.Fields{"config": "ReadConfig"}).Infof("config.IndexPath: %#v", config.IndexPath)
	return &config, nil
}

// WriteConfig writes the configuration information
func (config *Config) WriteConfig() error {
	err := os.MkdirAll(configDirectoryPath, 0700)
	if err != nil {
		return err
	}
	log.WithFields(logrus.Fields{"config": "WriteConfig"}).Infof("configDirectoryPath: %#v", configDirectoryPath)
	data, err := yaml.Marshal(config)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(configFilePath(), data, 0600)
}

func FindLocalConfig() (string) {
//func (config *Config) FindLocalConfig() (string) {
	configFilePath := configFilePath()
	if configFilePath == "" {
		log.WithFields(logrus.Fields{"config": "getConfiguration"}).Info("error while getting global configuration data.")
	}
	return configFilePath
}

// findLocalConfig returns the path to the local config file.
// It searches the current directory and all parent directories for a config file.
// If no config file is found, findLocalConfig returns an empty string.
func configFilePath() string {
	curdir, err := os.Getwd()
	if err != nil {
		curdir = "."
	}
	log.WithFields(logrus.Fields{"config": "configFilePath", "curdir": curdir}).Infof("currendir: %#v", curdir)
	path, err := filepath.Abs(curdir)
	if err != nil || path == "" {
		return ""
	}
	log.WithFields(logrus.Fields{"config": "configFilePath", "path": path}).Infof("path: %#v", path)
	//lp := ""
	for _, cfgPrefixPath := range configPrefixPaths {
		log.WithFields(logrus.Fields{"config": "configFilePath", "cfgPrefixPath": cfgPrefixPath}).Infof("path: %#v", path)
		for _, cfgFormat := range configFormats {
			confpath := filepath.Join(path, cfgPrefixPath, fmt.Sprintf("%s.%s", ProgramName, cfgFormat))
			// log.WithFields(logrus.Fields{"config": "configFilePath", "cfgPrefixPath": cfgPrefixPath, "cfgFormat": cfgFormat}).Infof("confpath: %#v", confpath)
			if _, err := os.Stat(confpath); err == nil {
				log.WithFields(logrus.Fields{"config": "configFilePath", "confpath": confpath, "cfgFormat": cfgFormat, "cfgPrefixPath": cfgPrefixPath}).Infof("FOUND: %#v", confpath)
				return confpath
			}

			// lp = path
			// path = filepath.Dir(path)
			// log.WithFields(logrus.Fields{"config": "configFilePath", "path": path}).Infof("path: %#v", path)
		}
	}

	return ""
}

func configFilePath2() string {
	configFilePath := path.Join(configDirectoryPath, fmt.Sprintf("%s.yaml", ProgramName))
	log.WithFields(logrus.Fields{"config": "configFilePath"}).Infof("configDirectoryPath: %#v", configDirectoryPath)
	log.WithFields(logrus.Fields{"config": "configFilePath"}).Infof("configFilePath: %#v", configFilePath)
	return configFilePath
}

func setConfig(path string) {
	// Default values
	viper.SetDefault("host.listen", "")
	viper.SetDefault("host.port", "4242")
	viper.SetDefault("host.hook", "hook")

	viper.SetDefault("repo.url", "https://github.com/roscopecoltran/sniperkit.git")
	viper.SetDefault("repo.path", "/shared/vcs/sniperkit")
	viper.SetDefault("repo.branch", "master")
	viper.SetDefault("repo.synccycle", 3600)

	viper.SetDefault("etcd.hosts", []string{"http://127.0.0.1:2379"})

	viper.SetDefault("auth.type", "ssh")
	viper.SetDefault("auth.ssh.key", "~/.ssh/id_rsa")
	viper.SetDefault("auth.ssh.public", "~/.ssh/id_rsa.pub")

	// Getting config from file
	viper.SetConfigName("sniperkit")
	viper.AddConfigPath("/etc/sniperkit/")
	viper.AddConfigPath("$HOME/.sniperkit")
	viper.AddConfigPath(".")
	if len(path) > 0 {
		viper.AddConfigPath(path)
	}
	err := viper.ReadInConfig()
	if err != nil {
		log.WithError(err).Warn("Couldn't read config file. Will use defaults.")
	}

	// Setting environment config
	viper.SetEnvPrefix("snk")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
}

