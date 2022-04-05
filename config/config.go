package config

import (
	"log"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/joho/godotenv"
)

var doOnce sync.Once
var config *Config

type Config struct {
	Host string
	Port string

	AppTimeout time.Duration

	MongoHost     string
	MongoPort     string
	MongoUsername string
	MongoPassword string
	MongoDbName   string

	RedisHost     string
	RedisPort     string
	RedisUsername string
	RedisPassword string
	RedisDbId     int

	PrivateKey []byte
	PublicKey  []byte

	JwtAccessTokenExpire      time.Duration
	JwtRefreshTokenExpire     time.Duration
	JwtRefreshTokenCookieName string
}

func load() *Config {
	c := Config{
		Host: os.Getenv("HOST"),
		Port: os.Getenv("PORT"),

		MongoHost:     os.Getenv("MONGO_HOST"),
		MongoPort:     os.Getenv("MONGO_PORT"),
		MongoUsername: os.Getenv("MONGO_USERNAME"),
		MongoPassword: os.Getenv("MONGO_PASSWORD"),
		MongoDbName:   os.Getenv("MONGO_DB_NAME"),

		RedisHost:     os.Getenv("REDIS_HOST"),
		RedisPort:     os.Getenv("REDIS_PORT"),
		RedisUsername: os.Getenv("REDIS_USERNAME"),
		RedisPassword: os.Getenv("REDIS_PASSWORD"),

		PrivateKey: []byte(os.Getenv("PRIVATE_KEY")),
		PublicKey:  []byte(os.Getenv("PUBLIC_KEY")),

		JwtRefreshTokenCookieName: os.Getenv("JWT_REFRESH_TOKEN_COOKIE_NAME"),
	}

	c.JwtAccessTokenExpire, _ = time.ParseDuration(os.Getenv("JWT_ACCESS_TOKEN_EXPIRE"))
	c.JwtRefreshTokenExpire, _ = time.ParseDuration(os.Getenv("JWT_REFRESH_TOKEN_EXPIRE"))
	c.AppTimeout, _ = time.ParseDuration(os.Getenv("APP_TIMEOUT"))
	c.RedisDbId, _ = strconv.Atoi(os.Getenv("REDIS_DB_ID"))

	return &c
}

func Load() *Config {
	doOnce.Do(func() {
		err := godotenv.Load()
		if err != nil {
			log.Println(err)
		}

		config = load()
	})
	return config
}

func Get() *Config {
	return config
}
