-- Drop all tables in reverse order (respecting foreign key dependencies)

-- Drop payments table
DROP INDEX IF EXISTS idx_payments_status;
DROP INDEX IF EXISTS idx_payments_user_id;
DROP INDEX IF EXISTS idx_payments_order_id;
DROP TABLE IF EXISTS payments;

-- Drop order_items table
DROP INDEX IF EXISTS idx_order_items_product_id;
DROP INDEX IF EXISTS idx_order_items_order_id;
DROP TABLE IF EXISTS order_items;

-- Drop orders table
DROP INDEX IF EXISTS idx_orders_status;
DROP INDEX IF EXISTS idx_orders_user_id;
DROP TABLE IF EXISTS orders;

-- Drop products table
DROP INDEX IF EXISTS idx_products_price;
DROP INDEX IF EXISTS idx_products_name;
DROP TABLE IF EXISTS products;
