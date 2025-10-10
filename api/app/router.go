package app

import (
	"fybe/app/router"
	"fybe/helper"
	"fybe/middleware"
	"fybe/model"

	"github.com/go-playground/validator/v10"
	"github.com/go-redis/redis_rate/v10"
	"github.com/gofiber/fiber/v2"
)

func InitRouter(f *fiber.App, app *model.App, redis helper.RedisHelper) {
	limiter := redis_rate.NewLimiter(redis.Client())

	api := f.Group("/api/v1")
	api.Use(middleware.RequestIDMiddleware())

	validator := validator.New()
	validator.RegisterValidation("slug", helper.ValidateSlug)

	router.AuthRouter(app.Controller.Auth, api, limiter)
	router.UserRouter(app.Controller.User, api, validator, limiter)
	router.TicketRouter(app.Controller.Ticket, api, validator, limiter)
	router.RegistrantRouter(app.Controller.Registrant, api, limiter)
}
