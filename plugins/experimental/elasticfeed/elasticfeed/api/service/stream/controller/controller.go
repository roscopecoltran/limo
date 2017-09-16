package controller

import (
	"github.com/roscopecoltran/elasticfeed/service/stream/controller/room"
)

func InitSession() {
	room.InitSessionManager()
}
