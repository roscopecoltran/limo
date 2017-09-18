package entities

type GraphRelation struct {
	Id			int
	Type       	string      `json:"type" yaml:"type"`
	Data       	interface{} `json:"data" yaml:"data"`
	Extensions 	interface{} `json:"extensions" yaml:"extensions"`

	StartNode	*GraphNode
	EndNode		*GraphNode
}