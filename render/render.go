package render

import (
	"html/template"
	"net/http"

	"github.com/gorilla/schema"
	"github.com/ian-kent/go-webapp/assets"
	"github.com/ian-kent/go-webapp/data/user"
	"github.com/ian-kent/go-webapp/logger"
	"github.com/ian-kent/go-webapp/session"
	"github.com/ian-kent/htmlform"
	"github.com/justinas/nosurf"
	"github.com/unrolled/render"
	"gopkg.in/bluesuncorp/validator.v5"
)

// Render is a global instance of github.com/unrolled/render.Render
var Render = New()

// New creates a new instance of github.com/unrolled/render.Render
func New() *render.Render {
	logger.Traceln(nil, "creating instance of render.Render")
	// TODO make configurable?
	return render.New(render.Options{
		Asset:      assets.Asset,
		AssetNames: assets.AssetNames,
		Delims:     render.Delims{Left: "[:", Right: ":]"},
		Layout:     "layout",
		Funcs: []template.FuncMap{template.FuncMap{
			"map": htmlform.Map,
			"ext": htmlform.Extend,
			"fnn": htmlform.FirstNotNil,
		}},
	})
}

// HTML is an alias to github.com/unrolled/render.Render.HTML
func HTML(w http.ResponseWriter, status int, name string, binding interface{}, htmlOpt ...render.HTMLOptions) {
	Render.HTML(w, status, name, binding, htmlOpt...)
}

// DefaultVars adds the default vars (User and Session) to the data map
func DefaultVars(req *http.Request, m map[string]interface{}) map[string]interface{} {
	if m == nil {
		logger.Traceln(req, "creating new template data map")
		m = make(map[string]interface{})
	}

	s, _ := session.Get(req)
	if s == nil {
		logger.Traceln(req, "session not found, returning original map")
		return m
	}

	logger.Traceln(req, "adding session to template data map")
	m["Session"] = s

	if uid, ok := s.Values["User"]; ok {
		logger.Traceln(req, "user found in session, getting user")
		u, _ := user.Get(uid.(string))
		if u != nil {
			logger.Traceln(req, "adding user to template data map")
			m["User"] = u
			return m
		}
		logger.Traceln(req, "user not found, returning original map")
		return m
	}

	logger.Traceln(req, "user not found in session, returning original map")
	return m
}

// Decoder is a global gorilla schema decoder
var decoder = func() *schema.Decoder {
	var decoder = schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)
	decoder.ZeroEmpty(true) // allows HTML fields to be "cleared"
	return decoder
}()

// Vtom converts validator.v5 errors to a map[string]interface{}
func Vtom(req *http.Request, errs *validator.StructErrors) func(field string) map[string]interface{} {
	logger.Traceln(req, "flattening errors")
	ef := errs.Flatten()
	return func(field string) map[string]interface{} {
		logger.Tracef(req, "looking up form error for %s", field)
		if e, ok := ef[field]; ok {
			logger.Tracef(req, "got form error for %s: %s", field, e.Tag)
			return map[string]interface{}{e.Tag: e}
		}
		logger.Traceln(req, "error not found in map")
		return nil
	}
}

// Ftos converts form data into a model
func Ftos(req *http.Request, model interface{}) error {
	logger.Traceln(req, "parsing request form")
	err := req.ParseForm()

	if err != nil {
		logger.Errorf(req, "error parsing request form: %s", err)
		return err
	}

	err = decoder.Decode(model, req.PostForm)
	if err != nil {
		logger.Errorf(req, "error decoding request form: %s", err)
	}

	return err
}

// Form creates a htmlform.Form from a model and http.Request
func Form(req *http.Request, model interface{}, errs *validator.StructErrors) htmlform.Form {
	logger.Traceln(req, "creating form from model")
	return htmlform.Create(model, Vtom(req, errs), []string{}, []string{}).WithCSRF(nosurf.FormFieldName, nosurf.Token(req))
}
