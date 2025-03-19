CREATE TABLE IF NOT EXISTS movies(
    movie_id uuid DEFAULT gen_random_uuid(),
    title VARCHAR NOT NULL,
    genre VARCHAR NOT NULL,
    director VARCHAR NOT NULL,
    year INTEGER NOT NULL CHECK (year > 0)
);