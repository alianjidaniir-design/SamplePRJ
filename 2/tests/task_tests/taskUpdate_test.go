package task_tests

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	_ "github.com/alianjidaniir-design/SamplePRJ/models/task"
	"github.com/alianjidaniir-design/SamplePRJ/services/core/route"
	"github.com/gofiber/fiber/v2"
)

func TestUpdateTask(t *testing.T) {
	app := fiber.New()
	route.SetupRoutes(app)

	createPayload := map[string]any{
		"body": map[string]any{
			"title":       "before",
			"description": "seed",
		},
	}
	createBody, _ := json.Marshal(createPayload)
	createReq, _ := http.NewRequest(http.MethodPost, "/task/create", bytes.NewReader(createBody))
	createReq.Header.Set("Content-Type", "application/json")
	createRes, err := app.Test(createReq)
	if err != nil {
		t.Fatalf("create api test failed: %v", err)
	}
	if createRes.StatusCode != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, createRes.StatusCode)
	}

	createResBody, err := io.ReadAll(createRes.Body)
	if err != nil {
		t.Fatalf("read create response failed: %v", err)
	}

	var createResult map[string]any
	if err := json.Unmarshal(createResBody, &createResult); err != nil {
		t.Fatalf("unmarshal create response failed: %v", err)
	}
	dataMap, ok := createResult["data"].(map[string]any)
	if !ok {
		t.Fatal("invalid create response: missing data")
	}
	taskMap, ok := dataMap["task"].(map[string]any)
	if !ok {
		t.Fatal("invalid create response: missing task")
	}
	taskIDFloat, ok := taskMap["id"].(float64)
	if !ok {
		t.Fatal("invalid create response: missing task.id")
	}

	updatePayload := map[string]any{
		"body": map[string]any{
			"taskID":      int64(taskIDFloat),
			"title":       "after",
			"description": "updated",
		},
	}
	updateBody, _ := json.Marshal(updatePayload)
	updateReq, _ := http.NewRequest(http.MethodPost, "/task/update", bytes.NewReader(updateBody))
	updateReq.Header.Set("Content-Type", "application/json")
	updateRes, err := app.Test(updateReq)
	if err != nil {
		t.Fatalf("update api test failed: %v", err)
	}
	if updateRes.StatusCode != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, updateRes.StatusCode)
	}
}
