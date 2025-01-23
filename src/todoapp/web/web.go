package web

import (
	"fmt"
	"net/http"
)

func Start(mux *http.ServeMux) {
	templateHandler := NewTemplateHandler()

	mux.Handle("/", templateHandler)
	mux.Handle("/lists/", templateHandler)
	fmt.Println("Template handler started!")
}
