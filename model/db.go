package model

import (
	"github.com/jinzhu/gorm"
	// Use the sqlite dialect
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
    // "github.com/qor/qor"
    // "github.com/qor/admin"
    // "github.com/boltdb/bolt"
)

// https://github.com/qor/qor-example/blob/master/db/db.go

// InitDB initializes the database at the specified path
func InitDB(filepath string, verbose bool) (*gorm.DB, error) {

	db, err := gorm.Open("sqlite3", filepath)
	if err != nil {
		return nil, err
	}

	db.LogMode(verbose)
	// , &Deps{}, &Patterns{}, &Snippets{}
	db.AutoMigrate(&Service{}, &Star{}, &Tag{}, &Topic{}, &LanguageDetected{}, &Tree{}, &Readme{}, &Academic{}, &Pkg{}, &Software{}, &Repo{}, &Keyword{}, &Pattern{})

	// Initalize
	// Admin := admin.New(&qor.Config{DB: db})

	return db, nil

}

