package config

import (
	"time"
)

var cfg *Config

var (
	DefaultConfigPath 			= "config/config_default.yml"
	ConfigPath        			= "config/config.yml"

	APIConfigPath     			= "config/components/api.yml"
	WebUIConfigPath   			= "config/components/webui.yml"
	SMTPConfigPath    			= "config/components/smtp.yml"

	LogingConfigPath			= "config/components/logging.yml"
	ContainersConfigPath		= "config/components/containers.yml"
	MacrosConfigPath			= "config/components/macros.yml"

	GatewayConfigPath   		= "config/components/gateway.yml"
	ErrorsConfigPath   			= "config/components/errors.yml"
	NotificationsConfigPath   	= "config/components/notifications.yml"
	SearchEngineConfigPath   	= "config/components/search.yml"
	TasksConfigPath   			= "config/components/tasks.yml"

	GeneralConfigPath 			= "config/components/general.yml"
	EnginesConfigPath 			= "config/components/engines.yml"
	LocalesConfigPath 			= "config/components/locales.yml"

	KeywordsConfigPath   		= "config/datasets/default_defs.yml"

)

type Config struct {

	DatabasePath 	string                    `yaml:"databasePath"`
	DatastorePath 	string                    `yaml:"datastorePath"`
	IndexPath    	string                    `yaml:"indexPath"`
	Services     	map[string]*ServiceConfig `yaml:"services"`
	Outputs     	map[string]*OutputConfig  `yaml:"outputs"`

	Debug   			bool			// run sniperkit-gateway api in debug mode

	App struct {
		Info 				AppProfileConfig
		Settings 			SettingsConfig
		FrontEnd 			WebUIConfig
		BackEnd 			WebUIConfig
		// Engines 				[]*EngineConfig `mapstructure:"engines"` // set of endpoint definitions
		Sentry 					SentryConfig
	}

	Auth AuthConfig

	WebUI struct {
		Default 	WebUIConfig
		FrontEnd 	WebUIConfig
		BackEnd 	WebUIConfig
	} 

	Gateway 		ServiceConfig 		`mapstructure:"service"`

	SearchEngine 	SearchEnginesConfig

	Database struct {

		Adapter  	string 				`env:"DBAdapter" default:"sqlite3"` 	// Options: mysql, postgres, sqlite3, mongodb
		Name     	string 				`env:"DBName" default:"qor_example"`	// Name of the database (which is a filepath for sqlite3
		Host     	string 				`env:"DBHost" default:"localhost"`		// For mysql, postgres, mongodb
		User     	string 				`env:"DBUser"`							// For mysql, postgres, mongodb
		Port     	string 				`env:"DBPort" default:"3306"`			// For mysql, postgres, mongodb
		Password 	string 				`env:"DBPassword"`						// For mysql, postgres, mongodb
		SSLMode		bool 				`default:"false"`						// For mysql, postgres

		Charset 	string 				`default:"utf8"`						// For mysql only
		ParseTime 	bool 				`default:"true"`						// For mysql only
		Local 		string 				`default:"Local"`						// For mysql only

		Mode 		string 				`default:"strong"`						// MongoDB only

		CreatedAt 	time.Time
		UpdatedAt 	time.Time

		SQLite 		SQLiteConfig
		BoltDB 		BoltDBConfig
		Graph 		GraphConfig

	}

	Dirs 				DirectoriesConfig
	Files 				FilesConfig
	Environment			EnvConfig    			
	Notifications 		NotificationConfig
	Logging     		LogConfig 				

	// Containers
	// Docker 				DockerConfig

}


// RegistryConfig represents information about the Registry of plugins, services and outputs
type RegistryConfig struct {
	Services     map[string]*ServiceConfig `yaml:"services"`
	Outputs      map[string]*OutputConfig  `yaml:"outputs"`
	// Modules      map[string]*ModuleConfig  	`yaml:"modules"`
	// Plugins      map[string]*PluginConfig  	`yaml:"plugins"`
	// Engines      map[string]*EngineConfig  	`yaml:"engines"`
	// Providers    map[string]*ProviderConfig  `yaml:"providers"`
	// Auths    	map[string]*AuthConfig  	`yaml:"auths"`
	// Patterns    	map[string]*PatternConfig  	`yaml:"patterns"`
	// Keywords    	map[string]*KeywordConfig  	`yaml:"keywords"`
	// Topics    	map[string]*TopicConfig  	`yaml:"topics"`
	// Analyzers    map[string]*AnalyzerConfig  `yaml:"analyzers"`
}

