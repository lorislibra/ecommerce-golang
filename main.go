package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/donnjedarko/paninaro/config"
	"github.com/donnjedarko/paninaro/infrastructures/db"
	"github.com/donnjedarko/paninaro/internal/routes"
	"github.com/donnjedarko/paninaro/internal/web"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	cfg := config.Load()

	fiberCfg := fiber.Config{
		ReadTimeout:  cfg.AppTimeout,
		WriteTimeout: cfg.AppTimeout,
		IdleTimeout:  cfg.AppTimeout,
	}

	app := fiber.New(fiberCfg)
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowCredentials: true,
	}))

	mongo := db.NewMongoInstance()
	mongo.Connect()

	redis := db.NewRedisInstance()
	redis.Connect()

	webRouter := web.NewWebRouter(app, mongo, redis)
	mainRouter := routes.NewMainRouter(webRouter)
	mainRouter.GetRoutes()

	// app.Static("/*", "../frontend/build")

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt)

	go func() {
		<-stopChan
		mongo.Disconnect()

		log.Println("Stopping...")

		app.Shutdown()
	}()

	app.Listen(fmt.Sprintf("%v:%v", cfg.Host, cfg.Port))

}
