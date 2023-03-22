package main

import (
	"fmt"
	"log"

	"github.com/ShikharY10/gbAUTH/cmd/admin"
	config "github.com/ShikharY10/gbAUTH/cmd/configs"
	"github.com/ShikharY10/gbAUTH/cmd/controllers/c_v1"
	"github.com/ShikharY10/gbAUTH/cmd/handlers"
	"github.com/ShikharY10/gbAUTH/cmd/middlewares"
	"github.com/ShikharY10/gbAUTH/cmd/routes/route_v1"
	"github.com/gin-gonic/gin"
)

func main() {
	env := config.LoadENV()

	if env.GIN_MODE == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	logger, err := admin.InitializeLogger(env, "AUTH")
	if err != nil {
		fmt.Println("1")
		log.Fatal(err)
	}

	mongoDB, err := config.ConnectMongoDB(env)
	if err != nil {
		fmt.Println("2")
		log.Fatal(err)
	}

	redis, err := config.ConnectRedis(env)
	if err != nil {
		fmt.Println("3")
		log.Fatal(err)
	}

	cloudinary := config.InitCloudinary(env)

	cache := handlers.InitializeCacheHandler(redis)
	dataBase := handlers.InitializeDataBase(mongoDB, logger)

	handler := &handlers.Handler{
		Logger:     logger,
		Cloudinary: cloudinary,
		Cache:      cache,
		DataBase:   dataBase,
	}

	middleware := middlewares.InitializeMiddleware(env, dataBase, cache)

	authController := &c_v1.AuthController{
		Handler:    handler,
		Middleware: middleware,
	}

	engine := gin.New()
	bashRoute := engine.Group("/api/v1")

	route_v1.AuthRoutes(bashRoute, authController)

	engine.Run(":" + env.AUTH_WEBSERVER_PORT)

}
