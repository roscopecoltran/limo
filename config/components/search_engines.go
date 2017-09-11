package components

import (
	"github.com/jinzhu/configor"
)

type SolrConfig struct {
	Active 				bool 			`default:"false"`	
}

type SphinxSearchConfig struct {
	Active 				bool 			`default:"false"`	
}

// SearchEnginesConfig represents information about a full-text search engine parameters.
type SearchEnginesConfig struct {
	Active 				bool 			`default:"false"`	
	Engines struct {
		ElasticSearch 	ElasticSearchConfig
		SphinxSearch 	SphinxSearchConfig
		Solr 			SolrConfig
	}
}
