CREATE TABLE IF NOT EXISTS privileges (
    id SERIAL PRIMARY KEY,
    privilege_title VARCHAR(20) NOT NULL,
    created_at TIMESTAMP NOT NULL
);