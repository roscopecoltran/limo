package graph

const DBS_GRAPH_DEFAULT_DRIVER 		= 		"cayley"
const DBS_GRAPH_DEFAULT_ADAPTER 	= 		"cayley"

var (
	DefaultClients 					= 		map[string]bool{"neo4j": true,  "cayley": true}
	ValidClients 					= 		[]string{"neo4j", "cayley", "dgraph"}
)
