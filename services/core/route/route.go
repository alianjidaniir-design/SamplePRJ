package route

import "github.com/gofiber/fiber/v2"

func SetupRoutes(app *fiber.App) map[string]string {
	return mergeMaps(
		SetupTaskRoute(app),
		SetupUserRoute(app),
	)
}

func mergeMaps(maps ...map[string]string) map[string]string {
	mergedMap := map[string]string{}
	for _, currentMap := range maps {
		for key, value := range currentMap {
			mergedMap[key] = value
		}
	}

	return mergedMap
}
