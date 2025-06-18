-- +goose Up
CREATE TABLE USERS (
    id UUID primary key,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    email VARCHAR(255) NOT NULL UNIQUE
);


-- +goose Down                  
DROP TABLE USERS;