package model

import (
	"time"
	"github.com/jinzhu/gorm"
	"github.com/qor/media"
	"github.com/qor/media/oss"
)

// User represents a GitHub user.
type User struct {
	gorm.Model
	// github
	ServiceID   	  uint 		 	`json:"service_id,omitempty" yaml:"service_id,omitempty"`
	Login             *string    	`json:"login,omitempty" yaml:"login,omitempty"`
	ID                *int       	`json:"id,omitempty" yaml:"id,omitempty"`
	AvatarURL         *string    	`json:"avatar_url,omitempty" yaml:"avatar_url,omitempty"`
	Avatar            AvatarImageStorage 
	HTMLURL           *string    	`json:"html_url,omitempty" yaml:"html_url,omitempty"`
	GravatarID        *string    	`json:"gravatar_id,omitempty" yaml:"gravatar_id,omitempty"`
	Name              *string    	`json:"name,omitempty" yaml:"name,omitempty"`
	Company           *string    	`json:"company,omitempty" yaml:"company,omitempty"`
	Blog              *string    	`json:"blog,omitempty" yaml:"blog,omitempty"`
	Location          *string    	`json:"location,omitempty" yaml:"location,omitempty"`
	Email             *string    	`json:"email,omitempty" yaml:"email,omitempty"`
	Hireable          *bool      	`json:"hireable,omitempty" yaml:"hireable,omitempty"`
	Bio               *string    	`json:"bio,omitempty" yaml:"bio,omitempty"`
	PublicRepos       *int       	`json:"public_repos,omitempty" yaml:"public_repos,omitempty"`
	PublicGists       *int       	`json:"public_gists,omitempty" yaml:"public_gists,omitempty"`
	Followers         *int       	`json:"followers,omitempty" yaml:"followers,omitempty"`
	Following         *int       	`json:"following,omitempty" yaml:"following,omitempty"`
	CreatedAt         time.Time 	`json:"created_at,omitempty" yaml:"created_at,omitempty"`
	UpdatedAt         time.Time 	`json:"updated_at,omitempty" yaml:"updated_at,omitempty"`
	SuspendedAt       time.Time 	`json:"suspended_at,omitempty" yaml:"suspended_at,omitempty"`
	Type              *string    	`json:"type,omitempty" yaml:"type,omitempty"`
	SiteAdmin         *bool      	`json:"site_admin,omitempty" yaml:"site_admin,omitempty"`
	TotalPrivateRepos *int       	`json:"total_private_repos,omitempty" yaml:"total_private_repos,omitempty"`
	OwnedPrivateRepos *int       	`json:"owned_private_repos,omitempty" yaml:"owned_private_repos,omitempty"`
	PrivateGists      *int       	`json:"private_gists,omitempty" yaml:"private_gists,omitempty"`
	DiskUsage         *int       	`json:"disk_usage,omitempty" yaml:"disk_usage,omitempty"`
	Collaborators     *int       	`json:"collaborators,omitempty" yaml:"collaborators,omitempty"`
	// extra
	Gender            string 		`default:"u" json:"gender,omitempty" yaml:"gender,omitempty"`
	Birthday          *time.Time 	`json:"birthday,omitempty" yaml:"birthday,omitempty"`
	Emails        	  []Email 		`gorm:"many2many:user_emails;" json:"emails,omitempty" yaml:"emails,omitempty"`
	Mentions          []Mention     `gorm:"many2many:user_mentions;"`
	Languages         []Language    `gorm:"many2many:user_languages;"`
}

type Email struct {
	gorm.Model
	email      	string
}


type UserResult struct {
	User  *User
	Error 	error
}

type AvatarImageStorage struct{ oss.OSS }

func (AvatarImageStorage) GetSizes() map[string]*media.Size {
	return map[string]*media.Size{
		"tiny":   {Width: 50, Height: 50},
		"small":  {Width: 120, Height: 120},
		"medium": {Width: 320, Height: 320},
		"big":    {Width: 460, Height: 460},
		"large":  {Width: 640, Height: 640},
		"xl":  	  {Width: 1000, Height: 1000},
	}
}

func (user User) DisplayName() string {
	return string(*user.Email)
}

func (user User) AvailableLocales() []string {
	return availableLocales
}

/*
func ExampleUsersService_ListAll() {
	client := github.NewClient(nil)
	opts := &github.UserListOptions{}
	for {
		users, _, err := client.Users.ListAll(opts)
		if err != nil {
			log.Fatalf("error listing users: %v", err)
		}
		if len(users) == 0 {
			break
		}
		opts.Since = *users[len(users)-1].ID
		// Process users...
	}
}
*/

/*
type User struct {
	gorm.Model
	Email                  string `form:"email"`
	Password               string
	Name                   string `form:"name"`
	Gender                 string
	Role                   string
	Birthday               *time.Time
	Balance                float32
	DefaultBillingAddress  uint `form:"default-billing-address"`
	DefaultShippingAddress uint `form:"default-shipping-address"`
	Addresses              []Address
	Avatar                 AvatarImageStorage

	// Confirm
	ConfirmToken string
	Confirmed    bool

	// Recover
	RecoverToken       string
	RecoverTokenExpiry *time.Time

	// Accepts
	AcceptPrivate bool `form:"accept-private"`
	AcceptLicense bool `form:"accept-license"`
	AcceptNews    bool `form:"accept-news"`
}

func (user User) DisplayName() string {
	return user.Email
}

func (user User) AvailableLocales() []string {
	return []string{"en-US", "zh-CN"}
}

type AvatarImageStorage struct{ oss.OSS }

func (AvatarImageStorage) GetSizes() map[string]*media.Size {
	return map[string]*media.Size{
		"small":  {Width: 50, Height: 50},
		"middle": {Width: 120, Height: 120},
		"big":    {Width: 320, Height: 320},
	}
}
*/