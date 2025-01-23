package web

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"regexp"
	"text/template"

	"github.com/MatthewArwenLangham/todoapp/store"
)

type templateHandler struct {
}

func NewTemplateHandler() *templateHandler {
	return &templateHandler{}
}

var (
	//Switch this to parameters???
	ToDoRe       = regexp.MustCompile(`^/lists/*$`)
	ToDoReWithID = regexp.MustCompile(`^/lists/[a-zA-Z0-9]+$`)
)

func (h *templateHandler) getHome(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles(
		"web/index.html",
		"web/header.html",
		"web/footer.html",
	))
	tmpl.Execute(w, nil)
}

func (h *templateHandler) getLists(w http.ResponseWriter, r *http.Request) {
	response, err := http.Get("http://127.0.0.1:8010/api/v1/lists")
	if err != nil {
		log.Fatalf("Failed to GET from api/v1/lists: %v", err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatalf("Failed to read response body: %v", err)
	}

	var lists map[string]store.List
	if err := json.Unmarshal(body, &lists); err != nil {
		log.Fatalf("Failed to parse JSON: %v", err)
	}

	tmpl := template.Must(template.ParseFiles(
		"web/lists.html",
		"web/header.html",
		"web/footer.html"))
	tmpl.Execute(w, lists)
}

func (h *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodGet && ToDoRe.MatchString(r.URL.Path):
		h.getLists(w, r)
		return
	case r.Method == http.MethodGet && r.URL.Path == "/":
		h.getHome(w, r)
		return
	}
}
