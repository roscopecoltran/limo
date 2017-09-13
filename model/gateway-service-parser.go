package model

// Parser reads a configuration file, parses it and returns the content as an init ServiceConfig struct
type Gateway_Parser interface {
	Gateway_Parse(configFile string) (Gateway_ServiceConfig, error)
}