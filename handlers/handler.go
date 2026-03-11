package handlers

import (
	"builderwireapi/models"
	"database/sql"
	"strings"

	"github.com/gofiber/fiber/v2"
	_ "github.com/lib/pq"
)

var DB *sql.DB

// Gets connection pool from db.go, which is passed in by main.go
func SetDB(db *sql.DB) {
	DB = db
}

func CreateProduct(c *fiber.Ctx) error {
	var (
		product    models.Product
		insertedID int
	)
	err := c.BodyParser(&product)

	//Input validation
	product.SKU = strings.TrimSpace(product.SKU)

	if len(product.SKU) == 0 {
		return c.Status(400).JSON(fiber.Map{"error": "SKU must not be empty"})
	}
	if len(product.SKU) > 10 {
		return c.Status(400).JSON(fiber.Map{"error": "SKU must not exceed 10 characters"})
	}
	if product.QuantityOnHand < 0 {
		return c.Status(400).JSON(fiber.Map{"error": "Quantity on hand must be 0 or more"})
	}
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	// var insertedID int
	err = DB.QueryRow("SELECT id FROM products WHERE sku=$1", product.SKU).Scan(&insertedID) // Checks if inserted product already exists. Refer to README for the reason
	if err != nil {
		if err == sql.ErrNoRows {
			_, err = DB.Exec(`
			INSERT INTO products(sku, quantity_on_hand)
			VALUES($1, $2)`, product.SKU, product.QuantityOnHand)

			if err != nil {
				return c.Status(500).JSON(fiber.Map{"error": "Error while inserting data into database"})
			}
			return c.Status(201).JSON(product)
		} else {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to create product. Please try again."})
		}
	} else {
		return c.Status(400).JSON(fiber.Map{"error": "Product already exists"})
	}
}

func PlaceOrder(c *fiber.Ctx) error {
	var (
		Order       models.Order
		retrievedID int
	)

	err := c.BodyParser(&Order)
	// Input validation
	if Order.Quantity <= 0 {
		return c.Status(400).JSON(fiber.Map{"error": "Quantity must be more than zero"})
	}

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	err = DB.QueryRow(`SELECT id FROM products WHERE id=$1`, Order.ProductID).Scan(&retrievedID) //Checks if product_id form the request exists in products table. Although the foreign key constraint takes care of this, this is for user-friendly errors

	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(400).JSON(fiber.Map{"error": "ID not found"})
		} else {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
	}

	//Start transaction
	tx, err := DB.Begin()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	// Lock the row
	var currentQty int
	err = tx.QueryRow("SELECT quantity_on_hand FROM products WHERE id = $1 FOR UPDATE", retrievedID).Scan(&currentQty)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	//Check for sufficient inventory
	if currentQty >= Order.Quantity {
		//Updates quantity in products table
		_, err = tx.Exec("UPDATE products SET quantity_on_hand = quantity_on_hand - $1 WHERE id = $2", Order.Quantity, retrievedID)
		if err != nil {
			tx.Rollback()
			return c.Status(500).JSON(fiber.Map{"error": "Couldn't update quantity on hand"})
		}

		//Inserts log into order table
		_, err = tx.Exec("INSERT INTO orders(product_id, quantity) VALUES($1, $2)", Order.ProductID, Order.Quantity)
		if err != nil {
			tx.Rollback()
			return c.Status(500).JSON(fiber.Map{"error": "Error while inserting data into database"})
		}

		err = tx.Commit()
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}

		return c.Status(200).JSON(fiber.Map{"Success": "Order placed successfully"})
	} else {
		tx.Rollback()
		return c.Status(400).JSON(fiber.Map{"error": "Not enough in inventory"})
	}
}

func HealthCheck(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "API OK!"})
}

