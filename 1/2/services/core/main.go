package main

import (
	"Fiber/API/2/services/core/route"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()
	routes := route.SetupRoute(app)

	fmt.Printf("routes: %+v\n -- ", routes)
	fmt.Print("Starting server on port 8080 --- ")
	fmt.Print("university")

	if err := app.Listen(":8080"); err != nil {
		panic(err)
	}

}
