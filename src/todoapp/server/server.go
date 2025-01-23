package server

import (
	"fmt"
	"net/http"

	"github.com/MatthewArwenLangham/todoapp/store"
)

func Start(mux *http.ServeMux, store store.Store) {
	listHandler := NewToDoHandler(store)

	mux.Handle("/api/v1", listHandler)
	mux.Handle("/api/v1/lists/", listHandler)

	fmt.Println("Server running on localhost:8010...")
	http.ListenAndServe(":8010", mux)
}
