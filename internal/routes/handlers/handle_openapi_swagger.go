package handlers

import (
	"net/http"

	"go.uber.org/zap"
)

func (h *Handler) HandleSwaggerUI(w http.ResponseWriter, r *http.Request) {
	html := `<!DOCTYPE html>
	<html>
		<head>
    		<title>Company Registry API</title>
    		<link rel="stylesheet" type="text/css" href="https://unpkg.com/swagger-ui-dist@5.22.0/swagger-ui.css" />
		</head>
		<body>
    		<div id="swagger-ui"></div>
    		<script src="https://unpkg.com/swagger-ui-dist@5.22.0/swagger-ui-bundle.js"></script>
    		<script>
        		SwaggerUIBundle({
            		url: '/openapi.json',
            		dom_id: '#swagger-ui',
            		presets: [SwaggerUIBundle.presets.apis, SwaggerUIBundle.presets.standalone]
        		});
    		</script>
		</body>
	</html>`
	w.Header().Set("Content-Type", "text/html")
	if _, err := w.Write([]byte(html)); err != nil {
		h.logger.Error("Couldn't write response", zap.Error(err))
	}
	h.logger.Info("Served Swagger UI", h.logger.ReqFields(r)...)
}

func (h *Handler) HandleGetOpenAPI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	http.ServeFile(w, r, "docs/openapi.json")
	h.logger.Info("Served OpenAPI JSON", h.logger.ReqFields(r)...)
}
