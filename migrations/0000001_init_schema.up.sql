CREATE TABLE user (
                        id UUID PRIMARY KEY,
                        username VARCHAR(255) NOT NULL UNIQUE,
                        email VARCHAR(255) NOT NULL UNIQUE,
                        password_hash VARCHAR(255) NOT NULL,
                        is_admin BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE TABLE book (
                      id UUID PRIMARY KEY,
                      genre VARCHAR(255) NOT NULL,
                      name VARCHAR(255) NOT NULL,
                      author VARCHAR(255) NOT NULL,
                      year INTEGER NOT NULL CHECK (year >= 0)
    );

CREATE TABLE user_book (
                           user_id UUID REFERENCES "user"(id),
                           book_id UUID REFERENCES book(id),
                           PRIMARY KEY (user_id, book_id)
);

CREATE INDEX idx_book_author ON book(author);
CREATE INDEX idx_book_year ON book(year);
CREATE INDEX idx_book_genre ON book(genre);
CREATE INDEX idx_user_book_user ON user_book(user_id);
CREATE INDEX idx_user_book_book ON user_book(book_id);