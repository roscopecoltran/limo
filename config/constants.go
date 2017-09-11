package config

import (
	"os"
)

// app
const ProgramName 					= "limo" 												// ProgramName is the program name
const Version 						= "undefined" 											// Version is the semver-compliant program version

// vcs / git
const GitCommit 					= "undefined" 											// GITCOMMIT indicates which git hash the binary was built off of

// paths
const _sep 							= string(os.PathSeparator) 	 							// path separator


// cfg paths (can be overriden by the config.yaml file)
const DefaultConfigFormatExtension 	= "yml"
const DefaultConfigPrefixPath 		= "config"
const DefaultComponentsPrefixPath 	= "components"
const DefaultDatasetsPrefixPath 	= "datasets"

// config files per components
const GeneralConfigBase     		= "general" 		// cfg - general
const APIConfigBase     			= "api"				// cfg - api
const WebUIConfigBase     			= "webui" 			// cfg - webui
const SMTPConfigBase     			= "smtp" 			// cfg - smtp
const ActivityConfigBase   			= "activity" 		// cfg - activity
const ContainersConfigBase     		= "containers" 		// cfg - containers
const MacrosConfigBase     			= "macros" 			// cfg - macros
const TasksConfigBase     			= "tasks"			// cfg - tasks
const GatewayConfigBase     		= "gateway"			// cfg - gateway
const ErrorsConfigBase     			= "errors" 			// cfg - errors
const NotificationsConfigBase     	= "notifications" 	// cfg - notifications
const SearchEngineConfigBase     	= "search" 			// cfg - search
const ProvidersConfigBase     		= "providers" 		// cfg - providers
const EnginesConfigBase     		= "engines" 		// cfg - engines
const LocalesConfigBase     		= "locales" 		// cfg - locales
const KeywordsConfigBase     		= "keywords" 		// cfg - keywords
const PatternsConfigBase     		= "patterns" 		// cfg - patterns
const CorpusConfigBase     			= "corpus" 			// cfg - corpus




