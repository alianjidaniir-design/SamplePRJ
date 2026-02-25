package route

import (
	. "github.com/alianjidaniir-design/SamplePRJ/controllers/task"
	"github.com/gofiber/fiber/v2"
)

var taskRoutes = map[string]string{
	"taskCreate": "/task/create",
	"taskList":   "/task/list",
}

func SetupTaskRoute(app *fiber.App) map[string]string {
	app.Post(taskRoutes["taskCreate"], Create)
	app.Post(taskRoutes["taskList"], List)
	return taskRoutes
}
