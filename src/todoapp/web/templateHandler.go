package web

import (
	"net/http"
	"regexp"
	"text/template"
)

type templateHandler struct {
}

func NewTemplateHandler() *templateHandler {
	return &templateHandler{}
}

var (
	//Switch this to parameters???
	ToDoRe       = regexp.MustCompile(`^/api/v1/lists/*$`)
	ToDoReWithID = regexp.MustCompile(`^/api/v1/lists/[a-zA-Z0-9]+$`)
)

func (h *templateHandler) serveTemplate(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("web/index.html"))
	tmpl.Execute(w, nil)
}

func (h *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.serveTemplate(w, r)
}
