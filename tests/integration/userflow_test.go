package integration_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Laelapa/CompanyRegistry/internal/routes/handlers"
	"github.com/google/uuid"
)

func TestUserFlow(t *testing.T) {
	app := setupApp(t)

	username := "user" + strings.ReplaceAll(uuid.NewString(), "-", "")
	password := "TestPassword123!"
	companyName := "company" + uuid.NewString()[:8]

	var accessToken string

	t.Run("User Signup", func(t *testing.T) {
		reqPayload := handlers.UserSignupRequest{
			Username: username,
			Password: password,
		}
		reqBody, _ := json.Marshal(reqPayload)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/signup", bytes.NewReader(reqBody))
	})
}
