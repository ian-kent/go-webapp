package user

import (
	"net/http"

	"github.com/ian-kent/go-webapp/logger"
	"github.com/ian-kent/go-webapp/session"
)

// LogoutHandler is a http.Handler for the logout page
func LogoutHandler(w http.ResponseWriter, req *http.Request) {
	sess, _ := session.Get(req)

	if _, ok := sess.Values["User"]; !ok {
		logger.Traceln(req, "user not signed in")
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}

	delete(sess.Values, "User")
	sess.Save(req, w)

	logger.Traceln(req, "user signed out, redirecting")
	http.Redirect(w, req, "/", http.StatusSeeOther)
}