// OutputConfig sontains configuration information for an output
type OutputConfig2 struct {
	SpinnerIndex    int `json:"spinner_index,omitempty" yaml:"spinner_index,omitempty"`
	SpinnerInterval int `json:"spinner_interval,omitempty" yaml:"spinner_interval,omitempty"`
	SpinnerColor    string `json:"spinner_color,omitempty" yaml:"spinner_color,omitempty"`
}

// AppProfileConfig sontains configuration information for an app
type AppProfileConfig struct {
	// app name
	Name string `default:"app name"`
	// version code of the configuration
	Version int `mapstructure:"version"`
	// settings
	Settings 	SettingsConfig	
	// port
	Port     	int 				`env:"AppPort" default:"8080"`
	// host
	Host     	string 				`env:"AppHost" default:"0.0.0.0"`
	// host
	ListenAddr 	string 				`env:"AppListen" default:":8080"`
	// Contacts
	Contacts []struct {
		Name  string
		Email string `required:"true"`
	}	
}

// SettingsConfig sontains configuration information for sources of settings variables
type SettingsConfig struct {
	File 	SettingsFileConfig
	Env 	EnvConfig
}

// SettingsFileConfig sontains configuration information for a settings file
type SettingsFileConfig struct {
	Name 		string `default:"settings_default.yml" json:"filename,omitempty" yaml:"filename,omitempty"`
	PrefixPath 	string `default:"./conf.d" json:"prefix_path,omitempty" yaml:"prefix_path,omitempty"`
	Path 		string `default:"./conf.d/settings_default.yml" default:"10" json:"filepath,omitempty" yaml:"filepath,omitempty"`
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
	Active 				bool 	`default:"false"`
	Required 			bool 	`default:"false"`
	/*Github github.Config
	Gitlab github.Config
	Google google.Config*/
}

// WebUIConfig sontains configuration information for a webui service if not routed by the API service 
type WebUIConfig struct {
	Active 				bool 	`default:"false"`
	ListenAddr 			string 	`default:":8090" json:"port" yaml:"port"`			
	Host 				string 	`default:"0.0.0.0" json:"host,omitempty" yaml:"host,omitempty"`			
	Port 				int 	`default:"8090" json:"port,omitempty" yaml:"port,omitempty"`			
	Debug   			bool	// run sniperkit-gateway backend in debug mode
    StaticPath 			string 	`json:"static_path,omitempty" yaml:"static_path,omitempty"` // Custom static path - leave it blank if you didn't change
    TemplatesPath 		string 	`json:"templates_path,omitempty" yaml:"templates_path,omitempty"` // Custom templates path - leave it blank if you didn't change
    DefaultTheme 		string 	`default:":sniperkit" json:"default_theme,omitempty" yaml:"default_theme,omitempty"` // ui theme
    DefaultLocale 		string 	`json:"default_locale,omitempty" yaml:"default_locale,omitempty"` // Default interface locale - leave blank to detect from browser information or use codes from the 'locales' config section			
    HTTPS 				ServerConfig
    Security 			SecurityConfig
}

// Database implements database connection configuration details
type DBConnectionConfig struct {
	Adapter  	string
	Type string // Options: mysql, postgres, sqlite3, mongodb
	DB string // Name of the database (which is a filepath for sqlite3
	Host string // For mysql, postgres, mongodb
	User string // For mysql, postgres, mongodb
	Password string // For mysql, postgres, mongodb

	SSLMode bool // For postgres only

	Charset string // For mysql only
	ParseTime bool // For mysql only
	Local string // For mysql only

	Mode string // MongoDB only

	Settings 	DBSettingsConfig

}

// DBSettingsConfig represents information about a database service settings (cache, logs,..).
type DBSettingsConfig struct {
	Cache struct {
		MaxCacheSize  int         `json:"max_cache_size" yaml:"max_cache_size"`
		CacheLifetime int         `json:"cache_lifetime" yaml:"cache_lifetime"`
	}
	Logging LogConfig
}

type GraphConfig struct {
	Active 				bool 	`default:"false"`	
}

type DockerConfig struct {
	Containers 			[]ContainerConfig
}

type ContainerConfig struct {
	Active 				bool 	`default:"false"`	
    ContainerName 		string 	`json:"container_name" yaml:"container_name"` 
    Image 				string 	`required:"true" json:"image" yaml:"image"`
}

type ElasticSearchConfig struct {
	Active 				bool 	`default:"false"`	
}

