package user_tests

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	_ "github.com/alianjidaniir-design/SamplePRJ/models/user"
	"github.com/alianjidaniir-design/SamplePRJ/services/core/route"
	"github.com/gofiber/fiber/v2"
)

func TestInfoUser(t *testing.T) {
	app := fiber.New()
	route.SetupRoutes(app)

	createPayload := map[string]any{
		"body": map[string]any{
			"username": "virasty",
			"email":    "virasty@example.com",
		},
	}
	createBodyBytes, _ := json.Marshal(createPayload)
	createReq, _ := http.NewRequest(http.MethodPost, "/user/create", bytes.NewReader(createBodyBytes))
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

	userMap, ok := dataMap["user"].(map[string]any)
	if !ok {
		t.Fatal("invalid create response: missing user")
	}

	userIDFloat, ok := userMap["id"].(float64)
	if !ok {
		t.Fatal("invalid create response: missing user.id")
	}

	infoPayload := map[string]any{
		"body": map[string]any{
			"userID": int64(userIDFloat),
		},
	}
	infoBodyBytes, _ := json.Marshal(infoPayload)
	infoReq, _ := http.NewRequest(http.MethodPost, "/user/info", bytes.NewReader(infoBodyBytes))
	infoReq.Header.Set("Content-Type", "application/json")
	infoRes, err := app.Test(infoReq)
	if err != nil {
		t.Fatalf("info api test failed: %v", err)
	}

	if infoRes.StatusCode != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, infoRes.StatusCode)
	}
}
