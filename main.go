package main

import (
	"log"
	"net/http"
	"time"

	"github.com/ian-kent/go-webapp/config"
	"github.com/ian-kent/go-webapp/handlers/home"
	"github.com/ian-kent/go-webapp/handlers/static"
	"github.com/ian-kent/go-webapp/handlers/timeout"
	"github.com/ian-kent/go-webapp/handlers/user"
	"github.com/ian-kent/go-webapp/logger"
	"github.com/ian-kent/go-webapp/session"

	"github.com/gorilla/pat"
	"github.com/ian-kent/go-webapp/render"
	"github.com/justinas/alice"
	"github.com/justinas/nosurf"
)

func main() {
	cfg, err := config.Get()
	if err != nil {
		log.Fatalf("error configuring app: %s", err)
	}

	// logging before this point must rely on setting LOGLEVEL env var
	if l, err := logger.LevelFromString(cfg.LogLevel); err == nil {
		logger.DefaultLevel = l
	} else {
		log.Fatalf("error setting log level: %s", err)
	}

	router := pat.New()
	static.Register(router)
	home.Register(router)
	user.Register(router)

	session.Init(cfg.SessionSecret, cfg.SessionName)

	chain := alice.New(logger.Handler /*, context.ClearHandler*/, timeoutHandler, withCsrf).Then(router)

	logger.Infof(nil, "listening on %s", cfg.BindAddr)
	err = http.ListenAndServe(cfg.BindAddr, chain)
	if err != nil {
		log.Fatalf("error listening on %s: %s", cfg.BindAddr, err)
	}
}

func withCsrf(h http.Handler) http.Handler {
	csrfHandler := nosurf.New(h)
	csrfHandler.SetFailureHandler(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		rsn := nosurf.Reason(req).Error()
		logger.Warnf(req, "failed csrf validation: %s", rsn)
		render.HTML(w, http.StatusBadRequest, "error", map[string]interface{}{"error": rsn})
	}))
	return csrfHandler
}

func timeoutHandler(h http.Handler) http.Handler {
	return timeout.Handler(h, 1*time.Second, http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		logger.Warnln(req, "request timed out")
		render.HTML(w, http.StatusRequestTimeout, "error", map[string]interface{}{"error": "Request timed out"})
	}))
}
