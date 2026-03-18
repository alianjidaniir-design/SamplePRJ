package route

import (
	"Fiber/API/2/controllers/Student"

	"github.com/gofiber/fiber/v2"
)

var studentRoute = map[string]string{
	"studentCreate": "/student/create",
}

func SetupStudentRoute(app *fiber.App) map[string]string {
	app.Post(studentRoute["studentCreate"], student.Create)
	return studentRoute
}
