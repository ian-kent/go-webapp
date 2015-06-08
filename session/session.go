package session

import (
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/ian-kent/go-webapp/logger"
)

type store struct {
	*sessions.CookieStore
	name string
}

func (s store) get(r *http.Request) (se *sessions.Session, e error) {
	logger.Traceln(r, "getting session")
	se, e = s.CookieStore.Get(r, s.name)
	if e != nil {
		logger.Tracef(r, "error getting session: %s", e)
		return
	}
	logger.Tracef(r, "got session %s", se.ID)
	return se, e
}

var s *store

// Init initialises the session store
func Init(secret []byte, name string) {
	logger.Infoln(nil, "initialising session storage")
	s = &store{sessions.NewCookieStore(secret), name}
}

// Get returns a session from the session store
func Get(r *http.Request) (*sessions.Session, error) {
	return s.get(r)
}
