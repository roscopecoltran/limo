package entities

import (
	"github.com/jmcvetta/neoism" 															// data-neo4j
	//"github.com/sirupsen/logrus" 															// logs-logrus
	//"github.com/davecgh/go-spew/spew" 													// debug-print
	//"github.com/k0kubun/pp" 																// debug-print
)

/*
Tag model
*/
type Tag struct {
	ID   int64  `json:"id" yaml:"id"`
	Name string `json:"name" yaml:"name"`
}

// ---------------------------------------------------------------------------

/*
GetTags returns collection of news
*/
func GetTags(db *neoism.Database) (*[]Tag, error) {
	var tags []Tag
	if err := db.Cypher(&neoism.CypherQuery{
		Statement:`MATCH (tag:Tag)
                RETURN DISTINCT ID(tag) as id, tag.name as name`,
		Result: &tags,
	}); err != nil {
		return nil, err
	}
	return &tags, nil
}
