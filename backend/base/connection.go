package database

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
)

type Server struct {
	Database *Database
}
type Database struct {
	Host       string
	Port       string
	Username   string
	Password   string
	Name       string
	Connection *sql.DB
}
type Cache struct {
	Categories map[string]interface{}
	Products   map[string]interface{}
	Users      map[string]interface{}
}

func CacheInit() *Cache {
	return &Cache{
		Categories: make(map[string]interface{}),
		Products:   make(map[string]interface{}),
		Users:      make(map[string]interface{}),
	}
}

var CacheData *Cache = CacheInit()

const categoriesCacheTTL = 60 * time.Second

func (db *Database) InvalidateCategoriesCache() {
	delete(CacheData.Categories, "data")
	delete(CacheData.Categories, "ts")
}

func (db *Database) Connect() error {

	db.Host = os.Getenv("SIJILPOS_DB_HOST")
	if db.Host == "" {
		db.Host = "localhost"
	}
	db.Port = os.Getenv("SIJILPOS_DB_PORT")
	if db.Port == "" {
		db.Port = "5432"
	}

	db.Username = os.Getenv("SIJILPOS_DB_USER")
	if db.Username == "" {
		db.Username = "postgres"
	}
	db.Password = os.Getenv("SIJILPOS_DB_PASSWORD")
	if db.Password == "" {
		log.Fatalln("Database password is not set in environment variables")
	}
	db.Name = os.Getenv("SIJILPOS_DB_NAME")
	if db.Name == "" {
		db.Name = "sijilposdb"
	}
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		db.Host, db.Port, db.Username, db.Password, db.Name)
	var err error
	db.Connection, err = sql.Open("postgres", connStr)
	if err != nil {
		return err
	}
	log.Println("Connected to the database successfully")

	return db.Connection.Ping()
}
func (db *Database) GetProducts(w http.ResponseWriter, r *http.Request) {
	if r.Context().Err() == context.Canceled {
		http.Error(w, "Request canceled by the client", http.StatusRequestTimeout)
		return
	}
	rows, err := db.Connection.Query("SELECT id, name, price, stock, description, barcode, image FROM products")
	if err != nil {
		http.Error(w, "Failed to fetch products", http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	products := []map[string]interface{}{}
	for rows.Next() {
		var id int
		var name string
		var price float64
		var stock int
		var description sql.NullString
		var barcode sql.NullString
		var image sql.NullString
		if err := rows.Scan(&id, &name, &price, &stock, &description, &barcode, &image); err != nil {
			http.Error(w, "Failed to scan product", http.StatusInternalServerError)
			return
		}
		product := map[string]interface{}{
			"id":          id,
			"name":        name,
			"price":       price,
			"stock":       stock,
			"description": description,
			"barcode":     barcode,
			"image":       image,
		}
		products = append(products, product)
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(products)
}
func (db *Database) GetCategories(w http.ResponseWriter, r *http.Request) {
	if r.Context().Err() == context.Canceled {
		http.Error(w, "Request canceled by the client", http.StatusRequestTimeout)
		return
	}

	if cached, ok := CacheData.Categories["data"]; ok {
		if tsRaw, okTs := CacheData.Categories["ts"]; okTs {
			if ts, okTime := tsRaw.(time.Time); okTime {
				if time.Since(ts) < categoriesCacheTTL {
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(cached)
					return
				}
			}
		}
	}
	rows, err := db.Connection.Query("SELECT category_id, name, created_at FROM categories")
	if err != nil {
		http.Error(w, "Failed to fetch categories", http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	categories := []map[string]interface{}{}
	for rows.Next() {
		var categoryID int
		var name string
		var created_at time.Time
		if err := rows.Scan(&categoryID, &name, &created_at); err != nil {
			http.Error(w, "Failed to scan category", http.StatusInternalServerError)
			return
		}
		category := map[string]interface{}{
			"id":         categoryID,
			"name":       name,
			"created_at": created_at,
		}
		categories = append(categories, category)
	}

	CacheData.Categories["data"] = categories
	CacheData.Categories["ts"] = time.Now()
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(categories)
}
