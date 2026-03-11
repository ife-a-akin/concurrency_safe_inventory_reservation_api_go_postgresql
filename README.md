Overview
This project implements a simple order placement API with inventory protection against overselling under concurrent requests.

Setup Instructions
1. Ensure Docker and Docker Compose are installed.

2. Start the database: 

docker compose up

3. Run the Go server:

go run main.go

4. API runs on:

http://localhost:8000


API Endpoints

POST /createproduct

Example request:
{
  "sku": "5X6-2FT",
  "quantity_on_hand": 370
}

Product flow:
Request → Validate Input → Check if product's SKU already exists
→ Insert product


POST /placeorder

Example request:
{
  "product_id": 9,
  "quantity": 100
}

Order flow:
Request → Validate Input → Begin Transaction
→ Lock Product Row (SELECT FOR UPDATE)
→ Verify Inventory → Update Inventory
→ Commit Transaction → Insert Order


GET /health - To confirm if API is running


Design Decisions

Queries are parameterized, and placeholders are used to avoid SQL injection

Inventory updates are handled using a database transaction with SELECT ... FOR UPDATE.

This ensures that when multiple concurrent requests attempt to place orders for the same product, the product row is locked until the transaction completes. This prevents race conditions and guarantees that inventory cannot be oversold.

Constraints have also been put on quantity_on_hand in products and quantity in orders to be 0 or more and more than 0 respectively

Added a check for entries in the products table even though the 'UNIQUE' SQL constraint was applied to the field

This ensured the next successful insert will not skip over the ID that was 'hidden' by the failed insert, hence, avoiding 'gaps' in the ID.


Assumptions

Products must be created before orders can be placed.
Inventory represents the currently available stock.


Improvements With More Time

Add timestamp to orders for better logginh
Add additional endpoints for listing products and orders.
Add basic indexing to make lookup faster