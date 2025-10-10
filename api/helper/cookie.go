package helper

import (
	"os"

	"github.com/gofiber/fiber/v2"
)

type CookieHelper interface {
	Set(value string, maxAge int)
	Get() string
	Delete()
}

type cookieHelper struct {
	c        *fiber.Ctx
	Name     string
	MaxAge   int
	HttpOnly bool
	Secure   bool
	SameSite string
	Path     string
	Domain   string
}

func NewCookieHelper(c *fiber.Ctx, name string) CookieHelper {
	return &cookieHelper{
		c:        c,
		Name:     name,
		MaxAge:   3600,
		HttpOnly: true,
		Secure:   os.Getenv("APP_ENV") == "production",
		SameSite: fiber.CookieSameSiteLaxMode,
		Path:     "/",
		Domain:   os.Getenv("COOKIE_DOMAIN"),
	}
}

func (ch *cookieHelper) Set(value string, maxAge int) {
	ch.c.Cookie(&fiber.Cookie{
		Name:     ch.Name,
		HTTPOnly: ch.HttpOnly,
		Secure:   ch.Secure,
		SameSite: ch.SameSite,
		Path:     ch.Path,
		Domain:   ch.Domain,
		MaxAge:   maxAge,
		Value:    value,
	})
}

func (ch *cookieHelper) Get() string {
	return ch.c.Cookies(ch.Name)
}

func (ch *cookieHelper) Delete() {
	ch.c.Cookie(&fiber.Cookie{
		Name:     ch.Name,
		HTTPOnly: ch.HttpOnly,
		Secure:   ch.Secure,
		SameSite: ch.SameSite,
		Path:     ch.Path,
		Domain:   ch.Domain,
		MaxAge:   -1,
		Value:    "",
	})
}
