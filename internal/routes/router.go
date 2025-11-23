package routes

import (
	"net/http"

	"github.com/Laelapa/CompanyRegistry/auth/tokenauthority"
	"github.com/Laelapa/CompanyRegistry/internal/middleware"
	"github.com/Laelapa/CompanyRegistry/internal/routes/handlers"
	"github.com/Laelapa/CompanyRegistry/internal/service"
	"github.com/Laelapa/CompanyRegistry/logging"

	"github.com/twmb/franz-go/pkg/kgo"
)

func Setup(
	staticDir string,
	logger *logging.Logger,
	service *service.Service,
	tokenAuthority *tokenauthority.TokenAuthority,
	kafkaClient *kgo.Client,
) *http.ServeMux {
	mux := http.NewServeMux()
	fileServer := http.StripPrefix("/static/", http.FileServer(http.Dir(staticDir)))

	h := handlers.New(logger, service, tokenAuthority, kafkaClient)

	// Wrapper for handlers that require authenticated access
	withAuth := func(handler func(http.ResponseWriter, *http.Request)) http.Handler {
		return middleware.AuthenticateWithJWT(tokenAuthority, logger)(http.HandlerFunc(handler))
	}

	mux.Handle("GET /static/", fileServer)

	mux.HandleFunc("POST /api/v1/login", h.HandleLogin)
	mux.HandleFunc("POST /api/v1/signup", h.HandleSignup)

	mux.HandleFunc("GET /api/v1/company/{name}", h.HandleGetCompanyByName)
	mux.Handle("POST /api/v1/company", withAuth(h.HandleCreateCompany))
	mux.Handle("PATCH /api/v1/company/{id}", withAuth(h.HandleUpdateCompany))
	mux.Handle("DELETE /api/v1/company/{id}", withAuth(h.HandleDeleteCompany))

	return mux
}
