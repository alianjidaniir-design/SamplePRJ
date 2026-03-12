package route

import (
	. "github.com/alianjidaniir-design/SamplePRJ/controllers/user"
	"github.com/gofiber/fiber/v2"
)

var userRoutes = map[string]string{
	"userCreate": "/user/create",
	"userInfo":   "/user/info",
}

func SetupUserRoute(app *fiber.App) map[string]string {
	app.Post(userRoutes["userCreate"], Create)
	app.Post(userRoutes["userInfo"], Info)
	return userRoutes
}
