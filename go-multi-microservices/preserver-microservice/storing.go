package main

import (
	gh "github.com/go-multi-microservices/common/generalhandling"
	rh "github.com/go-multi-microservices/common/rabbitmqhandling"
)

func main() {
	rh.Consumer(gh.AmqpUrl, gh.RabbitQueueName)
}
