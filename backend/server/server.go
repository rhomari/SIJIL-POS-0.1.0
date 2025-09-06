package server

import (
	"database"
	"log"
	"net/http"
	"os"
)

type Server struct {
	Database *database.Database
}

func NewServer() *Server {
	db := &database.Database{}
	err := db.Connect()
	if err != nil {
		panic(err)
	}
	return &Server{
		Database: db,
	}
}
func (server *Server) Start() error {

	server.SetupRoutes()

	port := os.Getenv("SIJILPOS_SERVER_PORT")
	if port == "" {
		port = "8080"
	}
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Println("Failed to start server:", err)
	}
	log.Printf("Server is running on port %s\n", port)

	return nil
}
func (server *Server) Stop() {

	println("Server stopped")
}
func (server *Server) Restart() {

	server.Stop()
	server.Start()
	println("Server restarted")
}
func (server *Server) SetupRoutes() {

	http.HandleFunc("/", http.FileServer(http.Dir("./static")).ServeHTTP)
	http.HandleFunc("/users", server.UsersHandler)
	http.HandleFunc("/products", server.ProductsHandler)
	http.HandleFunc("/orders", server.OrdersHandler)
}
