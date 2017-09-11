package config

import (
	"path"
	"os"
	"fmt"
)

// locales
var	clientLocale 					= "undefined"

// user-agents
var HeaderName 						= fmt.Sprintf("X-%X", ProgramName) 						// LimoHeaderName is the name of the custom KrakenD header
var HeaderValue 					= fmt.Sprintf("Version %s", Version) 					// HeaderValue is the value of the custom Limo header
var UserAgent 						= fmt.Sprintf("%s Version %s", ProgramName, Version)	// UserAgent is the value of the user agent header sent to the backends

// cfg - default
var DefaultConfigPath 				= path.Join(DefaultConfigPrefixPath, "config_default.yml")
var ConfigPath        				= path.Join(DefaultConfigPrefixPath, "config.yml")

var configPrefixPaths  				= [4]string{ path.Join("~", fmt.Sprintf("%s", ProgramName), fmt.Sprintf("%s", ProgramName)),
												 path.Join("shared", "conf.d", fmt.Sprintf("%s", ProgramName)),
												 path.Join("..", "shared", "conf.d", fmt.Sprintf("%s", ProgramName)),
												 path.Join("..", "..", "shared", "conf.d", fmt.Sprintf("%s", ProgramName))}

// tmp dir
var osTmpDir 						= os.TempDir()

// data 
var configFormats 					= [4]string{ "yaml", "json", "xml", "toml" }
var dumpDataFormats 				= [11]string{ "md", "csv", "yaml", "json", "xlsx", "xml", "tsv", "mysql", "postgres", "html", "ascii" }

// default paths
var GeneralConfigPath 				= path.Join(DefaultConfigPrefixPath, DefaultComponentsPrefixPath, fmt.Sprintf("%s.%s", GeneralConfigBase, DefaultConfigFormatExtension))
var APIConfigPath     				= path.Join(DefaultConfigPrefixPath, DefaultComponentsPrefixPath, fmt.Sprintf("%s.%s", APIConfigBase, DefaultConfigFormatExtension))
var WebUIConfigPath   				= path.Join(DefaultConfigPrefixPath, DefaultComponentsPrefixPath, WebUIConfigBase, ".", DefaultConfigFormatExtension)
var SMTPConfigPath    				= path.Join(DefaultConfigPrefixPath, DefaultComponentsPrefixPath, SMTPConfigBase, ".", DefaultConfigFormatExtension)
var ActivityConfigPath				= path.Join(DefaultConfigPrefixPath, DefaultComponentsPrefixPath, ActivityConfigBase, ".", DefaultConfigFormatExtension)
var ContainersConfigPath			= path.Join(DefaultConfigPrefixPath, DefaultComponentsPrefixPath, ContainersConfigBase, ".", DefaultConfigFormatExtension)
var MacrosConfigPath				= path.Join(DefaultConfigPrefixPath, DefaultComponentsPrefixPath, MacrosConfigBase, ".", DefaultConfigFormatExtension)

var KeywordsConfigPath 				= path.Join(DefaultConfigPrefixPath, DefaultComponentsPrefixPath, KeywordsConfigBase, ".", DefaultConfigFormatExtension)
var PatternsConfigPath 				= path.Join(DefaultConfigPrefixPath, DefaultComponentsPrefixPath, PatternsConfigBase, ".", DefaultConfigFormatExtension)
var CorpusConfigPath 				= path.Join(DefaultConfigPrefixPath, DefaultComponentsPrefixPath, CorpusConfigBase, ".", DefaultConfigFormatExtension)
var TasksConfigPath					= path.Join(DefaultConfigPrefixPath, DefaultComponentsPrefixPath, TasksConfigBase, ".", DefaultConfigFormatExtension)
var GatewayConfigPath   			= path.Join(DefaultConfigPrefixPath, DefaultComponentsPrefixPath, GatewayConfigBase, ".", DefaultConfigFormatExtension)
var ErrorsConfigPath   				= path.Join(DefaultConfigPrefixPath, DefaultComponentsPrefixPath, ErrorsConfigBase, ".", DefaultConfigFormatExtension)
var NotificationsConfigPath   		= path.Join(DefaultConfigPrefixPath, DefaultComponentsPrefixPath, NotificationsConfigBase, ".", DefaultConfigFormatExtension)
var SearchEngineConfigPath   		= path.Join(DefaultConfigPrefixPath, DefaultComponentsPrefixPath, SearchEngineConfigBase, ".", DefaultConfigFormatExtension)
var ProvidersConfigPath 			= path.Join(DefaultConfigPrefixPath, DefaultComponentsPrefixPath, ProvidersConfigBase, ".", DefaultConfigFormatExtension)
var EnginesConfigPath 				= path.Join(DefaultConfigPrefixPath, DefaultComponentsPrefixPath, EnginesConfigBase, ".", DefaultConfigFormatExtension)
var LocalesConfigPath 				= path.Join(DefaultConfigPrefixPath, DefaultComponentsPrefixPath, LocalesConfigBase, ".", DefaultConfigFormatExtension)
