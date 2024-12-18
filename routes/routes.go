package routes

import (
	"dukebward/search/db"

	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
)

// the backtick is a struct tag and you can do this
// with forms and json
// need to match the name of the name attribute in the form
type settingsForm struct {
	Amount   int  `form:amount`
	SearchOn bool `form:searchOn`
	AddNew   bool `form:addNew`
}

// export function so capital letter start
func SetRoutes(app *fiber.App) {
	app.Get("/", AuthMiddleware, DashboardHandler)
	app.Post("/", AuthMiddleware, LoginPostHandler)
	app.Get("/login", LoginHandler)
	app.Post("/login", LoginPostHandler)
	app.Get("/create", func(c *fiber.Ctx) error {
		u := &db.User{}
		u.CreateAdmin()
		return c.SendString("admin created")
	})
}

// non exported function
func render(c *fiber.Ctx, component templ.Component, options ...func(*templ.ComponentHandler)) error {
	componentHandler := templ.Handler(component)
	for _, o := range options {
		o(componentHandler)
	}
	return adaptor.HTTPHandler(componentHandler)(c)
}
