package service

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"os"
	"github.com/sirupsen/logrus"
	"github.com/fatih/color"
	"github.com/hoop33/entrevista"
	"github.com/roscopecoltran/sniperkit-limo/model"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
	"github.com/deckarep/gosx-notifier"
)

// https://github.com/cloudflavor/shep/blob/master/pkg/services/system.go

var	log 	= logrus.New()

func init() {

	// logs
	log.Out = os.Stdout
	// log.Formatter = new(prefixed.TextFormatter)

	formatter := new(prefixed.TextFormatter)
	formatter.FullTimestamp = true

	// Set specific colors for prefix and timestamp
	formatter.SetColorScheme(&prefixed.ColorScheme{
		PrefixStyle:    "blue+b",
		TimestampStyle: "white+h",
	})

	log.Formatter = formatter

}

func gosxnotifierTest() {
    //At a minimum specifiy a message to display to end-user.
    note := gosxnotifier.NewNotification("Check your Apple Stock!")

    //Optionally, set a title
    note.Title = "It's money making time 💰"

    //Optionally, set a subtitle
    note.Subtitle = "My subtitle"

    //Optionally, set a sound from a predefined set.
    note.Sound = gosxnotifier.Basso

    //Optionally, set a group which ensures only one notification is ever shown replacing previous notification of same group id.
    note.Group = "com.unique.yourapp.identifier"

    //Optionally, set a sender (Notification will now use the Safari icon)
    note.Sender = "com.apple.Safari"

    //Optionally, specifiy a url or bundleid to open should the notification be
    //clicked.
    note.Link = "http://www.yahoo.com" //or BundleID like: com.apple.Terminal

    //Optionally, an app icon (10.9+ ONLY)
    note.AppIcon = "gopher.png"

    //Optionally, a content image (10.9+ ONLY)
    note.ContentImage = "gopher.png"

    //Then, push the notification
    err := note.Push()

    //If necessary, check error
    if err != nil {
        log.Println("Uh oh!")
    }
}

// Service represents a service
type Service interface {
	Login(ctx context.Context) (string, error)
	// GetStars(ctx context.Context, starChan chan<- *model.StarResult, token, user string)
	GetStars(ctx context.Context, starChan chan<- *model.StarResult, token, user string, subChannels bool, subChannelsJobs uint)
	// GetUserInfos(ctx context.Context, starChan chan<- *model.StarResult, token)
	// GetReadmes(ctx context.Context, starChan chan<- *model.StarResult, token)
	// GetReadmes(ctx context.Context, starChan chan<- *model.StarResult, token, user string, name string)
	// GetRepos(ctx context.Context, starChan chan<- *model.StarResult, user string, name string)
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
	log.WithFields(logrus.Fields{"service": "Name"}).Infof("name: %#v", name)
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

func checkMapHasKey(entitiesMap map[string]interface{}, key string) string {
	if val, ok := entitiesMap[key]; ok {
		if val != nil {
			return val.(string)
		}
	}
	return ""
}



