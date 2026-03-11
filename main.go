package main

import (
	// "fmt"
	"log"

	"builderwireapi/db"
	"builderwireapi/handlers"
	"database/sql"

	"github.com/gofiber/fiber/v2"
)

// Create products table
func setupProductsTable(db *sql.DB) error {
	_, err := db.Exec(`
        CREATE TABLE IF NOT EXISTS products (
            id SERIAL PRIMARY KEY,
            sku VARCHAR(10) UNIQUE NOT NULL,
            quantity_on_hand INT CHECK (quantity_on_hand >= 0) NOT NULL
        )
    `)
	return err
}

// Create orders table
func setupOrderTable(db *sql.DB) error {
	_, err := db.Exec(`
        CREATE TABLE IF NOT EXISTS orders (
            id SERIAL PRIMARY KEY,
            product_id INT NOT NULL,
            quantity INT CHECK (quantity > 0) NOT NULL,
			FOREIGN KEY (product_id) REFERENCES products(id)
        )
    `)
	return err
}

func main() {
	db.ConnectDB()

	err := setupProductsTable(db.DB)
	if err != nil {
		log.Fatal("Database setup failed: ", err)
	}

	err = setupOrderTable(db.DB)
	if err != nil {
		log.Fatal("Order table setup failed: ", err)
	}

	//Makes the database connection pool from db.go accessible in handlers.go
	handlers.SetDB(db.DB)

	app := fiber.New() //Fiber instance initialized

	app.Post("/products", handlers.CreateProduct)
	app.Post("/orders", handlers.PlaceOrder)
	app.Get("/health", handlers.HealthCheck)

	log.Println("Server running...")
	app.Listen(":8000") //Server runs on this port as seen in README
}
