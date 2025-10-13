package main

import (
	"fybe/app"
	"fybe/helper"
	"os"
	"time"

	"github.com/gofiber/contrib/swagger"
	"go.uber.org/zap"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"

	_ "github.com/joho/godotenv/autoload"
)

// @title		FyAPI
// @version	1.0
// @host		localhost:3000
// @BasePath	/api/v1
func main() {
	logger := app.InitLogger()
	defer logger.Sync()

	redis := helper.NewRedisHelper(os.Getenv("REDIS_ADDR"), os.Getenv("REDIS_PWD"), 0, "app:", 10*time.Second)

	db := app.InitDB(logger)
	router := app.InitApp(db, logger, redis)

	fApp := fiber.New()

	fApp.Use(recover.New())
	fApp.Static("/upload", "upload")

	if os.Getenv("APP_ENV") == "development" {
		fApp.Use(swagger.New(swagger.Config{
			FilePath: "./docs/swagger.json",
		}))
	}
	fApp.Use(cors.New(cors.Config{
		AllowOrigins:     os.Getenv("FE_DOMAIN"),
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowMethods:     "GET, POST, PATCH, DELETE",
		AllowCredentials: true,
	}))

	app.InitRouter(fApp, router, redis)
	if err := fApp.Listen(":" + os.Getenv("BE_PORT")); err != nil {
		logger.Error("failed to start server: ", zap.Error(err))
	}
}
