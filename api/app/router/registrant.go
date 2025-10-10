package router

import (
	"fybe/controller"
	"fybe/helper"
	"fybe/middleware"
	"os"

	"github.com/go-redis/redis_rate/v10"
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
)

func RegistrantRouter(c controller.RegistrantController, app fiber.Router, limiter *redis_rate.Limiter) {
	registrant := app.Group("/registrant")
	registrant.Use(middleware.LimiterMiddleware(limiter, redis_rate.PerSecond(10)))

	registrant.Use(jwtware.New(jwtware.Config{
		SigningKey:   jwtware.SigningKey{Key: []byte(os.Getenv("A_SECRET"))},
		TokenLookup:  "cookie:access_token",
		ErrorHandler: helper.JWTErrorHandler,
	}), middleware.ParseUID())

	registrant.Get("/", c.GetByUserID)
	registrant.Get("/:ticketID", c.GetByTicketID)
	registrant.Post("/:ticketID", c.Create)
	registrant.Get("/gen/:ticketID", c.Generate)
	registrant.Get("/verify/:id/:ticketID", c.Verify)
}