type KibanaConfig struct {
	Active 				bool 	`default:"false"`	
	ProxiedBy 			bool 	`default:"nginx"`	
	ReverseProxy 		struct { 
		Configuration 	struct { 
			Nginx 		NginxConfig
			Caddy 		CaddyConfig
			Apache2 	Apache2Config
		}
	}
}

type LogStashConfig struct {
	Active 				bool 	`default:"false"`	
}

type FileBeatConfig struct {
	Active 				bool 	`default:"false"`	
}

type CaddyConfig struct {
	Active 				bool 	`default:"false"`	
}

type Apache2Config struct {
	Active 				bool 	`default:"false"`	
}

type NginxConfig struct {
	Active 				bool 	`default:"false"`	
}

type ElkConfig struct {
	Active 				bool 	`default:"false"`
	ElasticSearch 		ElasticSearchConfig	
	Kibana 				KibanaConfig
	FileBeat 			FileBeatConfig
	LogStash 			LogStashConfig
}

type SolrConfig struct {
	Active 				bool 	`default:"false"`	
}

type SphinxSearchConfig struct {
	Active 				bool 	`default:"false"`	
}

// SearchEnginesConfig represents information about a full-text search engine parameters.
type SearchEnginesConfig struct {
	Active 				bool 	`default:"false"`	
	Engines struct {
		ElasticSearch 	ElasticSearchConfig
		SphinxSearch 	SphinxSearchConfig
		Solr 			SolrConfig
	}
}

type ElkStackConfig struct {
	Elk ElkConfig
}

type BoltDBConfig struct {
	FilePath    string 		`env:"DBFilePath" default:"./shared/data/sniperkit.boltdb" json:"file_path,omitempty" yaml:"file_path,omitempty"`
	Buckets 	[]string 	`json:"buckets,omitempty" yaml:"buckets,omitempty"`
}

type SQLiteConfig struct {
	Version 	int `mapstructure:"version"`							// version code of the configuration
	FilePath    string `env:"DBFilePath" default:"./shared/data/sniperkit.sqlite" json:"file_path,omitempty" yaml:"file_path,omitempty"` // local filepath	
}

// config.Files.Extensions.Allowed
type FilesConfig struct {
	Active 				bool 			`default:"true"`
	Extensions struct {
		Allowed 	[]string 	`default:"go|py|md|cpp|h|php|java|cmake|txt|ini|yml|yaml|toml|ini|conf|log|js|html|htm|jx|jsx" json:"allowed,omitempty" yaml:"allowed,omitempty"`
		Blocked 	[]string 	`default:"epub|mobi|mp3|flac|mkv|avi|log" json:"forbidden,omitempty" yaml:"forbidden,omitempty"`
	} 
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

type LogConfig struct {

	Active 				bool 			`default:"true"`

	Access struct {
		AccessLogFilePath      	string `yaml:"access_log_filepath,omitempty"`
		AccessLogFileExtension 	string `yaml:"access_log_fileextension,omitempty"`
		AccessLogMaxSize       	int    `yaml:"access_log_max_size,omitempty"`
		AccessLogMaxBackups    	int    `yaml:"access_log_max_backups,omitempty"`
		AccessLogMaxAge        	int    `yaml:"access_log_max_age,omitempty"`
	}

	Error struct {
		ErrorLogFilePath       	string `yaml:"error_log_filepath,omitempty"`
		ErrorLogFileExtension  	string `yaml:"error_log_fileextension,omitempty"`
		ErrorLogMaxSize        	int    `yaml:"error_log_max_size,omitempty"`
		ErrorLogMaxBackups     	int    `yaml:"error_log_max_backups,omitempty"`
		ErrorLogMaxAge         	int    `yaml:"error_log_max_age,omitempty"`
	}

}

type DirectoriesConfig struct {
	Shared 				string `default:"shared" json:"shared_dir,omitempty" yaml:"shared_dir,omitempty"`
	Conf 				string `default:"conf.d" json:"conf_dir,omitempty" yaml:"conf_dir,omitempty"`
	Data 				string `default:"data" json:"conf_dir,omitempty" yaml:"data_dir,omitempty"`
	Load 				string `default:"load" json:"conf_dir,omitempty" yaml:"load_dir,omitempty"`
	Logs 				string `default:"logs" json:"conf_dir,omitempty" yaml:"logs_dir,omitempty"`
	Certs 				string `default:"certs" json:"certs_dir,omitempty" yaml:"certs_dir,omitempty"`
	Debug 				string `default:"debug" json:"debug_dir,omitempty" yaml:"debug_dir,omitempty"`
}

func GetConfig() *Config {
	return cfg
}

