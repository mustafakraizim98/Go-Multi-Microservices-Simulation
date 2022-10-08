package generalhandling

import "log"

const (
	AmqpUrl         string = "amqp://username:password@localhost:9001/"
	RabbitQueueName string = "TweetsQueue"
	RedisUrl        string = "localhost:6379"
	RedisDataKey    string = "tweets"
)

func HandleError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
