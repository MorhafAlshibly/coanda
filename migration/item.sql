CREATE TABLE item (
    id VARCHAR(255) NOT NULL,
    type VARCHAR(255) NOT NULL,
    data JSON NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    expires_at DATETIME NULL DEFAULT NULL,
    PRIMARY KEY (id, type)
) ENGINE = InnoDB;
-- CREATE TRIGGER item_cleanup
-- AFTER
-- UPDATE ON item FOR EACH ROW BEGIN DELETE item
-- WHERE expires_at < NOW()
--     AND expires_at IS NOT NULL
-- END;