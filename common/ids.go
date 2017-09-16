package common

import (
	"github.com/roscopecoltran/elasticfeed/resource"
)

func AdminChannelID(admin *resource.Admin) string {
	return GetMd5(admin.Id + admin.Org.Id)
}
