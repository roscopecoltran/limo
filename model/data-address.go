package model

import (
	"fmt"
	"github.com/jinzhu/gorm"
)

type Address struct {
	gorm.Model
	UserID      uint   `json:"user_id,omitempty" yaml:"user_id,omitempty"`
	Active      bool   `json:"active,omitempty" yaml:"active,omitempty"`
	Disabled    bool   `default:"false" json:"disabled,omitempty" yaml:"disabled,omitempty"`
	ContactName string `form:"contact-name" json:"contact_name,omitempty" yaml:"contact_name,omitempty"`
	Phone       string `form:"phone" json:"phone,omitempty" yaml:"phone,omitempty"`
	City        string `form:"city" json:"city,omitempty" yaml:"city,omitempty"`
	Address1    string `form:"address1" json:"address1,omitempty" yaml:"address1,omitempty"`
	Address2    string `form:"address2" json:"address2,omitempty" yaml:"address2,omitempty"`
	ExtraInfo   string `json:"extra_info,omitempty" yaml:"extra_info,omitempty"`
	DoorCode    string `json:"door_code,omitempty" yaml:"door_code,omitempty"`
	ZipCode     string `json:"zip_code,omitempty" yaml:"zip_code,omitempty"`
	CountryISO2 string `json:"country_iso2,omitempty" yaml:"country_iso2,omitempty"`
	CountryISO3 string `json:"country_iso3,omitempty" yaml:"country_iso3,omitempty"`
}

func (address Address) Stringify() string {
	return fmt.Sprintf("%v, %v, %v", address.Address2, address.Address1, address.City)
}
