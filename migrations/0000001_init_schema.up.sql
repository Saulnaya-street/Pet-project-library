CREATE TABLE users (
                       id UUID PRIMARY KEY,
                       username VARCHAR(255) NOT NULL UNIQUE,
                       email VARCHAR(255) NOT NULL UNIQUE,
                       password_hash VARCHAR(255) NOT NULL,
                       is_admin BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE TABLE books (
                       id UUID PRIMARY KEY,
                       genre VARCHAR(255) NOT NULL,
                       name VARCHAR(255) NOT NULL,
                       author VARCHAR(255) NOT NULL,
                       year SERIAL
);

CREATE TABLE user_book (
                           user_id UUID REFERENCES users(id),
                           book_id UUID REFERENCES books(id),
                           PRIMARY KEY (user_id, book_id)
);

CREATE INDEX idx_book_author ON books(author);
CREATE INDEX idx_book_year ON books(year);
CREATE INDEX idx_book_genre ON books(genre);