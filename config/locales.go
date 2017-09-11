package config

import (
	"strings"
	"os"
	// logs
	"github.com/sirupsen/logrus"
)

// locale gets the locale from the LANG env var if not set.
func locale() string {
	if clientLocale != "" {
		log.WithFields(logrus.Fields{"action": "locale", "step": "clientLocale", "service": "locales"}).Warn("Could not get the locale from the LANG env")
		return clientLocale
	}
	var lang string
	if os.Getenv("LANG") != "" {
		lang = os.Getenv("LANG")
		if lang == "" {
			return ""
		}
	} else {
		lang = clientLocaleDefault
	}
	locale := strings.Split(lang, ".")[0]
	return locale
}

// 
/*
// Region returns the users region code.
// Eg. "US", "GB", etc
func Region() string {
	l := locale()

	tag, err := language.Parse(l)
	if err != nil {
		return defaultRegion
	}

	region, _ := tag.Region()

	return region.String()
}

// Language returns the users language code.
// Eg. "en", "es", etc
func Language() string {
	l := locale()

	tag, err := language.Parse(l)
	if err != nil {
		return defaultLanguage
	}

	base, _ := tag.Base()

	return base.String()
}

// SetClientLocale sets the locale of the client
// connecting to the web server.
func SetClientLocale(locale string) {
	clientLocale = locale
}
*/

