package rabbitmqhandling

import (
	"context"
	"encoding/json"
	"log"
	"time"

	gh "github.com/go-multi-microservices/common/generalhandling"
	jh "github.com/go-multi-microservices/common/jsonhandling"
	reh "github.com/go-multi-microservices/common/redishandling"
	amqp "github.com/rabbitmq/amqp091-go"
)

func Producer(amqpUrl string, queueName string, data interface{}) {
	conn, err := amqp.Dial(amqpUrl)
	gh.HandleError(err, "Cannot connect to AMQP")
	defer func(conn *amqp.Connection) {
		err := conn.Close()
		if err != nil {
			gh.HandleError(err, "Error occurred while closing the connection of AMQP")
		}
	}(conn)

	ch, err1 := conn.Channel()
	gh.HandleError(err1, "Cannot create a amqpChannel")
	defer func(ch *amqp.Channel) {
		err := ch.Close()
		if err != nil {
			gh.HandleError(err, "Error occurred while closing the amqpChannel")
		}
	}(ch)

	q, err2 := ch.QueueDeclare(
		queueName,
		false,
		false,
		false,
		false,
		nil,
	)
	gh.HandleError(err2, "Could not declare "+queueName+" queue")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	body := jh.JsonMarshal(data)

	err3 := ch.PublishWithContext(
		ctx,
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        body,
		},
	)
	gh.HandleError(err3, "Error publishing tweet")

	log.Println("RabbitMQ: tweet published to the queue successfully")
}

func Consumer(amqpUrl string, queueName string) {
	conn, err := amqp.Dial(amqpUrl)
	gh.HandleError(err, "Cannot connect to AMQP")
	defer func(conn *amqp.Connection) {
		err := conn.Close()
		if err != nil {
			gh.HandleError(err, "Error occurred while closing the connection of AMQP")
		}
	}(conn)

	ch, err1 := conn.Channel()
	gh.HandleError(err1, "Cannot create a amqpChannel")
	defer func(ch *amqp.Channel) {
		err := ch.Close()
		if err != nil {
			gh.HandleError(err, "Error occurred while closing the amqpChannel")
		}
	}(ch)

	q, err2 := ch.QueueDeclare(
		queueName,
		false,
		false,
		false,
		false,
		nil,
	)
	gh.HandleError(err2, "Could not declare "+queueName+" queue")

	tweets, err3 := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	gh.HandleError(err3, "Could not register consumer")

	forever := make(chan string)

	go func() {
		for d := range tweets {
			var tweet jh.Tweet
			err := json.Unmarshal(d.Body, &tweet)
			gh.HandleError(err, "Error occurred while unmarshal retrieved tweets from rabbitmq:queue")

			check := reh.Get(reh.InitRedisConn(), gh.RedisDataKey)

			var tweets []jh.Tweet

			if len(check) == 0 || check == "null" || check == "" {
				tweets = append(tweets, tweet)
				reh.Set(reh.InitRedisConn(), gh.RedisDataKey, jh.JsonMarshal(tweets), 0)
				log.Printf("Data stored successfully for the first time to redis database")
			} else {
				err := json.Unmarshal([]byte(check), &tweets)
				gh.HandleError(err, "Error occurred while unmarshal tweets from redis:database")
				tweets = append(tweets, tweet)
				reh.Set(reh.InitRedisConn(), gh.RedisDataKey, jh.JsonMarshal(tweets), 0)
				log.Printf("Data stored successfully to redis database")
			}
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
