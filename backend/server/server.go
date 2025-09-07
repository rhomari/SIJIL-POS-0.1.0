package server

import (
	database "SIJIL-POS/base"
	middleware "SIJIL-POS/middleware"
	utils "SIJIL-POS/utils"
	"log"
	"net/http"
	"os"
)

type Server struct {
	Database   *database.Database
	Middleware *middleware.Middleware
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

	port := os.Getenv("SIJILPOS_SERVER_PORT")
	if port == "" {
		port = "8080"
	}
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Println("Failed to start server:", err)
	}
	log.Printf(utils.Blue+"Server is running on port %s\n"+utils.Reset, port)

	return nil
}

func (server *Server) SetupRoutes() {

	http.HandleFunc("/", http.FileServer(http.Dir("./static")).ServeHTTP)

	http.Handle("/api/products", server.Middleware.Serve(http.HandlerFunc(server.handleProducts)))
	http.Handle("/api/categories", server.Middleware.Serve(http.HandlerFunc(server.handleCategories)))

}
func (server *Server) SetHeaders(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
}
func (server *Server) handleProducts(w http.ResponseWriter, r *http.Request) {
	log.Println(utils.Yellow + "Handling /api/products request" + utils.Reset)
	server.SetHeaders(w)
	switch r.Method {
	case http.MethodGet:
		server.Database.GetProducts(w, r)
	case http.MethodPost:
		//server.Database.CreateProduct(w, r)
	}
}

func (server *Server) handleCategories(w http.ResponseWriter, r *http.Request) {
	log.Println("Handling /api/categories request")
	server.SetHeaders(w)
	switch r.Method {
	case http.MethodGet:
		server.Database.GetCategories(w, r)
	case http.MethodPost:
		// Implement category creation logic here, e.g.:
		// server.Database.CreateCategory(w, r)
	}
}
