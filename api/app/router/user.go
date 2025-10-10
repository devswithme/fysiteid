package router

import (
	"fybe/controller"
	"fybe/helper"
	"fybe/middleware"
	"fybe/model/dto"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/go-redis/redis_rate/v10"
	jwtware "github.com/gofiber/contrib/jwt"

	"github.com/gofiber/fiber/v2"
)

func UserRouter(c controller.UserController, app fiber.Router, validator *validator.Validate, limiter *redis_rate.Limiter) {
	user := app.Group("/user")
	user.Use(middleware.LimiterMiddleware(limiter, redis_rate.PerSecond(10)))
	user.Use(jwtware.New(jwtware.Config{
		SigningKey:   jwtware.SigningKey{Key: []byte(os.Getenv("A_SECRET"))},
		TokenLookup:  "cookie:access_token",
		ErrorHandler: helper.JWTErrorHandler,
	}), middleware.ParseUID())

	user.Get("/me", c.Me)
	user.Use(middleware.ValidatorMiddleware[dto.UserUpdate](validator)).Patch("/me", c.Update)

	user_public := app.Group("/public/user")
	user_public.Get("/:username", c.GetByUsername)
	user_public.Get("/id/:ticketID", c.GetByTicketID)
}
