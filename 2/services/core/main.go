package main

import (
	"fmt"
	"log"

	_ "github.com/alianjidaniir-design/SamplePRJ/models/task"
	_ "github.com/alianjidaniir-design/SamplePRJ/models/user"
	"github.com/alianjidaniir-design/SamplePRJ/services/core/route"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()
	routes := route.SetupRoutes(app)

	fmt.Println("Virasty-style API running on :8080")
	fmt.Printf("routes: %+v\n", routes)

	if err := app.Listen(":8080"); err != nil {
		log.Fatal(err)
	}
}
