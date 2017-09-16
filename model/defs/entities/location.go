package entities

import (
	"github.com/jmcvetta/neoism" 															// data-neo4j
	//"github.com/sirupsen/logrus" 															// logs-logrus
	//"github.com/davecgh/go-spew/spew" 													// debug-print
	//"github.com/k0kubun/pp" 																// debug-print
)

/*
Location model
*/
type Location struct {
	ID   		int64  		`json:"id" yaml:"id"`
	Name 		string 		`json:"name" yaml:"name"`
}

// ---------------------------------------------------------------------------

/*
GetLocations returns collection of news
*/
func GetLocations(db *neoism.Database) (*[]Location, error) {
	var locations []Location
	if err := db.Cypher(&neoism.CypherQuery{
		Statement:`MATCH (location:Location)
                RETURN DISTINCT ID(location) as id, location.name as name`,
		Result: &locations,
	}); err != nil {
		return nil, err
	}
	return &locations, nil
}
