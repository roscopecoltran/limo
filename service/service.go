package service

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	log "github.com/sirupsen/logrus"
	"github.com/fatih/color"
	"github.com/hoop33/entrevista"
	"github.com/roscopecoltran/sniperkit-limo/model"
)

// Service represents a service
type Service interface {
	Login(ctx context.Context) (string, error)
	GetStars(ctx context.Context, starChan chan<- *model.StarResult, token, user string)
	// GetRepos(ctx context.Context, starChan chan<- *model.StarResult, token, user string)
	GetTrending(ctx context.Context, trendingChan chan<- *model.StarResult, token, language string, verbose bool)
	GetEvents(ctx context.Context, eventChan chan<- *model.EventResult, token, user string, page, count int)
}

var services = make(map[string]Service)

func registerService(service Service) {
	services[Name(service)] = service
}

// Name returns the name of a service
func Name(service Service) string {
	parts := strings.Split(reflect.TypeOf(service).String(), ".")
	name := strings.ToLower(parts[len(parts)-1])
	log.WithFields(log.Fields{"service": "Name"}).Infof("name: %#v", name)
	return name
}

// ForName returns the service for a given name, or an error if it doesn't exist
func ForName(name string) (Service, error) {
	if service, ok := services[strings.ToLower(name)]; ok {
		return service, nil
	}
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
