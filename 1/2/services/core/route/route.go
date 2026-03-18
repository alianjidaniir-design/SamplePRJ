package route

import "github.com/gofiber/fiber/v2"

func SetupRoute(app *fiber.App) map[string]string {
	return mergeMaps(
		SetupStudentRoute(app),
	)
}

func mergeMaps(maps ...map[string]string) map[string]string {
	result := map[string]string{}
	for _, m := range maps {
		for k, v := range m {
			result[k] = v
		}

	}
	return result

}
