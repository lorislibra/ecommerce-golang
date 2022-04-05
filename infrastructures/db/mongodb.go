package db

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/donnjedarko/paninaro/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type MongoInstance interface {
	Db() *mongo.Database
	Client() *mongo.Client
	Connect()
	Disconnect()
}

type mongoInstance struct {
	once   sync.Once
	client *mongo.Client
	db     *mongo.Database
}

func NewMongoInstance() MongoInstance {
	cfg := config.Get()

	conUrl := fmt.Sprintf(
		"mongodb+srv://%s:%s@%s/%s/?retryWrites=true&w=majority",
		cfg.MongoUsername, cfg.MongoPassword, cfg.MongoHost, cfg.MongoDbName,
	)
	log.Println(conUrl)
	opt := options.Client().ApplyURI(conUrl)

	client, err := mongo.NewClient(opt)
	if err != nil {
		log.Fatalln(err.Error())
	}

	database := client.Database(cfg.MongoDbName)

	return &mongoInstance{
		client: client,
		db:     database,
	}
}

func (m *mongoInstance) Db() *mongo.Database {
	return m.db
}

func (m *mongoInstance) Client() *mongo.Client {
	return m.client
}

func (m *mongoInstance) Connect() {
	m.once.Do(func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute*1)
		defer cancel()

		if err := m.client.Connect(ctx); err != nil {
			log.Fatalln(err.Error())
		}

		if err := m.client.Ping(ctx, readpref.Primary()); err != nil {
			log.Fatalln(err.Error())
		}
		log.Println("Mongo: PING OK")
	})
}

func (m *mongoInstance) Disconnect() {
	m.client.Disconnect(context.Background())
}
