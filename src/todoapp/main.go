package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/MatthewArwenLangham/todoapp/cli"
	"github.com/MatthewArwenLangham/todoapp/server"
	"github.com/MatthewArwenLangham/todoapp/store"
)

func main() {
	cliMode := flag.Bool("cli", false, "Enable CLI mode")

	flag.Parse()

	if *cliMode {
		cli.Run()
		return
	}

	store := store.NewMemStore()
	store.LoadFromFile()
	listHandler := server.NewToDoHandler(store)

	mux := http.NewServeMux()
	mux.Handle("/", listHandler)
	mux.Handle("/lists/", listHandler)

	fmt.Println("Server running on localhost:8010...")
	http.ListenAndServe(":8010", mux)
}
