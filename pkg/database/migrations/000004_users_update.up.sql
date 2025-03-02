-- Table Roles
CREATE TABLE roles (
    id SERIAL PRIMARY KEY,
    role_name VARCHAR(20) UNIQUE NOT NULL
);

INSERT INTO roles (role_name) VALUES 
('admin'),
('seller'),
('customer');
-- End Table Roles

-- Table Users
ALTER TABLE users RENAME COLUMN picture TO images;
ALTER TABLE users ADD COLUMN first_name VARCHAR(30);
ALTER TABLE users ADD COLUMN last_name VARCHAR(30);
ALTER TABLE users DROP COLUMN role;
ALTER TABLE users ADD COLUMN role_id INT REFERENCES roles(id) ON DELETE SET NULL;
ALTER TABLE users ALTER COLUMN created_at DROP DEFAULT;
ALTER TABLE users ALTER COLUMN updated_at DROP DEFAULT;

-- Drop Old ID
ALTER TABLE orders DROP CONSTRAINT orders_order_by_fkey;
ALTER TABLE carts DROP CONSTRAINT carts_order_by_fkey;
ALTER TABLE refresh_token DROP CONSTRAINT refresh_token_user_id_fkey;
ALTER TABLE users DROP COLUMN id;

-- Update New User ID
ALTER TABLE users ADD COLUMN user_id VARCHAR(6) UNIQUE;

UPDATE users SET user_id = 'U' || LPAD(CAST(row_number AS TEXT), 5, '0')
FROM (
    SELECT user_id, ROW_NUMBER() OVER (ORDER BY created_at) AS row_number
    FROM users
) AS temp
WHERE users.user_id = temp.user_id;

ALTER TABLE users ADD PRIMARY KEY (user_id);

CREATE SEQUENCE user_seq START 1;

ALTER TABLE users ALTER COLUMN user_id
SET DEFAULT 'U' || LPAD(NEXTVAL('user_seq')::TEXT, 5, '0');

