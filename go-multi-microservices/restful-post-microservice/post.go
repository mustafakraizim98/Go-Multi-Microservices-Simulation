package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	gh "github.com/go-multi-microservices/common/generalhandling"
	jh "github.com/go-multi-microservices/common/jsonhandling"
	rh "github.com/go-multi-microservices/common/rabbitmqhandling"
)

func ginPostHttpRequest(context *gin.Context) {
	var tweet jh.Tweet

	if err := context.BindJSON(&tweet); err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"response": "bad request"})
	} else {
		now := time.Now().Format("2006-01-02 15:04:05")
		tweet.CreatedAt, err = time.Parse("2006-01-02 15:04:05", now)
		gh.HandleError(err, "Error occurred while parsing the creation time")
		context.IndentedJSON(http.StatusCreated, tweet)

		rh.Producer(gh.AmqpUrl, gh.RabbitQueueName, tweet)
	}
}

func main() {
	router := gin.Default()

	router.POST("/tweet", ginPostHttpRequest)

	err := router.Run("localhost:9090")
	if err != nil {
		return
	}
}
