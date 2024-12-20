package routes

import (
	"bytes"
	"dukebward/search/db"
	"dukebward/search/views"

	"github.com/gofiber/fiber/v2"
)

type searchInput struct {
	Term string `json:"term"`
}

// search handler
func HandleSearch(c *fiber.Ctx) error {
	input := searchInput{}
	if err := c.BodyParser(&input); err != nil {
		c.Status(500)
		c.Append("content-type", "application/json")
		return c.JSON(fiber.Map{
			"success": false,
			"message": "Invalid input",
			"data":    nil,
		})
	}
	if input.Term == "" {
		c.Status(500)
		c.Append("content-type", "application/json")
		return c.JSON(fiber.Map{
			"success": false,
			"message": "Invalid input",
			"data":    nil,
		})
	}
	idx := &db.SearchIndex{}
	data, err := idx.FullTextSearch(input.Term)
	if err != nil {
		c.Status(500)
		c.Append("content-type", "application/json")
		return c.JSON(fiber.Map{
			"success": false,
			"message": "Invalid input",
			"data":    nil,
		})
	}

	// Convert your templ results component to a string
	var buf bytes.Buffer
	views.SearchResults(data).Render(c.Context(), &buf)
	return c.SendString(buf.String())

}
