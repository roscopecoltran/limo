package entities

import (
	"github.com/jmcvetta/neoism" 															// data-neo4j
	//"github.com/sirupsen/logrus" 															// logs-logrus
	//"github.com/davecgh/go-spew/spew" 													// debug-print
	//"github.com/k0kubun/pp" 																// debug-print
)

/*
Company model
*/
type Company struct {
	ID   		int64  	`json:"id" yaml:"id"`
	Name 		string 	`json:"name" yaml:"name"`
}

// ---------------------------------------------------------------------------

/*
GetCompanies returns collection of news
*/
func GetCompanies(db *neoism.Database) (*[]Company, error) {

	var companies []Company
	if err := db.Cypher(&neoism.CypherQuery{
		Statement:`MATCH (company:Company)
                RETURN DISTINCT ID(company) as id, company.name as name`,
		Result: &companies,
	}); err != nil {
		return nil, err
	}
	return &companies, nil
}
