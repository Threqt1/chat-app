package database

import (
	"chat-app/pkg/snowflake"
	"chat-app/repository"
	"context"
	"os"

	"github.com/redis/go-redis/v9"
)

const (
	MAIN_DATABASE = 0
	TEST_DATABASE = 1
)

// initializeRedis intiializes the Redis database connection and returns the connected client.
func intitializeRedis(dbType int) (*redis.Client, error) {
	var Addr string
	var Password string
	switch dbType {
	case 0:
		Addr = os.Getenv("REDIS_CONNECTION")
		Password = os.Getenv("REDIS_PASSWORD")
	case 1:
		Addr = os.Getenv("REDIS_LOCAL_CONNECTION")
		Password = os.Getenv("REDIS_LOCAL_PASSWORD")
	}

	//Connect to Redis
	connection := redis.NewClient(&redis.Options{
		Addr:     Addr,
		Password: Password,
		DB:       0,
	})

	//Check if Redis connection was succesful
	_, err := connection.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}

	return connection, nil
}

type redisDatabaseProvider struct {
	client            *redis.Client
	snowflakeProvider repository.SnowflakeService
}

// NewRedisDatabaseProvider creates a new Redis database service.
func NewRedisDatabaseProvider(dbType int) (redisDatabaseProvider, error) {
	rdb := redisDatabaseProvider{}

	snowflakeProvider, err := snowflake.NewSnowflakeProvider()
	if err != nil {
		return rdb, err
	}
	rdb.snowflakeProvider = snowflakeProvider

	connection, err := intitializeRedis(dbType)
	if err != nil {
		return rdb, err
	}
	rdb.client = connection

	return rdb, nil
}
