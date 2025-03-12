-- ==========================
-- CREATE TABLE: products
-- ==========================
CREATE TABLE products (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    category VARCHAR(50) NOT NULL,
    sku VARCHAR(50) UNIQUE NOT NULL,
    price DECIMAL(10,2) NOT NULL,
    stock INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Add trigger to update 'updated_at' column
CREATE OR REPLACE FUNCTION update_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_products
BEFORE UPDATE ON products
FOR EACH ROW EXECUTE FUNCTION update_updated_at();

-- ==========================
-- CREATE TABLE: product_variants
-- ==========================
CREATE TABLE product_variants (
    id BIGSERIAL PRIMARY KEY,
    product_id BIGINT NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    sku VARCHAR(50) UNIQUE NOT NULL,
    main_attribute_1 VARCHAR(50),
    main_attribute_2 VARCHAR(50),
    attributes JSONB, -- Stores additional attributes
    price DECIMAL(10,2) NOT NULL,
    stock INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- ==========================
-- CREATE TABLE: orders (Partitioned)
-- ==========================
CREATE TABLE orders (
    id BIGSERIAL PRIMARY KEY,
    customer_name VARCHAR(255) NOT NULL,
    total_price DECIMAL(10,2) NOT NULL,
    status VARCHAR(20) CHECK (status IN ('pending', 'shipped', 'delivered', 'cancelled')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
) PARTITION BY RANGE (created_at);

-- Creating partitions for 2025
CREATE TABLE orders_2025_03 PARTITION OF orders
FOR VALUES FROM ('2025-03-01') TO ('2025-03-31');

CREATE TABLE orders_default PARTITION OF orders DEFAULT;

-- ==========================
-- CREATE TABLE: order_items (Partitioned)
-- ==========================
CREATE TABLE order_items (
    id BIGSERIAL PRIMARY KEY,
    order_id BIGINT NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    variant_id BIGINT NOT NULL REFERENCES product_variants(id) ON DELETE CASCADE,
    quantity INT NOT NULL,
    unit_price DECIMAL(10,2) NOT NULL,
    product_name VARCHAR(255) NOT NULL,
    product_category VARCHAR(50) NOT NULL,
    product_sku VARCHAR(50) NOT NULL,
    attributes JSONB,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
) PARTITION BY RANGE (created_at);

-- Creating partitions for 2025
CREATE TABLE order_items_2025_03 PARTITION OF order_items
FOR VALUES FROM ('2025-03-01') TO ('2025-03-31');

CREATE TABLE order_items_default PARTITION OF order_items DEFAULT;

-- ==========================
-- CREATE TABLE: stock_adjustments
-- ==========================
CREATE TABLE stock_adjustments (
   id BIGSERIAL PRIMARY KEY,
   product_id BIGINT NOT NULL REFERENCES products(id) ON DELETE CASCADE,
   adjustment_type VARCHAR(50) CHECK (adjustment_type IN ('restock', 'damage', 'sale', 'return')),
   quantity INT NOT NULL,
   reason TEXT,
   created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- ==========================
-- CREATE TABLE: promotions
-- ==========================
CREATE TABLE promotions (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    discount_percentage DECIMAL(5,2) CHECK (discount_percentage >= 0 AND discount_percentage <= 100),
    start_date TIMESTAMP NOT NULL,
    end_date TIMESTAMP NOT NULL
);

-- ==========================
-- CREATE TABLE: product_promotions (Many-to-Many)
-- ==========================
CREATE TABLE product_promotions (
    product_id BIGINT NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    promotion_id BIGINT NOT NULL REFERENCES promotions(id) ON DELETE CASCADE,
    PRIMARY KEY (product_id, promotion_id)
);

-- ==========================
-- CREATE TABLE: order_promotions (Many-to-Many)
-- ==========================
CREATE TABLE order_promotions (
    order_id BIGINT NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    promotion_id BIGINT NOT NULL REFERENCES promotions(id) ON DELETE CASCADE,
    PRIMARY KEY (order_id, promotion_id)
);


-- Insert sample data into products (50 products)
INSERT INTO products (name, category, sku, price, stock)
SELECT
    'Product ' || FLOOR(RANDOM() * 10000),
    'Category-' || FLOOR(RANDOM() * 10),
    'SKU-' || FLOOR(RANDOM() * 100000),
    ROUND(RANDOM() * 500 + 50, 2),
    FLOOR(RANDOM() * 100 + 10)
FROM generate_series(1, 50);

-- Insert sample data into product_variants (100 variants)
INSERT INTO product_variants (product_id, sku, main_attribute_1, main_attribute_2, price, stock)
SELECT
    p.id,
    'SKU-' || FLOOR(RANDOM() * 100000),
    'Attribute1-' || FLOOR(RANDOM() * 10),
    'Attribute2-' || FLOOR(RANDOM() * 10),
    p.price + ROUND(RANDOM() * 50, 2),
    FLOOR(RANDOM() * 50 + 5)
FROM products p
CROSS JOIN generate_series(1, 2) -- 2 variants per product
LIMIT 100;

-- Insert sample data into orders (100 orders)
INSERT INTO orders (customer_name, total_price, created_at)
SELECT
    'Customer ' || FLOOR(RANDOM() * 5000),
    ROUND(RANDOM() * 1000 + 100, 2),
    NOW() - INTERVAL '1 day' * FLOOR(RANDOM() * 365)
FROM generate_series(1, 100);

-- Insert sample data into order_items (500 order items)
INSERT INTO order_items (order_id, variant_id, quantity, unit_price, product_name, product_category, product_sku, attributes)
SELECT
    o.id,
    pv.id,
    FLOOR(RANDOM() * 5 + 1),  -- Quantity between 1 and 5
    pv.price,
    p.name,
    p.category,
    pv.sku,
    '{}'::jsonb -- Empty JSON for now
FROM orders o
JOIN product_variants pv ON pv.product_id = (SELECT id FROM products ORDER BY RANDOM() LIMIT 1)
CROSS JOIN generate_series(1, 5) -- Up to 5 items per order
LIMIT 500;
