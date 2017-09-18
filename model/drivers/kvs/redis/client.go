package redis

/*
refs:
- https://github.com/silentred/toolkit/blob/master/service/redis.go
- https://github.com/feedlabs/feedify/blob/master/redis/client.go
- 
*/

import (
	"errors"
	"github.com/fzzy/radix/extra/pubsub"
	"github.com/fzzy/radix/redis"
	// redis "gopkg.in/redis.v5"
	"github.com/roscopecoltran/sniperkit-limo/config" 							// app-config
	"github.com/sirupsen/logrus"												// logs-logrus
	prefixed "github.com/x-cray/logrus-prefixed-formatter" 						// logs-logrus
)

const PKG_REDIS_LABEL_CLUSTER 		= 		"dbs"
const PKG_REDIS_LABEL_GROUP 		= 		"dbs-kvs"
const PKG_REDIS_LABEL_PREFIX 		= 		"dbs-kvs-redis"
const PKG_REDIS_LABEL_DRIVER 		= 		"radix/redis" 											// radix/redis, redis.v5

var dbs 							*model.RootDrivers
var log 							= logrus.New()

func init() {
	log.Out 						= 		os.Stdout 							// logs-logrus
	formatter 						:= 		new(prefixed.TextFormatter) 		// logs-logrus
	log.Formatter 					= 		formatter 							// logs-logrus
	log.Level 						= 		logrus.DebugLevel 					// logs-logrus
}

// Gorm Resources
type RedisRes 	struct {
	Ok 				bool  			`default:"false" json:"-" yaml:"-"`
	Cli 			*gorm.DB 		`json:"-" yaml:"-"`
}

type RedisClient struct {
	host     string
	port     string
	protocol string
}

func (r RedisClient) Cmd(command string, args ...interface{}) error {
	c, err := redis.Dial(r.protocol, r.host + ":" + r.port)
	if err != nil {
		return errors.New("Redis dial error")
	}
	c.Cmd(command, args)
	return nil
}

func (r RedisClient) _subscribe(channel []string, callback func(bool, string, string)) error {
	c, err := redis.Dial(r.protocol, r.host + ":" + r.port)
	if err != nil {
		return errors.New("Redis dial error")
	}

	psc := pubsub.NewSubClient(c)
	psr := psc.Subscribe(channel)
	for {
		psr = psc.Receive()
		callback(psr.Timeout(), psr.Message, psr.Channel)
	}

	return nil
}

func (r RedisClient) Subscribe(channel []string, callback func(bool, string, string)) {
	go r._subscribe(channel, callback)
}

func New() (*RedisClient) {
	host := config.GetConfigKey("redis::host")
	port := config.GetConfigKey("redis::port")
	protocol := config.GetConfigKey("redis::protocol")
	return &RedisClient{host, port, protocol}
}

/*
// NewRedisClient get a redis client
func NewRedisClient(cfg config.RedisConfig) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisInstance.Address(),
		DB:       cfg.RedisInstance.Db,
		Password: cfg.RedisInstance.Pwd, // no password set
	})

	if cfg.Ping {
		if err := client.Ping().Err(); err != nil {
			log.Fatal(err)
		}
	}

	return client
}

func initRedis(app Application) error {
	if app.GetConfig().Redis.InitRedis {
		redis := NewRedisClient(app.GetConfig().Redis)
		app.Set("redis", redis, nil)
	}
	return nil
}

*/

