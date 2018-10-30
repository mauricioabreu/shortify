package main

import (
	"log"

	"github.com/go-redis/redis"
)

// RedisClient : create a new client to handle Redis operations
func RedisClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	pong, err := client.Ping().Result()
	log.Println(pong, err)
	return client
}
