package user

import (
	"github.com/gorilla/pat"
	"github.com/ian-kent/go-webapp/logger"
)

// Register creates routes for each static resource
func Register(r *pat.Router) {
	logger.Debugln(nil, "registering handlers for user package")
	r.Path("/signup").Methods("GET", "POST").HandlerFunc(SignUpHandler)
	r.Path("/login").Methods("GET", "POST").HandlerFunc(LoginHandler)
	r.Path("/logout").Methods("GET").HandlerFunc(LogoutHandler)
}
