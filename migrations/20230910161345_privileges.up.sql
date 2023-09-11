CREATE TABLE IF NOT EXISTS privileges (
    id SERIAL PRIMARY KEY,
    privilege_title VARCHAR(20) NOT NULL UNIQUE,
    created_at TIMESTAMP NOT NULL
);