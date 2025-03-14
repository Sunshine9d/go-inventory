package postgres

import (
	"database/sql"
	"log"
)

// MigratePostgres runs the database migrations for PostgreSQL
func MigratePostgres(db *sql.DB) error {
	query := `
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

	-- Create products table
	CREATE TABLE products (
	    id BIGSERIAL PRIMARY KEY,
	    name VARCHAR(255) NOT NULL,
	    category VARCHAR(50) NOT NULL,
	    sku UUID UNIQUE NOT NULL DEFAULT gen_random_uuid(),
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

	-- Create product_variants table
	CREATE TABLE product_variants (
	    id BIGSERIAL PRIMARY KEY,
	    product_id BIGINT NOT NULL REFERENCES products(id) ON DELETE CASCADE,
	    sku UUID UNIQUE NOT NULL DEFAULT gen_random_uuid(),
	    main_attribute_1 VARCHAR(50),
	    main_attribute_2 VARCHAR(50),
	    attributes JSONB,
	    price DECIMAL(10,2) NOT NULL,
	    stock INT NOT NULL,
	    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	-- Create partitioned orders table
	CREATE TABLE orders (
	    id BIGSERIAL,
	    customer_name VARCHAR(255) NOT NULL,
	    total_price DECIMAL(10,2) NOT NULL,
	    status VARCHAR(20) CHECK (status IN ('pending', 'shipped', 'delivered', 'cancelled')),
	    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	    PRIMARY KEY (id, created_at)
	) PARTITION BY RANGE (created_at);

	CREATE TABLE orders_2025_03 PARTITION OF orders
	    FOR VALUES FROM ('2025-03-01') TO ('2025-03-31');

	CREATE TABLE orders_default PARTITION OF orders DEFAULT;

	-- Create order_items table
	CREATE TABLE order_items (
	    id BIGSERIAL,
	    order_id BIGINT NOT NULL,
	    order_created_at TIMESTAMP NOT NULL,
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

	CREATE TABLE order_items_2025_03 PARTITION OF order_items
	    FOR VALUES FROM ('2025-03-01') TO ('2025-03-31');

	CREATE TABLE order_items_default PARTITION OF order_items DEFAULT;

	-- Create stock_adjustments table
	CREATE TABLE stock_adjustments (
	    id BIGSERIAL PRIMARY KEY,
	    product_id BIGINT NOT NULL REFERENCES products(id) ON DELETE CASCADE,
	    adjustment_type VARCHAR(50) CHECK (adjustment_type IN ('restock', 'damage', 'sale', 'return')),
	    quantity INT NOT NULL,
	    reason TEXT,
	    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	-- Create promotions table
	CREATE TABLE promotions (
	    id BIGSERIAL PRIMARY KEY,
	    name VARCHAR(255) NOT NULL,
	    discount_percentage DECIMAL(5,2) CHECK (discount_percentage >= 0 AND discount_percentage <= 100),
	    start_date TIMESTAMP NOT NULL,
	    end_date TIMESTAMP NOT NULL
	);

	-- Create product_promotions table
	CREATE TABLE product_promotions (
	    product_id BIGINT NOT NULL REFERENCES products(id) ON DELETE CASCADE,
	    promotion_id BIGINT NOT NULL REFERENCES promotions(id) ON DELETE CASCADE,
	    PRIMARY KEY (product_id, promotion_id)
	);

	-- Create order_promotions table
	CREATE TABLE order_promotions (
	    order_id BIGINT NOT NULL,
	    order_created_at TIMESTAMP NOT NULL,
	    promotion_id BIGINT NOT NULL REFERENCES promotions(id) ON DELETE CASCADE,
	    PRIMARY KEY (order_id, order_created_at, promotion_id),
	    FOREIGN KEY (order_id, order_created_at) REFERENCES orders(id, created_at) ON DELETE CASCADE
	                              
	);
	`

	_, err := db.Exec(query)
	if err != nil {
		log.Printf("❌ Failed to run PostgreSQL migrations: %v\n", err)
		return err
	}

	log.Println("✅ PostgreSQL migrations completed successfully!")
	return nil
}
