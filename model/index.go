package model

import (
	"fmt"
	"github.com/google/go-github/github"
	jsoniter "github.com/json-iterator/go"
	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/analysis/analyzers/keyword_analyzer"
	"github.com/blevesearch/bleve/analysis/analyzers/simple_analyzer"
	"github.com/blevesearch/bleve/analysis/language/en"
)

const mapping = `
{
	"mappings": {
		"_default_": {
			"dynamic_templates": [
				{
					"strings": {
						"match_mapping_type": "string",
						"mapping": {
							"index": "not_analyzed"
						}
					}
				}
			],
			"properties": {
				"body": {"type": "string", "index": "analyzed"}
			}
		}
	}
}`

// InitIndex initializes the search index at the specified path
func InitIndex(filepath string) (bleve.Index, error) {
	index, err := bleve.Open(filepath)

	// Doesn't yet exist (or error opening) so create a new one
	if err != nil {
		index, err = bleve.New(filepath, buildIndexMapping())
		if err != nil {
			return nil, err
		}
	}
	return index, nil
}

func buildIndexMapping() *bleve.IndexMapping {
	simpleTextFieldMapping := bleve.NewTextFieldMapping()
	simpleTextFieldMapping.Analyzer = simple_analyzer.Name

	englishTextFieldMapping := bleve.NewTextFieldMapping()
	englishTextFieldMapping.Analyzer = en.AnalyzerName

	keywordFieldMapping := bleve.NewTextFieldMapping()
	keywordFieldMapping.Analyzer = keyword_analyzer.Name

	starMapping := bleve.NewDocumentMapping()
	starMapping.AddFieldMappingsAt("Name", simpleTextFieldMapping)
	starMapping.AddFieldMappingsAt("FullName", simpleTextFieldMapping)
	starMapping.AddFieldMappingsAt("Description", englishTextFieldMapping)
	starMapping.AddFieldMappingsAt("Language", keywordFieldMapping)
	starMapping.AddFieldMappingsAt("Tags.Name", keywordFieldMapping)
	starMapping.AddFieldMappingsAt("Topics.Name", keywordFieldMapping)
	starMapping.AddFieldMappingsAt("Languages.Name", keywordFieldMapping)
	// starMapping.AddFieldMappingsAt("Readmes.Content", keywordFieldMapping)

	indexMapping := bleve.NewIndexMapping()
	indexMapping.AddDocumentMapping("Star", starMapping)

	// indexMapping.AddDocumentMapping("Repo", starMapping)


	return indexMapping
}

// https://github.com/dastergon/strgz/blob/master/lib/bleve.go
func ShowResults(results *bleve.SearchResult, index bleve.Index) {
	if len(results.Hits) < 1 {
		fmt.Println(results)
	}
	for _, val := range results.Hits {
		id := val.ID
		doc, err := index.Document(id)
		if err != nil {
			fmt.Println(err)
		}
		for _, field := range doc.Fields {
			repo := github.Repository{}
			jsoniter.Unmarshal(field.Value(), &repo)
			fmt.Printf("%s - %s (%s)\n\t%s\n", *repo.Name, *repo.Description, *repo.Language, *repo.HTMLURL)
		}
	}
}
