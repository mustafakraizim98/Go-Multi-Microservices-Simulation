package redishandling

import (
	"context"
	"fmt"
	"time"

	gh "github.com/go-multi-microservices/common/generalhandling"
	"github.com/go-redis/redis"
)

type Database struct {
	Client *redis.Client
}

func NewClient(address string) (*Database, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: "",
		DB:       0,
	})

	if err := client.Ping().Err(); err != nil {
		return nil, fmt.Errorf("%s: %s", context.TODO(), "no matching record found in redis database")
	}

	return &Database{
		Client: client,
	}, nil
}

func Set(databse *Database, key string, value interface{}, expiration time.Duration) {
	err := databse.Client.Set(key, value, expiration).Err()
	gh.HandleError(err, "failed while setting value to redis database")
}

func Get(databse *Database, key string) string {
	value, err := databse.Client.Get(key).Result()
	if err == redis.Nil {
		return ""
	} else {
		gh.HandleError(err, "failed while retrieving data from redis database")
	}
	return value
}

func InitRedisConn() *Database {
	database, err := NewClient(gh.RedisUrl)
	gh.HandleError(err, "Error occurred while connecting to redis database")
	return database
}
