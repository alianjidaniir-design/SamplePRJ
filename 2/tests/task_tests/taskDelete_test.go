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

func TestSoftDeleteTask(t *testing.T) {
	app := fiber.New()
	route.SetupRoutes(app)

	createPayload := map[string]any{
		"body": map[string]any{
			"title":       "to-delete",
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

	taskID := int64(taskIDFloat)

	listBeforeReq, _ := http.NewRequest(http.MethodGet, "/task/list?page=1&perPage=100", nil)
	listBeforeRes, err := app.Test(listBeforeReq)
	if err != nil {
		t.Fatalf("list (before delete) api test failed: %v", err)
	}
	if listBeforeRes.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(listBeforeRes.Body)
		t.Fatalf("expected status %d, got %d body=%s", http.StatusOK, listBeforeRes.StatusCode, string(body))
	}

	listBeforeResBody, err := io.ReadAll(listBeforeRes.Body)
	if err != nil {
		t.Fatalf("read list (before delete) response failed: %v", err)
	}

	var listBeforeResult map[string]any
	if err := json.Unmarshal(listBeforeResBody, &listBeforeResult); err != nil {
		t.Fatalf("unmarshal list (before delete) response failed: %v", err)
	}
	listBeforeData, ok := listBeforeResult["data"].(map[string]any)
	if !ok {
		t.Fatal("invalid list (before delete) response: missing data")
	}
	totalBeforeFloat, ok := listBeforeData["total"].(float64)
	if !ok {
		t.Fatal("invalid list (before delete) response: missing total")
	}
	tasksBeforeAny, ok := listBeforeData["tasks"].([]any)
	if !ok {
		t.Fatal("invalid list (before delete) response: missing tasks array")
	}

	foundBefore := false
	for _, item := range tasksBeforeAny {
		taskMap, ok := item.(map[string]any)
		if !ok {
			continue
		}
		idFloat, ok := taskMap["id"].(float64)
		if !ok {
			continue
		}
		if int64(idFloat) == taskID {
			foundBefore = true
			break
		}
	}
	if !foundBefore {
		t.Fatalf("expected created task id=%d to exist before delete", taskID)
	}

	deletePayload := map[string]any{
		"body": map[string]any{
			"taskID": taskID,
		},
	}
	deleteBody, _ := json.Marshal(deletePayload)
	deleteReq, _ := http.NewRequest(http.MethodPost, "/task/delete", bytes.NewReader(deleteBody))
	deleteReq.Header.Set("Content-Type", "application/json")
	deleteRes, err := app.Test(deleteReq)
	if err != nil {
		t.Fatalf("delete api test failed: %v", err)
	}
	if deleteRes.StatusCode != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, deleteRes.StatusCode)
	}

	listReq, _ := http.NewRequest(http.MethodGet, "/task/list?page=1&perPage=100", nil)
	listRes, err := app.Test(listReq)
	if err != nil {
		t.Fatalf("list api test failed: %v", err)
	}
	if listRes.StatusCode != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, listRes.StatusCode)
	}

	listResBody, err := io.ReadAll(listRes.Body)
	if err != nil {
		t.Fatalf("read list response failed: %v", err)
	}

	var listResult map[string]any
	if err := json.Unmarshal(listResBody, &listResult); err != nil {
		t.Fatalf("unmarshal list response failed: %v", err)
	}
	listData, ok := listResult["data"].(map[string]any)
	if !ok {
		t.Fatal("invalid list response: missing data")
	}
	tasksAny, ok := listData["tasks"].([]any)
	if !ok {
		t.Fatal("invalid list response: missing tasks array")
	}

	foundAfter := false
	for _, item := range tasksAny {
		taskMap, ok := item.(map[string]any)
		if !ok {
			continue
		}
		idFloat, ok := taskMap["id"].(float64)
		if !ok {
			continue
		}
		if int64(idFloat) == taskID {
			foundAfter = true
			break
		}
	}
	if foundAfter {
		t.Fatalf("expected deleted task id=%d to not appear in list after soft delete", taskID)
	}

	totalAfterFloat, ok := listData["total"].(float64)
	if !ok {
		t.Fatal("invalid list response: missing total")
	}
	if int(totalAfterFloat) != int(totalBeforeFloat)-1 {
		t.Fatalf("expected total to decrement by 1 after soft delete (before=%d after=%d)", int(totalBeforeFloat), int(totalAfterFloat))
	}
}
