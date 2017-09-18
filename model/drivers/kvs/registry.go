package kvs

const DBS_KVS_DEFAULT_DRIVER 	= 		"gorm"
const DBS_KVS_DEFAULT_ADAPTER 	= 		"sqlite3"

var (
	DefaultClients 				= map[string]bool{"boltdb": true, "etcd": true}
)