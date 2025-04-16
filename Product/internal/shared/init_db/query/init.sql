CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS products_query (
                                              id UUID PRIMARY KEY,
                                              name VARCHAR(255) NOT NULL,
    description TEXT,
    price DECIMAL(10,2) NOT NULL,
    stock INTEGER NOT NULL,
    last_updated TIMESTAMP DEFAULT NOW()
    );
