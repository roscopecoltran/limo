package model

import (
	"errors"
	"github.com/jinzhu/gorm"
	"strconv"
	"time"
	//"github.com/qor/sorting"
	"github.com/google/go-github/github"
	"github.com/olivere/nullable"
	"github.com/sirupsen/logrus"
)

type UserInfoVCS struct {
	gorm.Model `json:"-" yaml:"-"`
	//sorting.SortingDESC
	ServiceID         uint               `gorm:"column:service_id" json:"service_id,omitempty" yaml:"service_id,omitempty"`
	UserID            string             `gorm:"column:user_id" json:"user_id,omitempty" yaml:"user_id,omitempty"`
	Login             *string            `gorm:"column:username" json:"username,omitempty" yaml:"username,omitempty"`
	LoginEmail        string             `gorm:"column:login_email" json:"login_email,omitempty" yaml:"login_email,omitempty"`
	AvatarURL         *string            `gorm:"column:avatar_url" json:"avatar_url,omitempty" yaml:"avatar_url,omitempty"`
	Avatar            AvatarImageStorage `json:"-" yaml:"-"`
	HTMLURL           *string            `gorm:"column:html_url" json:"html_url,omitempty" yaml:"html_url,omitempty"`
	GravatarID        *string            `gorm:"column:gravatar_id" json:"gravatar_id,omitempty" yaml:"gravatar_id,omitempty"`
	Name              *string            `gorm:"column:name" json:"name,omitempty" yaml:"name,omitempty"`
	Birthday          *time.Time         `gorm:"column:birthday" json:"birthday,omitempty" yaml:"birthday,omitempty"`
	Company           *string            `gorm:"column:company" json:"company,omitempty" yaml:"company,omitempty"`
	Blog              *string            `gorm:"column:blog" json:"blog,omitempty" yaml:"blog,omitempty"`
	Location          *string            `gorm:"column:location" json:"location,omitempty" yaml:"location,omitempty"`
	Hireable          *bool              `gorm:"column:hireable" json:"hireable,omitempty" yaml:"hireable,omitempty"`
	Bio               *string            `gorm:"column:bio" json:"bio,omitempty" yaml:"bio,omitempty"`
	PublicRepos       int                `gorm:"column:public_repos" json:"public_repos,omitempty" yaml:"public_repos,omitempty"`
	PublicGists       int                `gorm:"column:public_gists" json:"public_gists,omitempty" yaml:"public_gists,omitempty"`
	Followers         int                `gorm:"column:followers" json:"followers,omitempty" yaml:"followers,omitempty"`
	Following         int                `gorm:"column:following" json:"following,omitempty" yaml:"following,omitempty"`
	CreatedAt         time.Time          `gorm:"column:user_created_at" json:"user_created_at,omitempty" yaml:"user_created_at,omitempty"`
	UpdatedAt         time.Time          `gorm:"column:user_updated_at" json:"user_updated_at,omitempty" yaml:"user_updated_at,omitempty"`
	SuspendedAt       time.Time          `gorm:"column:user_suspended_at" json:"user_suspended_at,omitempty" yaml:"user_suspended_at,omitempty"`
	AccountType       *string            `gorm:"column:account_type" json:"account_type,omitempty" yaml:"account_type,omitempty"`
	SiteAdmin         *bool              `gorm:"column:site_admin" json:"site_admin,omitempty" yaml:"site_admin,omitempty"`
	TotalPrivateRepos int                `gorm:"column:total_private_repos" json:"total_private_repos,omitempty" yaml:"total_private_repos,omitempty"`
	OwnedPrivateRepos int                `gorm:"column:owned_private_repos" json:"owned_private_repos,omitempty" yaml:"owned_private_repos,omitempty"`
	PrivateGists      int                `gorm:"column:private_gists" json:"private_gists,omitempty" yaml:"private_gists,omitempty"`
	DiskUsage         *int               `gorm:"column:disk_usage" json:"disk_usage,omitempty" yaml:"disk_usage,omitempty"`
	Collaborators     int                `gorm:"column:collaborators" json:"collaborators,omitempty" yaml:"collaborators,omitempty"`
	AlternativeEmails []Email            `gorm:"many2many:alternative_emails;" json:"alternative_emails,omitempty" yaml:"alternative_emails,omitempty"`
	ProgLanguages     []Language         `gorm:"many2many:user_programming_languages;" json:"user_programming_languages,omitempty" yaml:"user_programming_languages,omitempty"`
}

