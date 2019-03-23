package utils

import (
	"fmt"

	"github.com/go-redis/redis"
)

func NewRedisClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	pong, err := client.Ping().Result()

	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Connected to Redis instance: " + pong)
	// Output: PONG <nil>
	return client
}

var RedisClient = NewRedisClient()
