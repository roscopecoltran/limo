package feedify

import (
	// Golang packages
	"fmt"
	"strconv"

	// Beego framework packages
	"github.com/astaxie/beego"

	// feedify packages
	"github.com/roscopecoltran/feedify/config"
	_ "github.com/roscopecoltran/feedify/stream/adapter/message"
	_ "github.com/roscopecoltran/feedify/graph/adapter"
)

func GetConfigKey(key string) string {
	return config.GetConfigKey(key)
}

func Banner() {
	fmt.Printf("Starting app '%s' on port '%s'\n", config.GetConfigKey("appname"), config.GetConfigKey("feedify::port"))
}

func SetStaticPath(url string, path string) *beego.App {
	return beego.SetStaticPath(url, path)
}

func Error(v ...interface{}) {
	beego.Error(v...)
}

func Run() {
	Banner()

	beego.HttpPort, _ = strconv.Atoi(config.GetConfigKey("feedify::port"))
	beego.Run()
}
