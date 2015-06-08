package user

import (
	"net/http"

	"github.com/ian-kent/go-webapp/data/user"
	"github.com/ian-kent/go-webapp/logger"
	"github.com/ian-kent/go-webapp/render"
	"github.com/ian-kent/go-webapp/session"
	"gopkg.in/bluesuncorp/validator.v5"
)

type signUpDataModel struct {
	Email    string `validate:"required,email,max=60" schema:"email"`
	Password string `validate:"required,min=8" schema:"password" htmlform:"type=password"`
}

// SignUpHandler is a http.Handler for the sign up page
func SignUpHandler(w http.ResponseWriter, req *http.Request) {
	var model signUpDataModel
	var errs *validator.StructErrors

	sess, _ := session.Get(req)

	if _, ok := sess.Values["User"]; ok {
		logger.Traceln(req, "user signed in, redirecting")
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}

	if req.Method == "POST" {
		if err := render.Ftos(req, &model); err != nil {
			logger.Errorf(req, "error parsing form: %s", err)
			render.HTML(w, http.StatusBadRequest, "error", err)
			return
		}

		if errs = validator.New("validate", validator.BakedInValidators).Struct(&model); errs == nil {
			logger.Tracef(req, "form is valid, creating user %s", model.Email)
			u, err := user.Create(model.Email, model.Password)
			if err != nil {
				logger.Errorf(req, "error creating user: %s", err)
				render.HTML(w, http.StatusInternalServerError, "error", err)
				return
			}

			sess.Values["User"] = u.Email
			sess.Save(req, w)

			logger.Traceln(req, "user created, redirecting")
			http.Redirect(w, req, "/", http.StatusSeeOther)
			return
		}

		logger.Tracef(req, "form errors: %s", errs)
	}

	form := render.Form(req, &model, errs)

	render.HTML(w, http.StatusOK, "user/signup", render.DefaultVars(req, map[string]interface{}{"Form": form}))
}
