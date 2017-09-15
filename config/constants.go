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

// locales
const clientLocaleDefault = "en_US.UTF-8"

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

// default [tag_format] pattern
const DEFAULT_TAG_FORMAT = `^.*\[\s*(.*)\s*\]\s*-\s*(.*)\r?\n?`

// for [type] section in rcfile
const TYPE_LIST_PATTERN = `\"(\S+)\"\s+=>\s+(\[\".*\"\])`

// format option
const (
	FORMAT_PRINT   = "print"
	FORMAT_JSON    = "json"
	FORMAT_UNITE   = "unite"
	FORMAT_SILENT  = "silent"
	FORMAT_NOCOLOR = "nocolor"
)

const (
	REMOTE_GITHUB    = "GITHUB"
	REMOTE_GITLAB    = "GITLAB"
	REMOTE_BITBUCKET = "BITBUCKET"
	REMOTE_ASANA     = "ASANA"
	REMOTE_NONE      = ""
)

// show option
const (
	SHOW_ALL   = "all"
	SHOW_CLEAN = "clean"
	SHOW_DIRTY = "dirty"
)

// parse depth limitation
const (
	PARSE_DEPTH_MIN = 1
	PARSE_DEPTH_MAX = 255
)

const (
	CONTEXT_LINES_MIN     = 1
	CONTEXT_LINES_MAX     = 255
	CONTEXT_LINES_DEFAULT = 15
)

// https://github.com/ykanda/watson-go/blob/master/constant.go
// setting file template
const LIMORC_TEMPLATE = `
# limo rc
# limo - inline issue manager
# [goosecode] labs

# Directories
[dirs]
./

# Tags
[tags]
fix
review
todo

# Ignores
[ignore]
.git
*.swp
`

const (
	// Time822 formt time for RFC 822
	Time822 = "02 Jan 2006 15:04:05 -0700" // "02 Jan 06 15:04 -0700"
)
