CREATE TABLE IF NOT EXISTS products (
    id SERIAL PRIMARY KEY,
    sku VARCHAR(10) UNIQUE NOT NULL,
    quantity_on_hand INT CHECK (quantity_on_hand >= 0) NOT NULL
);

CREATE TABLE IF NOT EXISTS orders (
    id SERIAL PRIMARY KEY,
    product_id INT NOT NULL,
    quantity INT CHECK (quantity > 0) NOT NULL,
    FOREIGN KEY (product_id) REFERENCES products(id)
);