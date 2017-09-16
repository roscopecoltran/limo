package model

import (
	"github.com/roscopecoltran/elasticfeed/service/stream"
)

type ServiceManager interface {

	GetStreamService() *stream.StreamService

	Init()
}
