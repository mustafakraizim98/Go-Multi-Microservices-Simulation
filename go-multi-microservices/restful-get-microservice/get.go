package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	gh "github.com/go-multi-microservices/common/generalhandling"
	jh "github.com/go-multi-microservices/common/jsonhandling"
	reh "github.com/go-multi-microservices/common/redishandling"
	"net/http"
	"sort"
	"strings"
)

type timeSlice []jh.Tweet

func (p timeSlice) Len() int {
	return len(p)
}

func (p timeSlice) Less(i, j int) bool {
	// After: for descending order
	return p[i].CreatedAt.After(p[j].CreatedAt)
}

func (p timeSlice) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func sorting(data []jh.Tweet) timeSlice {
	dateSortedReviews := make(timeSlice, 0, len(data))
	for _, d := range data {
		dateSortedReviews = append(dateSortedReviews, d)
	}
	sort.Sort(dateSortedReviews)
	return dateSortedReviews
}

func getDataFromRedisDatabase() []jh.Tweet {
	var tweets []jh.Tweet

	data := reh.Get(reh.InitRedisConn(), gh.RedisDataKey)
	if len(data) == 0 {
		return tweets
	} else {
		err := json.Unmarshal([]byte(data), &tweets)
		gh.HandleError(err, "GET REQUEST | Error occurred while unmarshal retrieved tweets from redis:database")
	}

	return tweets
}

func ginGetHttpRequest(context *gin.Context) {
	if len(getDataFromRedisDatabase()) == 0 {
		context.IndentedJSON(http.StatusInsufficientStorage, gin.H{"response": "redis database is empty"})
	} else {
		context.IndentedJSON(http.StatusOK, sorting(getDataFromRedisDatabase()))
	}
}

func getDataFromRedisDatabaseByParameter(context *gin.Context) {
	if len(getDataFromRedisDatabase()) == 0 {
		context.IndentedJSON(http.StatusInsufficientStorage, gin.H{"response": "redis database is empty"})
	} else {
		var founded []jh.Tweet

		counter := 0
		for _, d := range getDataFromRedisDatabase() {
			if strings.EqualFold(d.Creator, context.Param("creator")) {
				founded = append(founded, d)
				counter++
			}
		}

		if counter == 0 {
			context.IndentedJSON(http.StatusNotFound, gin.H{"response": "record not found"})
		} else {
			context.IndentedJSON(http.StatusAccepted, sorting(founded))
		}
	}
}

func main() {
	router := gin.Default()

	router.GET("/tweet/list", ginGetHttpRequest)
	router.GET("/tweet/list/:creator", getDataFromRedisDatabaseByParameter)

	err := router.Run("localhost:9091")
	if err != nil {
		return
	}
}
