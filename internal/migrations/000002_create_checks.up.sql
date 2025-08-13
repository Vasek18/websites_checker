CREATE TABLE checks (
    id SERIAL PRIMARY KEY,
    url TEXT NOT NULL,
    check_timestamp TIMESTAMPTZ NOT NULL,
    response_time_ms INT,
    http_status INT,
    regex_match BOOLEAN,
    error TEXT
);