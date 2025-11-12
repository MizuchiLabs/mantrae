package server

import (
	"embed"
	"log/slog"
	"net/http"

	"github.com/vearutop/statigz"
)

//go:embed openapi/openapi.json
var oapi embed.FS

const specHTML = `<!doctype html>
<html>
  <head>
    <title>Scalar API Reference</title>
    <meta charset="utf-8" />
    <meta
      name="viewport"
      content="width=device-width, initial-scale=1" />
  </head>
  <body>
    <div id="app"></div>
    <script src="https://cdn.jsdelivr.net/npm/@scalar/api-reference"></script>
    <script>
      Scalar.createApiReference('#app', {
        url: '/openapi.json',
        proxyUrl: 'https://proxy.scalar.com',
		  theme: 'elysiajs',
      })
    </script>
  </body>
</html>`

func (s *Server) OpenAPIHandler() {
	// Serve OpenAPI spec
	s.mux.Handle("/openapi.json", statigz.FileServer(oapi, statigz.FSPrefix("openapi")))

	// Serve Spec UI
	s.mux.HandleFunc("/openapi", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		if _, err := w.Write([]byte(specHTML)); err != nil {
			slog.Error("failed to write elements HTML", "error", err)
		}
	})
}
