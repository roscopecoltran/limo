package service

import (
	"context" 													// go-core
	"fmt" 														// go-core
	"reflect" 													// go-core
	"strings" 													// go-core
	"os" 														// go-core
	"github.com/fatih/color" 									// cli-output
	"github.com/hoop33/entrevista" 								// cli	
	"github.com/roscopecoltran/sniperkit-limo/model" 			// data
	"github.com/sirupsen/logrus" 								// logs
	prefixed "github.com/x-cray/logrus-prefixed-formatter" 		// logs
)

// https://github.com/cloudflavor/shep/blob/master/pkg/services/system.go
var	log 	= logrus.New()

func init() {
	// logs
	log.Out 		= 	os.Stdout
	formatter 		:= 	new(prefixed.TextFormatter)
	log.Formatter 	= 	formatter
	log.Level 		= 	logrus.DebugLevel
}

// Service represents a service
type Service interface {
	Login(ctx context.Context) (string, error)
	GetStars(ctx context.Context, starChan chan<- *model.StarResult, token, user string, isAugmented bool)
	GetTrending(ctx context.Context, trendingChan chan<- *model.StarResult, token, language string, verbose bool) 	// , isAugmented bool
	GetEvents(ctx context.Context, eventChan chan<- *model.EventResult, token, user string, page, count int)
}

// Add providers and engines here !

var services = make(map[string]Service)

func registerService(service Service) {
	services[Name(service)] = service
}

// Name returns the name of a service
func Name(service Service) string {
	parts 	:= strings.Split(reflect.TypeOf(service).String(), ".")
	name 	:= strings.ToLower(parts[len(parts)-1])
	log.WithFields(logrus.Fields{
		"service": 		"service/service.go",
		"prefix": 		"svc-core",
		"call": 		"Name(...)",
		"parts": 		parts,
		"name": 		name,
		}).Info("returning service name...")
	return name
}

// ForName returns the service for a given name, or an error if it doesn't exist
func ForName(name string) (Service, error) {
	if service, ok := services[strings.ToLower(name)]; ok {
		return service, nil
	}
	log.WithFields(logrus.Fields{
		"service": 		"service/service.go",
		"prefix": 		"svc-core",
		"call": 		"ForName(...)",
		"name": 		name,
		}).Errorf("Service '%s' not found", name)
	return &NotFound{}, fmt.Errorf("Service '%s' not found", name)
}

func createInterview() *entrevista.Interview {
	interview := entrevista.NewInterview()
	interview.ShowOutput = func(message string) {
		fmt.Print(color.GreenString(message))
	}
	interview.ShowError = func(message string) {
		color.Red(message)
	}
	return interview
}

func checkMapHasKey(entitiesMap map[string]interface{}, key string) string {
	if val, ok := entitiesMap[key]; ok {
		if val != nil {
			return val.(string)
		}
	}
	return ""
}