func NewUserFromGithub(userVCS *github.User) (*UserInfoVCS, error) {

	// Require the GitHub UserID
	if userVCS.ID == nil {
		errMsg := errors.New("ID from GitHub User is required")
		log.WithError(errMsg).WithFields(
			logrus.Fields{
				"method.name": "NewUserFromGithub(...)",
				"src.file":    "model/vcs-user.go",
				"prefix":      "vs-github-user",
				"var.userVCS": userVCS,
			}).Error("missing identifier")
		return nil, errMsg
	}

	userAccountEmail := nullable.StringWithDefault(userVCS.Email, "hidden@github.com") // Set 'public_repos' count to 0 if nil

	publicReposCount := nullable.IntWithDefault(userVCS.PublicRepos, 0)             // Set 'public_repos' count to 0 if nil
	publicGistsCount := nullable.IntWithDefault(userVCS.PublicGists, 0)             // Set 'public_gists' count to 0 if nil
	followersCount := nullable.IntWithDefault(userVCS.Followers, 0)                 // Set 'followers' count to 0 if nil
	followingCount := nullable.IntWithDefault(userVCS.Following, 0)                 // Set 'following' count to 0 if nil
	collaboratorsCount := nullable.IntWithDefault(userVCS.Collaborators, 0)         // Set 'collaborators' count to 0 if nil
	totalPrivateReposCount := nullable.IntWithDefault(userVCS.TotalPrivateRepos, 0) // Set 'total_private_repos' count to 0 if nil
	ownedPrivateReposCount := nullable.IntWithDefault(userVCS.OwnedPrivateRepos, 0) // Set 'owned_private_repos' count to 0 if nil
	privateGistsCount := nullable.IntWithDefault(userVCS.PrivateGists, 0)           // Set 'private_gists' count to 0 if nil

	createdAt, _ := time.Parse(defaultDateShort, userVCS.CreatedAt.String())     // convert 'created_at' date to short format "2017-01-02 15:04:05 -0700 UTC"
	updatedAt, _ := time.Parse(defaultDateShort, userVCS.UpdatedAt.String())     // convert 'created_at' date to short format "2017-01-02 15:04:05 -0700 UTC"
	suspendedAt, _ := time.Parse(defaultDateShort, userVCS.SuspendedAt.String()) // convert 'created_at' date to short format "2017-01-02 15:04:05 -0700 UTC"

	// general info
	userMetaInfo := &User{
		Email:      userAccountEmail,
		AvatarURL:  userVCS.AvatarURL,
		HTMLURL:    userVCS.HTMLURL,
		GravatarID: userVCS.GravatarID,
		Name:       userVCS.Name,
		Company:    userVCS.Company,
		Blog:       userVCS.Blog,
		Location:   userVCS.Location,
		Hireable:   userVCS.Hireable,
		Bio:        userVCS.Bio,
		// Gender:          userVCS.Gender,
		// Birthday:        userBirthday,
		// Emails:        	extraInfo.UserInfo.Emails,
	}

	log.WithFields(logrus.Fields{
		"method.name":      "NewUserFromGithub(...)",
		"src.file":         "model/vcs-user.go",
		"prefix":           "vcs-github-user",
		"group":            "user-meta-info",
		"var.userMetaInfo": userMetaInfo,
	}).Info("user common profile info")

	// VCS related info
	userVcsGithubMetaInfo := &UserInfoVCS{
		UserID:            strconv.Itoa(*userVCS.ID),
		LoginEmail:        userAccountEmail,
		Login:             userVCS.Login,
		PublicRepos:       publicReposCount,
		PublicGists:       publicGistsCount,
		Followers:         followersCount,
		Following:         followingCount,
		CreatedAt:         createdAt,
		UpdatedAt:         updatedAt,
		SuspendedAt:       suspendedAt,
		TotalPrivateRepos: totalPrivateReposCount,
		OwnedPrivateRepos: ownedPrivateReposCount,
		PrivateGists:      privateGistsCount,
		Collaborators:     collaboratorsCount,
		AccountType:       userVCS.Type,
		DiskUsage:         userVCS.DiskUsage,
		SiteAdmin:         userVCS.SiteAdmin,
	}

	log.WithFields(logrus.Fields{
		"method.name":               "NewUserFromGithub(...)",
		"src.file":                  "model/vcs-user.go",
		"prefix":                    "vs-github-user",
		"group":                     "user-vcs-info",
		"var.userVcsGithubMetaInfo": userVcsGithubMetaInfo, // missing serviceID
	}).Info("user vcs related info")

	return userVcsGithubMetaInfo, nil

}

// CreateOrUpdateStar creates or updates a star and returns true if the star was created (vs updated)
func CreateOrUpdateUserVCS(db *gorm.DB, userVCS *UserInfoVCS, service *Service) (bool, error) {
	// Get existing by remote ID and service ID
	var existing UserInfoVCS
	if db.Where("user_id = ? AND service_id = ?", userVCS.UserID, service.ID).First(&existing).RecordNotFound() {
		userVCS.ServiceID = service.ID
		err := db.Create(userVCS).Error
		return err == nil, err
	}
	userVCS.ID = existing.ID
	userVCS.ServiceID = service.ID
	userVCS.CreatedAt = existing.CreatedAt
	return false, db.Save(userVCS).Error
}

func (user UserInfoVCS) DisplayLogin() string {
	return string(*user.Login)
}
