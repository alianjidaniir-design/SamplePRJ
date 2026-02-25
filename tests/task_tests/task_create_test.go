package task_tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	_ "github.com/alianjidaniir-design/SamplePRJ/models/task"
	"github.com/alianjidaniir-design/SamplePRJ/services/core/route"
	"github.com/gofiber/fiber/v2"
)

func TestCreateTask(t *testing.T) {
	app := fiber.New()
	route.SetupRoutes(app)

	payload := map[string]any{
		"body": map[string]any{
			"title":       "Learn Virasty style",
			"description": "Build clean layered project",
		},
	}

	bodyBytes, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("marshal payload failed: %v", err)
	}

	req, err := http.NewRequest(http.MethodPost, "/task/create", bytes.NewReader(bodyBytes))
	if err != nil {
		t.Fatalf("build request failed: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("api test failed: %v", err)
	}

	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, res.StatusCode)
	}
}
