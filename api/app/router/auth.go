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

func AuthRouter(c controller.AuthController, app fiber.Router, limiter *redis_rate.Limiter) {
	login := app.Group("/login")
	login.Use(middleware.LimiterMiddleware(limiter, redis_rate.PerSecond(100)))
	login.Get("/google", c.Login)
	login.Get("/google/callback", c.Callback)

	auth := app.Group("/auth")
	auth.Use(jwtware.New(jwtware.Config{
		SigningKey:   jwtware.SigningKey{Key: []byte(os.Getenv("R_SECRET"))},
		TokenLookup:  "cookie:refresh_token",
		ErrorHandler: helper.JWTErrorHandler,
	}), middleware.ParseUID())

	auth.Post("/refresh", c.Refresh)
	auth.Post("/logout", c.Logout)
}
