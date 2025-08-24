CREATE TABLE monitored_urls (
    id SERIAL PRIMARY KEY,
    url TEXT NOT NULL UNIQUE,
    check_interval_sec INT NOT NULL CHECK (check_interval_sec BETWEEN 5 AND 300),
    regex_pattern TEXT
);