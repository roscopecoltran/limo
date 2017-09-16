package entities

import (
	"github.com/jmcvetta/neoism" 															// data-neo4j
	//"github.com/sirupsen/logrus" 															// logs-logrus
	//"github.com/davecgh/go-spew/spew" 													// debug-print
	//"github.com/k0kubun/pp" 																// debug-print
)

/*
NewsProvider model
*/
type NewsProvider struct {
	ID   		int64  				`json:"id" yaml:"id"`
	Name 		string 				`json:"name" yaml:"name"`
}

// ---------------------------------------------------------------------------

/*
GetNewsProviders returns collection of news
*/
func GetNewsProviders(db *neoism.Database) (*[]NewsProvider, error) {
	var newsproviders []NewsProvider
	if err := db.Cypher(&neoism.CypherQuery{
		Statement:`MATCH (newsprovider:NewsProvider)
                RETURN DISTINCT ID(newsprovider) as id, newsprovider.name as name`,
		Result: &newsproviders,
	}); err != nil {
		return nil, err
	}
	return &newsproviders, nil
}
