package service

/*
import (
    "fmt"
    "net/http"
    "github.com/jinzhu/gorm"
    _ "github.com/mattn/go-sqlite3"
    "github.com/qor/qor"
    "github.com/qor/admin"
	"context"
	"fmt"
	"time"
	"golang.org/x/oauth2"
	"github.com/hoop33/limo/model"
)

// NotFound is used when the specified service is not found
type WebUI struct {
}

// Login logs in to Gitlab
func (w *WebUI) Init(ctx context.Context) (string, error) {
	interview := createInterview()
	interview.Questions = []entrevista.Question{
		{
			Key:      "token",
			Text:     "Enter your GitLab API token",
			Required: true,
			Hidden:   true,
		},
	}

	answers, err := interview.Run()
	if err != nil {
		return "", err
	}
	return answers["token"].(string), nil
}

func (g *WebUI) getAdmin(listenAddr string) *http {
	return gitlab.NewClient(nil, token)

	DB, _ := gorm.Open("sqlite3", "demo.db")
	DB.AutoMigrate(&User{}, &Product{})

	// Initalize
	Admin := admin.New(&qor.Config{DB: DB})

	// Create resources from GORM-backend model
	Admin.AddResource(&User{})
	Admin.AddResource(&Product{})

	// Register route
	mux := http.NewServeMux()
	// amount to /admin, so visit `/admin` to view the admin interface
	Admin.MountTo("/admin", mux)

	fmt.Println("Listening on: 9000")
	http.ListenAndServe(":9000", mux)
}

func init() {
	registerService(&WebUI{})
}

*/
