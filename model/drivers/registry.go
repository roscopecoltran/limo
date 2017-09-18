package rdb

const DBS_RDB_DEFAULT_DRIVER 		= 		"gorm"
const DBS_RDB_DEFAULT_ADAPTER 		= 		"sqlite3"

var (
	ValidClients 					= 		[]string{"sqlite3", "postgres", "mysql"}
)
