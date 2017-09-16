package entities

import (
	"github.com/jmcvetta/neoism" 															// data-neo4j
	//"github.com/sirupsen/logrus" 															// logs-logrus
	//"github.com/davecgh/go-spew/spew" 													// debug-print
	//"github.com/k0kubun/pp" 																// debug-print
)

/*
Category model
*/
type Category struct {
	ID   		int64 		`json:"id" yaml:"id"`
	Name 		string 		`json:"name" yaml:"name"`
}

// ---------------------------------------------------------------------------

/*
GetCategories returns collection of news
*/
func GetCategories(db *neoism.Database) (*[]Category, error) {

	var categories []Category
	if err := db.Cypher(&neoism.CypherQuery{
		Statement:`MATCH (category:Category)
                RETURN DISTINCT ID(category) as id, category.name as name`,
		Result: &categories,
	}); err != nil {
		return nil, err
	}
	return &categories, nil
}
