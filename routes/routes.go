package routes

import (
	"time"

	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/gofiber/fiber/v2/middleware/cache"
)

// export function so capital letter start
func SetRoutes(app *fiber.App) {
	app.Get("/", AuthMiddleware, DashboardHandler)
	app.Post("/", AuthMiddleware, DashboardPostHandler)

	app.Get("/login", LoginHandler)
	app.Post("/login", LoginPostHandler)
	app.Post("/logout", LogoutHandler)

	app.Post("/search", HandleSearch)
	app.Use("/search", cache.New(cache.Config{
		Next: func(c *fiber.Ctx) bool {
			return c.Query("noCache") == "true"
		},
		Expiration:   30 * time.Minute,
		CacheControl: true,
	}))

	// manual user creation
	// app.Get("/create", func(c *fiber.Ctx) error {
	// 	u := &db.User{}
	// 	u.CreateAdmin()
	// 	return c.SendString("admin created")
	// })
}

// non exported function
func render(c *fiber.Ctx, component templ.Component, options ...func(*templ.ComponentHandler)) error {
	componentHandler := templ.Handler(component)
	for _, o := range options {
		o(componentHandler)
	}
	return adaptor.HTTPHandler(componentHandler)(c)
}
