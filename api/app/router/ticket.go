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

func TicketRouter(c controller.TicketController, app fiber.Router, validator *validator.Validate, limiter *redis_rate.Limiter) {
	ticket := app.Group("/ticket")
	ticket.Use(middleware.LimiterMiddleware(limiter, redis_rate.PerSecond(10)))

	ticket.Use(jwtware.New(jwtware.Config{
		SigningKey:   jwtware.SigningKey{Key: []byte(os.Getenv("A_SECRET"))},
		TokenLookup:  "cookie:access_token",
		ErrorHandler: helper.JWTErrorHandler,
	}), middleware.ParseUID())

	ticket.Get("/", c.Get)
	ticket.Get("/:id", c.GetByID)
	ticket.Delete("/:id", c.Delete)
	ticket.Use(middleware.ValidatorMiddleware[dto.TicketCreate](validator)).Use(middleware.LimiterMiddleware(limiter, redis_rate.PerSecond(5))).Post("/", c.Create)
	ticket.Use(middleware.ValidatorMiddleware[dto.TicketCreate](validator)).Patch("/:id", c.Update)

	ticket_public := app.Group("/public/ticket")
	ticket_public.Get("/:username", c.GetPublicByUsername)
	ticket_public.Get("/id/:id", c.GetPublicByID)
}
