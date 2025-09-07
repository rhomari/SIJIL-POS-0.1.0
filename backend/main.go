package main

import (
	"SIJIL-POS/server"
)

func main() {
	srv := server.NewServer()
	srv.SetupRoutes()

	srv.Start()
}
