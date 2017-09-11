package model

import (
	"errors"
	"fmt"
	"strings"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

// LanguageDetected represents a languageDetected in the database
type LanguageDetected struct {
	gorm.Model
	Name        			string 				`gorm:"-" yaml:"name,omitempty" json:"name,omitempty"`
	Type        			string 				`yaml:"type,omitempty" json:"type,omitempty"`
	Group       			string 				`yaml:"group,omitempty" json:"group,omitempty"`
	AceMode     			string 				`yaml:"ace_mode,omitempty" json:"ace_mode,omitempty"`
	IsPopular   			bool   				`yaml:"is_popular,omitempty" json:"is_popular,omitempty"`
	IsUnpopular 			bool   				`yaml:"is_unpopular,omitempty" json:"is_unpopular,omitempty"`
	LanguageDetectedCount 	int    				`gorm:"-"`
	Stars     				[]Star 				`gorm:"many2many:star_languagesDetected;"`
}

// Detection represents a language detection result
type LanguageDetection struct {
	gorm.Model
	Path                   string    			`yaml:"path,omitempty" json:"path,omitempty"`
	Type                   string    			`yaml:"type,omitempty" json:"type,omitempty"`
	ExtName                string    			`yaml:"extname,omitempty" json:"extname,omitempty"`
	MimeType               string    			`yaml:"mime_type,omitempty" json:"mime_type,omitempty"`
	ContentType            string    			`yaml:"content_type,omitempty" json:"content_type,omitempty"`
	Disposition            string    			`yaml:"disposition,omitempty" json:"disposition,omitempty"`
	IsDocumentation        bool      			`yaml:"is_documentation,omitempty" json:"is_documentation,omitempty"`
	IsLarge                bool      			`yaml:"is_large,omitempty" json:"is_large,omitempty"`
	IsGenerated            bool      			`yaml:"is_generated,omitempty" json:"is_generated,omitempty"`
	IsText                 bool      			`yaml:"is_text,omitempty" json:"is_text,omitempty"`
	IsImage                bool      			`yaml:"is_image,omitempty" json:"is_image,omitempty"`
	IsBinary               bool      			`yaml:"is_binary,omitempty" json:"is_binary,omitempty"`
	IsVendored             bool      			`yaml:"is_vendored,omitempty" json:"is_vendored,omitempty"`
	IsHighRatioOfLongLines bool      			`yaml:"is_high_ratio_of_long_lines,omitempty" json:"is_high_ratio_of_long_lines,omitempty"`
	IsViewable             bool      			`yaml:"is_viewable,omitempty" json:"is_viewable,omitempty"`
	IsSafeToColorize       bool      			`yaml:"is_safe_to_colorize,omitempty" json:"is_safe_to_colorize,omitempty"`
	Language               *LanguageDetected 	`yaml:"language,omitempty" json:"language,omitempty"`
}

// Result is the result details of a detection
type LanguageDetectionResult struct {
	Success    				bool       			`yaml:"success" json:"success"`
	Message    				string     			`yaml:"message,omitempty" json:"message,omitempty"`
	Result     				*LanguageDetection 	`yaml:"result" json:"result"`
	IsBinary   				bool       			`yaml:"binary" json:"binary"`
	IsLarge    				bool       			`yaml:"large" json:"large"`
	IsExcluded 				bool       			`yaml:"excluded" json:"excluded"`
}

// FindLanguageDetecteds finds all languagesDetected
func FindLanguageDetecteds(db *gorm.DB) ([]LanguageDetected, error) {
	var languagesDetected []LanguageDetected
	db.Order("name").Find(&languagesDetected)
	return languagesDetected, db.Error
}

// FindLanguageDetectedsWithStarCount finds all languagesDetected and gets their count of stars
func FindLanguageDetectedsWithStarCount(db *gorm.DB) ([]LanguageDetected, error) {
	var languagesDetected []LanguageDetected

	rows, err := db.Raw(`
		SELECT T.NAME, COUNT(ST.LANGUAGE_ID) AS STARCOUNT
		FROM LANGUAGES_DETECTED T
		LEFT JOIN STAR_LANGUAGES_DETECTED ST ON T.ID = ST.LANGUAGE_ID
		WHERE T.DELETED_AT IS NULL
		GROUP BY T.ID
		ORDER BY T.NAME`).Rows()

	if err != nil {
		return languagesDetected, err
	}

	defer func() {
		err := rows.Close()
		if err != nil {
			log.WithError(err).WithFields(logrus.Fields{"service": "FindLanguageDetectedsWithStarCount", "action": "languagesDetected"}).Fatalf("error: %#s", err)
		}
	}()

	for rows.Next() {
		var languageDetected LanguageDetected
		if err = rows.Scan(&languageDetected.Name, &languageDetected.StarCount); err != nil {
			return languagesDetected, err
		}
		languagesDetected = append(languagesDetected, languageDetected)
	}
	return languagesDetected, db.Error
}

// FindLanguageDetectedByName finds a languageDetected by name
func FindLanguageDetectedByName(db *gorm.DB, name string) (*LanguageDetected, error) {
	var languageDetected LanguageDetected
	if db.Where("lower(name) = ?", strings.ToLower(name)).First(&languageDetected).RecordNotFound() {
		return nil, db.Error
	}
	return &languageDetected, db.Error
}

// FindOrCreateLanguageDetectedByName finds a languageDetected by name, creating if it doesn't exist
func FindOrCreateLanguageDetectedByName(db *gorm.DB, name string) (*LanguageDetected, bool, error) {
	var languageDetected LanguageDetected
	if db.Where("lower(name) = ?", strings.ToLower(name)).First(&languageDetected).RecordNotFound() {
		languageDetected.Name = name
		err := db.Create(&languageDetected).Error
		return &languageDetected, true, err
	}
	return &languageDetected, false, nil
}

// LoadStars loads the stars for a language
func (languageDetected *LanguageDetected) LoadStars(db *gorm.DB, match string) error {
	// Make sure languageDetected exists in database, or we will panic
	var existing LanguageDetected
	if db.Where("id = ?", languageDetected.ID).First(&existing).RecordNotFound() {
		return fmt.Errorf("LanguageDetected '%d' not found", languageDetected.ID)
	}

	if match != "" {
		var stars []Star
		db.Raw(`
			SELECT *
			FROM LANGUAGES_DETECTED S
			INNER JOIN STAR_LANGUAGES_DETECTED ST ON S.ID = ST.STAR_ID
			WHERE S.DELETED_AT IS NULL
			AND ST.LANGUAGE_ID = ?
			AND LOWER(S.FULL_NAME) LIKE ?
			ORDER BY S.FULL_NAME`,
			languageDetected.ID,
			fmt.Sprintf("%%%s%%", strings.ToLower(match))).Scan(&stars)
		languageDetected.Stars = stars
		return db.Error
	}
	return db.Model(languageDetected).Association("Stars").Find(&languageDetected.Stars).Error
}

// Rename renames a languageDetected -- new name must not already exist
func (languageDetected *LanguageDetected) Rename(db *gorm.DB, name string) error {
	// Can't rename to the same name
	if name == languageDetected.Name {
		return errors.New("You can't rename to the same name")
	}

	// If they're just changing case, allow. Otherwise, block the change
	if strings.ToLower(name) != strings.ToLower(languageDetected.Name) {
		existing, err := FindLanguageDetectedByName(db, name)
		if err != nil {
			log.WithError(err).WithFields(logrus.Fields{"service": "FindLanguageDetectedsWithStarCount", "action": "FindLanguageDetectedByName"}).Fatalf("error: %#s", err)
			return err
		}
		if existing != nil {
			return fmt.Errorf("LanguageDetected '%s' already exists", existing.Name)
		}
	}
	languageDetected.Name = name
	return db.Save(languageDetected).Error
}

// Delete deletes a languageDetected and disassociates it from any stars
func (languageDetected *LanguageDetected) Delete(db *gorm.DB) error {
	if err := db.Model(languageDetected).Association("Stars").Clear().Error; err != nil {
		return err
	}
	return db.Delete(languageDetected).Error
}
