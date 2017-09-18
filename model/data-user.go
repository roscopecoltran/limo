package model

import (
	"github.com/jinzhu/gorm"
	"time"
	//"github.com/qor/sorting"
	"github.com/qor/media"
	"github.com/qor/media/oss"
	"github.com/sirupsen/logrus"
)

// User represents a GitHub user.
type User struct {
	gorm.Model `json:"-" yaml:"-"`
	//sorting.SortingDESC
	ServiceID              uint               `gorm:"column:service_id" json:"service_id,omitempty" yaml:"service_id,omitempty"`
	OpenID                 int                `gorm:"column:open_id" json:"open_id" yaml:"open_id"`
	Login                  *string            `gorm:"column:login" json:"login,omitempty" yaml:"login,omitempty"`
	Role                   string             `gorm:"column:role" json:"role,omitempty" yaml:"role,omitempty"`
	AvatarURL              *string            `gorm:"column:avatar_url" json:"avatar_url,omitempty" yaml:"avatar_url,omitempty"`
	Avatar                 AvatarImageStorage `json:"-" yaml:"-"`
	HTMLURL                *string            `gorm:"column:html_url" json:"html_url,omitempty" yaml:"html_url,omitempty"`
	GravatarID             *string            `gorm:"column:gravatar_id" json:"gravatar_id,omitempty" yaml:"gravatar_id,omitempty"`
	Name                   *string            `gorm:"column:name" json:"name,omitempty" yaml:"name,omitempty"`
	Company                *string            `gorm:"column:company" json:"company,omitempty" yaml:"company,omitempty"`
	Blog                   *string            `gorm:"column:blog" json:"blog,omitempty" yaml:"blog,omitempty"`
	Location               *string            `gorm:"column:location" json:"location,omitempty" yaml:"location,omitempty"`
	Languages              []Language         `gorm:"many2many:user_languages;"`
	Email                  string             `gorm:"column:email" json:"email,omitempty" yaml:"email,omitempty"`
	Position               string             `gorm:"column:position" json:"position,omitempty" yaml:"position,omitempty"`
	Hireable               *bool              `gorm:"column:hireable" json:"hireable,omitempty" yaml:"hireable,omitempty"`
	Bio                    *string            `gorm:"column:bio" json:"bio,omitempty" yaml:"bio,omitempty"`
	Balance                float64            `gorm:"column:balance" json:"balance,omitempty" yaml:"balance,omitempty"`
	DefaultBillingAddress  uint               `gorm:"column:default_billing_address" form:"default-billing-address" json:"default_billing_address,omitempty" yaml:"default_billing_address,omitempty"`
	DefaultShippingAddress uint               `gorm:"column:default_shipping_address" form:"default-shipping-address" json:"default_shipping_address,omitempty" yaml:"default_shipping_address,omitempty"`
	ConfirmToken           string             `gorm:"column:confirm_token" json:"confirm_token,omitempty" yaml:"confirm_token,omitempty"`
	Confirmed              bool               `gorm:"column:confirmed" json:"confirmed,omitempty" yaml:"confirmed,omitempty"`
	RecoverToken           string             `gorm:"column:recover_token" json:"recover_token,omitempty" yaml:"recover_token,omitempty"`
	RecoverTokenExpiry     *time.Time         `gorm:"column:recover_token_expiry" json:"recover_token_expiry,omitempty" yaml:"recover_token_expiry,omitempty"`
	// Accepts
	AcceptPrivate bool `gorm:"column:accept_private" form:"accept-private" json:"accept_private,omitempty" yaml:"accept_private,omitempty"`
	AcceptLicense bool `gorm:"column:accept_license" form:"accept-license" json:"accept_license,omitempty" yaml:"accept_license,omitempty"`
	AcceptNews    bool `gorm:"column:accept_news" form:"accept-news" json:"accept_news,omitempty" yaml:"accept_news,omitempty"`
	// extra
	Gender            string      `default:"unkown" gorm:"column:gender" json:"gender,omitempty" yaml:"gender,omitempty"`
	Addresses         []Address   `gorm:"many2many:user_physical_address;" json:"physical_address,omitempty" yaml:"physical_address,omitempty"`
	AlternativeEmails []Email     `gorm:"many2many:user_alternative_emails;" json:"alternative_emails,omitempty" yaml:"alternative_emails,omitempty"`
	Birthday          *time.Time  `gorm:"column:birthday" json:"birthday,omitempty" yaml:"birthday,omitempty"`
	Vcs               UserInfoVCS `gorm:"many2many:user_vcs_infos;" json:"vcs_infos,omitempty" yaml:"vcs_infos,omitempty"`
	// Authoring       		[]Mention   	`gorm:"many2many:user_authoring;" json:"authoring,omitempty" yaml:"authoring,omitempty"`
	// Languages         	[]Language    	`gorm:"many2many:user_languages;"`
}

type Email struct {
	gorm.Model
	Email           string        `gorm:"column:email" json:"email,omitempty" yaml:"email,omitempty"`
	EmailType       string        `default:"pro" gorm:"column:email_type" json:"email_type,omitempty" yaml:"email_type,omitempty"`
	Disabled        bool          `default:"false" gorm:"column:disabled" json:"disabled,omitempty" yaml:"disabled,omitempty"`
	RefServiceID    uint          `gorm:"column:ref_service_id" json:"ref_service_id,omitempty" yaml:"ref_service_id,omitempty"`
	RefServiceURL   string        `gorm:"column:ref_service_url" json:"ref_service_url,omitempty" yaml:"ref_service_url,omitempty"`
	RefExternalUrls []ExternalURL `gorm:"column:ref_external_urls" json:"ref_external_urls,omitempty" yaml:"ref_external_urls,omitempty"`
}

type UserResult struct {
	User  *User
	Error error
}

type AvatarImageStorage struct{ oss.OSS }

func (AvatarImageStorage) GetSizes() map[string]*media.Size {
	return map[string]*media.Size{
		"tiny":   {Width: 50, Height: 50},
		"small":  {Width: 120, Height: 120},
		"medium": {Width: 320, Height: 320},
		"big":    {Width: 460, Height: 460},
		"large":  {Width: 640, Height: 640},
		"xl":     {Width: 1000, Height: 1000},
	}
}

func (user User) DisplayName() string {
	return string(*user.Name)
}

func (user User) DisplayLogin() string {
	return string(*user.Login)
}

func (user User) DisplayEmail() string {
	return string(user.Email)
}

func (user User) AvailableLocales() []string {
	return availableLocales
}

// CreateOrUpdateStar creates or updates a star and returns true if the star was created (vs updated)
func CreateOrUpdateUser(db *gorm.DB, user *User, service *Service) (bool, error) {
	// Get existing by remote ID and service ID
	var existing User
	if db.Where("email = ? AND service_id = ?", user.Email, service.ID).First(&existing).RecordNotFound() {
		user.ServiceID = service.ID
		err := db.Create(user).Error
		log.WithError(err).WithFields(
			logrus.Fields{
				"method.name": "CreateOrUpdateUser(...)",
				"src.file":    "model/data-user.go",
				"prefix":      "user-profile",
				"var.user":    user,
			}).Error("database error")
		return err == nil, err
	}
	user.ID = existing.ID
	user.ServiceID = service.ID
	user.CreatedAt = existing.CreatedAt
	return false, db.Save(user).Error
}
