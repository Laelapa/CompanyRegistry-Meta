package integration_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Laelapa/CompanyRegistry/internal/app"
	"github.com/Laelapa/CompanyRegistry/internal/routes/handlers"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
		w := sendPostRequest(app, "/api/v1/signup", reqPayload, "")

		require.Equal(t, http.StatusCreated, w.Code)

		var authResp handlers.AuthResponse
		err := json.Unmarshal(w.Body.Bytes(), &authResp)
		require.NoError(t, err)
		require.NotEmpty(t, authResp.AccessToken)

		accessToken = authResp.AccessToken
	})

	t.Run("Signup with existing username fails", func(t *testing.T) {
		reqPayload := handlers.UserSignupRequest{
			Username: username,
			Password: password,
		}
		w := sendPostRequest(app, "/api/v1/signup", reqPayload, "")

		require.Equal(t, http.StatusConflict, w.Code)
	})

	t.Run("Create Company", func(t *testing.T) {
		var ec int32 = 50
		reg := true
		reqPayload := handlers.CreateCompanyRequest{
			Name:          companyName,
			EmployeeCount: &ec,
			Registered:    &reg,
			CompanyType:   "Corporation",
		}
		w := sendPostRequest(app, "/api/v1/company", reqPayload, accessToken)

		require.Equal(t, http.StatusCreated, w.Code)

		var resp handlers.CompanyResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)
		require.Equal(t, companyName, resp.Name)
		assert.Equal(t, ec, resp.EmployeeCount)
		assert.Equal(t, reg, resp.Registered)
		assert.Equal(t, "Corporation", resp.CompanyType)
	})

	t.Run("Get Company", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/api/v1/company/"+companyName, nil)
		w := httptest.NewRecorder()
		app.ServeHTTP(w, r)

		require.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Create duplicate Company fails", func(t *testing.T) {
		var ec int32 = 50
		reg := true
		reqPayload := handlers.CreateCompanyRequest{
			Name:          companyName,
			EmployeeCount: &ec,
			Registered:    &reg,
			CompanyType:   "Corporation",
		}
		w := sendPostRequest(app, "/api/v1/company", reqPayload, accessToken)

		require.Equal(t, http.StatusConflict, w.Code)
	})
}

func sendPostRequest(app *app.App, url string, body any, accessToken string) *httptest.ResponseRecorder {
	reqBody, _ := json.Marshal(body)
	r := httptest.NewRequest(http.MethodPost, url, bytes.NewReader(reqBody))
	r.Header.Set("Content-Type", "application/json")
	if accessToken != "" {
		r.Header.Set("Authorization", "Bearer "+accessToken)
	}
	w := httptest.NewRecorder()
	app.ServeHTTP(w, r)
	return w
}
