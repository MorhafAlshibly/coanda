CREATE TABLE item (
    id VARCHAR(255) NOT NULL,
    data JSON,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP,
    PRIMARY KEY (id)
);

CREATE INDEX idx_expires_at ON item (expires_at);