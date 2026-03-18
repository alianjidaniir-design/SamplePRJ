package student_tests

import (
	"Fiber/API/2/services/core/route"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestCreateTask(t *testing.T) {
	app := fiber.New()
	route.SetupRoute(app)

}
