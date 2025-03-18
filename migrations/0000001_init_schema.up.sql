CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE "user" (
                        id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                        username VARCHAR(255) NOT NULL UNIQUE,
                        email VARCHAR(255) NOT NULL UNIQUE,
                        password_hash VARCHAR(255) NOT NULL,
                        is_admin BOOLEAN NOT NULL DEFAULT FALSE,
                        created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE genre (
                       id_genre UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                       name_genre VARCHAR(255) NOT NULL UNIQUE,
                       created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE book (
                      id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                      id_user UUID NOT NULL REFERENCES "user"(id),
                      id_genre UUID NOT NULL REFERENCES genre(id_genre),
                      name VARCHAR(255) NOT NULL,
                      author VARCHAR(255) NOT NULL,
                      year INTEGER NOT NULL,
                      created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                      updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_book_author ON book(author);
CREATE INDEX idx_book_year ON book(year);
CREATE INDEX idx_book_genre ON book(id_genre);