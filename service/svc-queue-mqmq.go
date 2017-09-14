package service

import (
	"time" 														// go-core
	"github.com/disintegration/mqmq"							// queue-mqmq
	"github.com/sirupsen/logrus" 								// logs-logrus
)

func mainMQMQ() {
	mqmqListenAddr := ""
	c := mqmq.NewClient()
	err := c.Connect(mqmqListenAddr)
	if err != nil {
		log.Fatalf("failed to connect: %s", err)
		log.WithError(err).WithFields(
			logrus.Fields{
				"src.file": 			"service/svc-queue-mqmq.go",
				"prefix": 				"svc-queue",
				"method.name": 			"mainMQMQ(...)",
				"var.mqmqListenAddr": 	mqmqListenAddr,
				}).Fatalf("failed to connect...")
	}

	for received := 0; received < 10; {
		queueName := "queue1"
		// Wait at most 1 minute for the next message.
		msg, err := c.Get(queueName, 1*time.Minute)
		if err == mqmq.ErrTimeout {
			continue // No message so far. Keep on waiting.
		} else if err != nil {
			log.WithError(err).WithFields(
				logrus.Fields{
					"src.file": 			"service/svc-queue-mqmq.go",
					"prefix": 				"svc-queue",
					"method.name": 			"mainMQMQ(...)",
					"var.mqmqListenAddr": 	mqmqListenAddr,
					"var.queueName": 		queueName,
					}).Fatalf("failed to get message...")
		}
		received++
		log.Printf("received: %s", string(msg))
	}

	c.Disconnect()

}