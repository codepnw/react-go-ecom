SET
    TIME ZONE 'Asia/Bangkok';

CREATE TYPE user_role AS ENUM ('user', 'admin', 'owner');

CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(100) UNIQUE NOT NULL,
    password VARCHAR(100) NOT NULL,
    picture VARCHAR(100) NOT NULL,
    role user_role NOT NULL DEFAULT 'user',
    enabled BOOLEAN DEFAULT true,
    address VARCHAR(255),
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now()
);

CREATE TABLE IF NOT EXISTS products (
    id SERIAL PRIMARY KEY,
    title VARCHAR(100) NOT NULL,
    description VARCHAR(100),
    price FLOAT NOT NULL,
    sole INT NULL NULL DEFAULT 0,
    quantity INT NULL NULL DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now()
);

CREATE TYPE order_status AS ENUM ('waiting', 'shipping', 'completed', 'canceled');

CREATE TABLE IF NOT EXISTS orders (
    id SERIAL PRIMARY KEY,
    cart_total FLOAT NOT NULL,
    order_by INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    amount INT,
    status order_status DEFAULT 'waiting',
    currentcy VARCHAR(10) DEFAULT 'THB',
    created_at TIMESTAMPTZ DEFAULT now()
);

CREATE TABLE IF NOT EXISTS product_order (
    id SERIAL PRIMARY KEY,
    product_id INT NULL NULL REFERENCES products(id) ON DELETE CASCADE,
    order_id INT NULL NULL REFERENCES orders(id) ON DELETE CASCADE,
    count INT,
    price FLOAT
);

CREATE TABLE IF NOT EXISTS categories (
    id SERIAL PRIMARY KEY,
    title VARCHAR(50) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS carts (
    id SERIAL PRIMARY KEY,
    cart_total INT NOT NULL,
    order_by INT NOT NULL REFERENCES users(id),
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now()
);

CREATE TABLE IF NOT EXISTS product_cart (
    id SERIAL PRIMARY KEY,
    cart_id INT NOT NULL REFERENCES carts(id) ON DELETE CASCADE,
    product_id INT NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    count INT NOT NULL,
    price FLOAT NOT NULL
);