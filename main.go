package main

import (
	"flag"
	"net/http"

	"github.com/MatthewArwenLangham/todoapp/cli"
	"github.com/MatthewArwenLangham/todoapp/server"
	"github.com/MatthewArwenLangham/todoapp/store"
	"github.com/MatthewArwenLangham/todoapp/web"
)

func main() {
	cliMode := flag.Bool("cli", false, "Enable CLI mode")

	flag.Parse()

	if *cliMode {
		cli.Run()
		return
	}
	mux := http.NewServeMux()
	web.Start(mux)
	store := store.NewMemStore()
	go store.LoadFromFile()
	server.Start(mux, store)
}
