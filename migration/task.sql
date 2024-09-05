CREATE TABLE task (
    id VARCHAR(255) NOT NULL,
    type VARCHAR(255) NOT NULL,
    data JSON NOT NULL,
    expires_at DATETIME NULL,
    completed_at DATETIME NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id, type)
) ENGINE = InnoDB;