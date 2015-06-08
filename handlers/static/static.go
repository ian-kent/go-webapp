package static

import (
	"mime"
	"net/http"
	"strings"

	"github.com/gorilla/pat"
	"github.com/ian-kent/go-webapp/assets"
	"github.com/ian-kent/go-webapp/logger"
	"github.com/ian-kent/go-webapp/render"
)

// Register creates routes for each static resource
func Register(r *pat.Router) {
	logger.Debugln(nil, "registering not found handler for static package")
	r.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		render.HTML(w, http.StatusNotFound, "error", map[string]interface{}{"error": "Page not found"})
	})

	logger.Debugln(nil, "registering static content handlers for static package")
	for _, file := range assets.AssetNames() {
		if strings.HasPrefix(file, "static/") {
			path := strings.TrimPrefix(file, "static")
			logger.Tracef(nil, "registering handler for static asset: %s", path)

			var mimeType string
			switch {
			case strings.HasSuffix(path, ".css"):
				mimeType = "text/css"
			case strings.HasSuffix(path, ".js"):
				mimeType = "application/javascript"
			default:
				mimeType = mime.TypeByExtension(path)
			}

			logger.Tracef(nil, "using mime type: %s", mimeType)

			r.Path(path).Methods("GET").HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				if b, err := assets.Asset("static" + path); err == nil {
					w.Header().Set("Content-Type", mimeType)
					w.Header().Set("Cache-control", "public, max-age=259200")
					w.WriteHeader(200)
					w.Write(b)
					return
				}
				// This should never happen!
				logger.Errorf(nil, "it happened ¯\\_(ツ)_/¯", path)
				r.NotFoundHandler.ServeHTTP(w, req)
			})
		}
	}
}
