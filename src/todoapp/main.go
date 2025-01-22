package main

import (
	"flag"

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
	server.Start(store)
}
