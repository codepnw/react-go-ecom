ALTER TABLE products RENAME COLUMN sole TO stock;

ALTER TABLE products ADD COLUMN category_id INT REFERENCES categories(id) ON DELETE SET NULL;

UPDATE products SET updated_at = NULL;
ALTER TABLE products ALTER COLUMN updated_at DROP DEFAULT; 

-- Drop Old Product ID & FK
ALTER TABLE product_order DROP CONSTRAINT product_order_product_id_fkey;
ALTER TABLE product_cart DROP CONSTRAINT product_cart_product_id_fkey;
ALTER TABLE products DROP CONSTRAINT products_pkey;
ALTER TABLE products DROP COLUMN id;

-- Update New Product ID
ALTER TABLE products ADD COLUMN product_id VARCHAR(10) UNIQUE;

UPDATE products SET product_id = 'P' || LPAD(CAST(row_number AS TEXT), 5, '0')
FROM (
    SELECT product_id, ROW_NUMBER() OVER (ORDER BY created_at) AS row_number
    FROM products
) AS temp
WHERE products.product_id = temp.product_id;

ALTER TABLE products ADD PRIMARY KEY (product_id);

CREATE SEQUENCE product_seq START 1;

ALTER TABLE products ALTER COLUMN product_id 
SET DEFAULT 'P' || LPAD(NEXTVAL('product_seq')::TEXT, 5, '0');