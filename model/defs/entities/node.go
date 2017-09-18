package entities

type GraphNode struct {
	Id				int
	Data			map[string]interface{} `json:"data" yaml:"data"`
	Extensions		map[string]interface{} `json:"extensions" yaml:"extensions"`
	Labels			[]string
}

func (n *GraphNode) GetRelation(int) (*GraphRelation) {
	return &GraphRelation{}
}