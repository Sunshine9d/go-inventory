-- ==========================
-- DROP EXISTING TABLES, PARTITIONS & TRIGGERS
-- ==========================

DO $$
DECLARE r RECORD;
BEGIN
    -- Drop all triggers first
FOR r IN (SELECT trigger_name, event_object_table FROM information_schema.triggers WHERE trigger_schema = 'public')
    LOOP
        EXECUTE 'DROP TRIGGER IF EXISTS ' || r.trigger_name || ' ON ' || r.event_object_table || ' CASCADE';
END LOOP;
END $$;

-- Drop dependent tables first
DROP TABLE IF EXISTS order_promotions CASCADE;
DROP TABLE IF EXISTS product_promotions CASCADE;
DROP TABLE IF EXISTS order_items CASCADE;
DROP TABLE IF EXISTS stock_adjustments CASCADE;

-- Drop main tables
DROP TABLE IF EXISTS orders CASCADE;
DROP TABLE IF EXISTS product_variants CASCADE;
DROP TABLE IF EXISTS products CASCADE;
DROP TABLE IF EXISTS promotions CASCADE;


-- ==========================
-- CREATE TABLE: products
-- ==========================
CREATE TABLE products (
                          id BIGSERIAL PRIMARY KEY,
                          name VARCHAR(255) NOT NULL,
                          category VARCHAR(50) NOT NULL,
                          sku UUID UNIQUE NOT NULL DEFAULT gen_random_uuid(),  -- ✅ Prevent duplicate SKUs
                          price DECIMAL(10,2) NOT NULL,
                          stock INT NOT NULL,
                          created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                          updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

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
                                  sku UUID UNIQUE NOT NULL DEFAULT gen_random_uuid(),
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
                        id BIGSERIAL,
                        customer_name VARCHAR(255) NOT NULL,
                        total_price DECIMAL(10,2) NOT NULL,
                        status VARCHAR(20) CHECK (status IN ('pending', 'shipped', 'delivered', 'cancelled')),
                        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                        PRIMARY KEY (id, created_at) -- Include partitioning column
) PARTITION BY RANGE (created_at);


-- ✅ Creating partitions for 2025
CREATE TABLE orders_2025_03 PARTITION OF orders
    FOR VALUES FROM ('2025-03-01') TO ('2025-03-31');

CREATE TABLE orders_default PARTITION OF orders DEFAULT;

-- ==========================
-- CREATE TABLE: order_items (Partitioned)
-- ==========================
CREATE TABLE order_items (
                             id BIGSERIAL,
                             order_id BIGINT NOT NULL,
                             order_created_at TIMESTAMP NOT NULL, -- Reference partitioned column
                             variant_id BIGINT NOT NULL REFERENCES product_variants(id) ON DELETE CASCADE,
                             quantity INT NOT NULL,
                             unit_price DECIMAL(10,2) NOT NULL,
                             product_name VARCHAR(255) NOT NULL,
                             product_category VARCHAR(50) NOT NULL,
                             product_sku VARCHAR(50) NOT NULL,
                             attributes JSONB,
                             created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                             PRIMARY KEY (id, created_at),
                             FOREIGN KEY (order_id, order_created_at) REFERENCES orders(id, created_at) ON DELETE CASCADE
) PARTITION BY RANGE (created_at);



-- ✅ Creating partitions for 2025
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
                                  order_id BIGINT NOT NULL,
                                  order_created_at TIMESTAMP NOT NULL, -- Reference partitioned column
                                  promotion_id BIGINT NOT NULL REFERENCES promotions(id) ON DELETE CASCADE,
                                  PRIMARY KEY (order_id, order_created_at, promotion_id),
                                  FOREIGN KEY (order_id, order_created_at) REFERENCES orders(id, created_at) ON DELETE CASCADE
);


-- ==========================
-- INSERT SAMPLE DATA
-- ==========================
-- Insert sample data into products (50 products)
INSERT INTO products (name, category, sku, price, stock)
SELECT
    'Product ' || i,
    'Category-' || FLOOR(RANDOM() * 10)::TEXT,
        gen_random_uuid(),  -- ✅ Generate valid UUID
    ROUND(CAST(RANDOM() * 500 + 50 AS NUMERIC), 2), -- Random price: 50 - 550
    FLOOR(RANDOM() * 100 + 10) -- Random stock: 10 - 110
FROM generate_series(1, 50000) AS i;



-- Insert sample data into product_variants (100 variants)
INSERT INTO product_variants (product_id, main_attribute_1, main_attribute_2, price, stock)
SELECT
    p.id,
    'Attribute1-' || FLOOR(RANDOM() * 10)::TEXT,
        'Attribute2-' || FLOOR(RANDOM() * 10)::TEXT,
        ROUND((RANDOM() * 50)::NUMERIC, 2) + p.price, -- Variant price based on product price
    FLOOR(RANDOM() * 50 + 5) -- Random stock: 5 - 55
FROM products p
         CROSS JOIN generate_series(1, 2) -- Each product gets 2 variants
    LIMIT 100000;


-- Insert sample data into orders (10000 orders)
INSERT INTO orders (customer_name, total_price, created_at)
SELECT
    'Customer ' || FLOOR(RANDOM() * 500000)::TEXT, -- More unique customers
        ROUND((RANDOM() * 1000 + 100)::NUMERIC, 2), -- Random price: 100 - 1100
    NOW() - INTERVAL '1 day' * FLOOR(RANDOM() * 365) -- Random date within the last year
FROM generate_series(1, 1000000);
