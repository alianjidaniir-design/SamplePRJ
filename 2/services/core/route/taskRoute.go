package route

import (
	. "github.com/alianjidaniir-design/SamplePRJ/controllers/task"
	"github.com/gofiber/fiber/v2"
)

var taskRoutes = map[string]string{
	"taskCreate": "/task/create",
	"taskList":   "/task/list",
	"taskUpdate": "/task/update",
	"taskDelete": "/task/delete",
}

func SetupTaskRoute(app *fiber.App) map[string]string {
	app.Post(taskRoutes["taskCreate"], Create)
	app.Get(taskRoutes["taskList"], List)
	app.Post(taskRoutes["taskUpdate"], Update)
	app.Post(taskRoutes["taskDelete"], Delete)
	return taskRoutes
}
