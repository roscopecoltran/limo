package news

import (
	"github.com/jmcvetta/neoism" 															// data-neo4j
	//"github.com/sirupsen/logrus" 															// logs-logrus
	//"github.com/davecgh/go-spew/spew" 													// debug-print
	//"github.com/k0kubun/pp" 																// debug-print
)

/*
Person model
*/
type Person struct {
	ID   int64  `json:"id" yaml:"id"`
	Name string `json:"name" yaml:"name"`
}

// ---------------------------------------------------------------------------

/*
GetPeople returns collection of news
*/
func GetPeople(db *neoism.Database) (*[]Person, error) {

	var people []Person
	if err := db.Cypher(&neoism.CypherQuery{
		Statement:`MATCH (person:Person)
                RETURN DISTINCT ID(person) as id, person.name as name`,
		Result: &people,
	}); err != nil {
		return nil, err
	}
	return &people, nil
}
