package service

/*
import (
	"github.com/nsqio/go-nsq"
)

// ref. https://github.com/toorop/tmail/blob/master/core/scope.go

var (
	NsqQueueProducer                 *nsq.Producer
)

// initMailQueueProducer init producer for queue
func initMailQueueProducer() (err error) {
	nsqCfg := nsq.NewConfig()
	nsqCfg.UserAgent = "tmail.queue"
	NsqQueueProducer, err = nsq.NewProducer("127.0.0.1:4150", nsqCfg)
	if Cfg.GetDebugEnabled() {
		NsqQueueProducer.SetLogger(NewNSQLogger(), nsq.LogLevelDebug)
	} else {
		NsqQueueProducer.SetLogger(NewNSQLogger(), nsq.LogLevelError)
	}
	return err
}

*/