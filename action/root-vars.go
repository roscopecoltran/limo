package action

import (
	"github.com/roscopecoltran/sniperkit-limo/config" 												// app-config
	"github.com/roscopecoltran/sniperkit-limo/model" 												// data-models
	"github.com/sirupsen/logrus" 																	// logs-logrus
	//"github.com/davecgh/go-spew/spew" 															// debug-print
	//"github.com/k0kubun/pp" 																		// debug-print
)

var (
	configuration 				*config.Config 														// cfg-init
	dbs 						= &model.DatabaseDrivers{} 											// data-drivers
	log 						= logrus.New() 														// logs-logrus
)

var (
	inputFile, outputFile 		string 																// ai-word-embed
	dimension, window     		int 																// ai-word-embed
	learningRate          		float64 															// ai-word-embed
)

var options struct {
	config 						string
	language 					string
	output   					string
	service  					string
	tag      					string
	verbose  					bool
	dir 		struct {
		tmp 					string 
		data 					string
		conf 					string
	}
}