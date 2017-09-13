package model

/*
import (
	// "log"
	"github.com/sirupsen/logrus"
	"github.com/garyburd/redigo/redis"
	// "github.com/roscopecoltran/sniperkit-limo/model"
	"github.com/roscopecoltran/sniperkit-limo/config"	
)

//CacheRedis connects to a redis node and encapsulates caching.
type CacheRedis interface {
	Get(key string)
	Set(key string, val string)
}

//CacheDrivers holds the redis information and implements CacheRedis interface.
type CacheDrivers struct {
	redisConn redis.Conn
}

//NewLinkAggCache generates a new cache.
func NewLinkCache(config *config.Config) *CacheDrivers {
	var c CacheDrivers
	// log.Println("Connecting to Redis instance on local port,", config.GetString("Redis.port"))
	//var conn, err = redis.Dial("tcp", ":"+config.GetString("Redis.port"))
	var conn, err = redis.Dial("tcp", ":"+config.Cache.Redis.Port
	if err != nil {
		log.WithError(err).WithFields(logrus.Fields{"action": "NewLinkCache", "step": "Dial", "service": "cache"}).Warnln("Error dialing redis instance", err)
	}
	c.redisConn = conn
	return &c
}

//Get fetches entry from Redis instance if it exists, else returns "".
func (cache *CacheDrivers) Get(key string) string {
	val, err := cache.redisConn.Do("GET", key)
	if err != nil {
		log.WithError(err).WithFields(logrus.Fields{"action": "NewLinkCache", "step": "Dial", "service": "cache"}).Warnln("Unable get from redis", err)
		return ""
	}
	if str, ok := val.(string); ok {
		return str
	}
	return ""
}

//Set saves to cache and overwrites any previous value.
func (cache *CacheDrivers) Set(key string, val string) {
	cache.redisConn.Send("SET", key, val)
	cache.redisConn.Flush()
}

*/

