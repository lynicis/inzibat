package recorder

import (
	"github.com/gofiber/fiber/v2"
)

// RegisterAdminRoutes registers the recorder admin API endpoints on the Fiber app.
// All routes are under /_inzibat/recorder/.
func RegisterAdminRoutes(app *fiber.App, store *Store) {
	group := app.Group("/_inzibat/recorder")

	group.Get("/entries", listEntriesHandler(store))
	group.Get("/session", getSessionHandler(store))
	group.Post("/clear", clearEntriesHandler(store))
}

func listEntriesHandler(store *Store) fiber.Handler {
	return func(c *fiber.Ctx) error {
		entries := store.List()
		return c.JSON(entries)
	}
}

func getSessionHandler(store *Store) fiber.Handler {
	return func(c *fiber.Ctx) error {
		session := store.Session()
		return c.JSON(session)
	}
}

func clearEntriesHandler(store *Store) fiber.Handler {
	return func(c *fiber.Ctx) error {
		store.Clear()
		return c.JSON(fiber.Map{
			"message": "all recorded entries cleared",
		})
	}
}
