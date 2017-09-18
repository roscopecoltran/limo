package model

import (
	"time"
	"github.com/jinzhu/gorm"
	// "github.com/sirupsen/logrus"
)

type ExternalURL struct {
	gorm.Model
	URL      			string 				`gorm:"column:external_url" json:"external_url,omitempty" yaml:"external_url,omitempty"`
	Patterns     		[]PatternEntry 		`gorm:"column:pattern_matched" json:"pattern_matched,omitempty" yaml:"pattern_matched,omitempty"`
	AuthRequired      	bool 				`default:"false" gorm:"column:auth_required" json:"auth_required,omitempty" yaml:"auth_required,omitempty"`
	AuthCredentialsID   int64 				`gorm:"column:auth_credentials_id" json:"auth_credentials_id,omitempty" yaml:"auth_credentials_id,omitempty"`
	LastUrlStatus     	uint 				`gorm:"column:last_status" json:"last_status,omitempty" yaml:"last_status,omitempty"`
	LastUrlSuccess      *time.Time 			`gorm:"column:last_url_success" json:"last_url_success,omitempty" yaml:"last_url_success,omitempty"`
	LastUrlVisit      	*time.Time 			`gorm:"column:last_visit" json:"last_visit,omitempty" yaml:"last_visit,omitempty"`
	ProxyRequired      	bool 				`default:"false" gorm:"column:proxy_required" json:"proxy_required,omitempty" yaml:"proxy_required,omitempty"`
	ProxyType      		uint 				`gorm:"column:proxy_type" json:"proxy_type,omitempty" yaml:"proxy_type,omitempty"`
	LinkPriority    	uint 				`default:"1" gorm:"column:link_priority" json:"link_priority,omitempty" yaml:"link_priority,omitempty"`
	LinkRank      		float64 			`default:"0" gorm:"column:link_rank" json:"link_rank,omitempty" yaml:"link_rank,omitempty"`
	PageRank      		float64 			`default:"0" gorm:"column:page_rank" json:"page_rank,omitempty" yaml:"page_rank,omitempty"`
}
