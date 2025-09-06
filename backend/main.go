package main

import (
	"github.com/rhomari/Sijil-POS/backend/server"
)

func main() {
	srv := server.NewServer()
	srv.Start()
}
