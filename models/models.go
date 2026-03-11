package models

// models for products POST requests
type Product struct {
	SKU            string `json:"sku"`
	QuantityOnHand int    `json:"quantity_on_hand"`
}

// models for products POST requests
type Order struct {
	ProductID int `json:"product_id"`
	Quantity  int `json:"quantity"`
}
