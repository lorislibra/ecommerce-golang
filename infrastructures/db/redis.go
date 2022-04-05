package db

import (
	"fmt"
	"log"
	"sync"

	"github.com/donnjedarko/paninaro/config"
	"github.com/go-redis/redis/v8"
)

type RedisInstance interface {
	Client() *redis.Client
	Connect()
}

type redisInstance struct {
	once   sync.Once
	client *redis.Client
}

func NewRedisInstance() RedisInstance {
	cfg := config.Get()

	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%v:%v", cfg.RedisHost, cfg.RedisPort),
		Username: cfg.RedisUsername,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDbId,
	})

	return &redisInstance{
		client: rdb,
	}
}

func (ri *redisInstance) Client() *redis.Client {
	return ri.client
}

func (ri *redisInstance) Connect() {
	ri.once.Do(func() {
		client := ri.Client()
		ctx := client.Context()

		ping, err := client.Ping(ctx).Result()
		if err != nil {
			log.Fatalln("Redis Connection: ", err)
		}

		log.Println("Redis: PING -> ", ping)
	})
}
