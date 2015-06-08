package user

import (
	"net/http"

	"github.com/ian-kent/go-webapp/data/user"
	"github.com/ian-kent/go-webapp/logger"
	"github.com/ian-kent/go-webapp/render"
	"github.com/ian-kent/go-webapp/session"
	"gopkg.in/bluesuncorp/validator.v5"
)

type loginDataModel struct {
	Email    string `validate:"required,email,max=60" schema:"email"`
	Password string `validate:"required,min=8" schema:"password" htmlform:"type=password"`
}

// LoginHandler is a http.Handler for the sign up page
func LoginHandler(w http.ResponseWriter, req *http.Request) {
	var model loginDataModel
	var errs *validator.StructErrors
	var loginError string

	sess, _ := session.Get(req)

	if _, ok := sess.Values["User"]; ok {
		logger.Traceln(req, "user already signed in, redirecting")
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}

	if req.Method == "POST" {
		if err := render.Ftos(req, &model); err != nil {
			logger.Tracef(req, "error parsing form: %s", err)
			render.HTML(w, http.StatusBadRequest, "error", err)
			return
		}

		if errs = validator.New("validate", validator.BakedInValidators).Struct(&model); errs == nil {
			logger.Tracef(req, "form is valid, getting user %s", model.Email)
			u, err := user.Get(model.Email)
			if err != nil {
				logger.Tracef(req, "error fetching user: %s", err)
				render.HTML(w, http.StatusInternalServerError, "error", err)
				return
			}

			if u == nil {
				logger.Tracef(req, "user not found")
				loginError = "invalid-email"
			} else {
				ok, err := u.ValidatePassword(model.Password)
				if err != nil {
					logger.Tracef(req, "error validating password: %s", err)
					render.HTML(w, http.StatusInternalServerError, "error", err)
					return
				}

				if !ok {
					logger.Traceln(req, "invalid password")
					loginError = "invalid-password"
				} else {
					logger.Traceln(req, "user signed in")
					sess.Values["User"] = u.Email
					sess.Save(req, w)
					http.Redirect(w, req, "/", http.StatusSeeOther)
					return
				}
			}
		}

		logger.Tracef(req, "form errors: %s", errs)
	}

	form := render.Form(req, &model, errs)

	render.HTML(w, http.StatusOK, "user/login", render.DefaultVars(req, map[string]interface{}{"Form": form, "Error": loginError}))
}
