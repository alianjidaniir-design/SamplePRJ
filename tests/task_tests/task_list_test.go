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

func TestListTask(t *testing.T) {
	app := fiber.New()
	route.SetupRoutes(app)

	createPayload := map[string]any{
		"body": map[string]any{
			"title":       "First task",
			"description": "seed",
		},
	}
	createBody, _ := json.Marshal(createPayload)
	createReq, _ := http.NewRequest(http.MethodPost, "/task/create", bytes.NewReader(createBody))
	createReq.Header.Set("Content-Type", "application/json")
	if _, err := app.Test(createReq); err != nil {
		t.Fatalf("create api test failed: %v", err)
	}

	listPayload := map[string]any{
		"body": map[string]any{
			"page":    1,
			"perPage": 10,
		},
	}
	listBody, _ := json.Marshal(listPayload)
	listReq, _ := http.NewRequest(http.MethodPost, "/task/list", bytes.NewReader(listBody))
	listReq.Header.Set("Content-Type", "application/json")
	listRes, err := app.Test(listReq)
	if err != nil {
		t.Fatalf("list api test failed: %v", err)
	}

	if listRes.StatusCode != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, listRes.StatusCode)
	}
}
