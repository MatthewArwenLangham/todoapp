package server

import (
	"fmt"
	"net/http"

	"github.com/MatthewArwenLangham/todoapp/store"
)

func Start(store store.Store) {
	listHandler := NewToDoHandler(store)

	mux := http.NewServeMux()
	mux.Handle("/", listHandler)
	mux.Handle("/lists/", listHandler)

	fmt.Println("Server running on localhost:8010...")
	http.ListenAndServe(":8010", mux)
}
