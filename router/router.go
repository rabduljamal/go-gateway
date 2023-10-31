package router

import (
	"github.com/gofiber/fiber/v2"
	metabase_handler "github.com/rabduljamal/gateway-snip/hendler/metabase"
)

// SetupRoutes func
func SetupRoutes(app *fiber.App) {
	// grouping
	api := app.Group("/api")
	v1 := api.Group("/v1")

	// user := v1.Group("/user")
	// user.Get("/", user_handler.GetAllUsers)
	// user.Get("/:id", user_handler.GetSingleUser)
	// user.Post("/", user_handler.CreateUser)
	// user.Put("/:id", user_handler.UpdateUser)
	// user.Delete("/:id", user_handler.DeleteUserByID)

	metabase := v1.Group("/metabase")
	metabase.Post("/", metabase_handler.GetMetabases)
}
