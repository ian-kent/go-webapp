package home

import (
	"net/http"

	"github.com/gorilla/pat"
	"github.com/ian-kent/go-webapp/logger"
	"github.com/ian-kent/go-webapp/render"
)

// Register creates routes for each home handler
func Register(r *pat.Router) {
	logger.Debugln(nil, "registering handlers for home package")
	r.Path("/").Methods("GET").HandlerFunc(Handler)
}

// Handler is a http.Handler for the home page
func Handler(w http.ResponseWriter, req *http.Request) {
	render.HTML(w, http.StatusOK, "index", render.DefaultVars(req, nil))
}
