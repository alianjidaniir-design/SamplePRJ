package user_tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	_ "github.com/alianjidaniir-design/SamplePRJ/models/user"
	"github.com/alianjidaniir-design/SamplePRJ/services/core/route"
	"github.com/gofiber/fiber/v2"
)

func TestCreateUser(t *testing.T) {
	app := fiber.New()
	route.SetupRoutes(app)

	payload := map[string]any{
		"body": map[string]any{
			"username": "virasty",
			"email":    "virasty@example.com",
		},
	}

	bodyBytes, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("marshal payload failed: %v", err)
	}

	req, err := http.NewRequest(http.MethodPost, "/user/create", bytes.NewReader(bodyBytes))
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
