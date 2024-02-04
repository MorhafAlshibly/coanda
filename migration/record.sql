CREATE TABLE record (
    name VARCHAR(255) NOT NULL,
    user_id BIGINT UNSIGNED NOT NULL,
    record BIGINT UNSIGNED NOT NULL,
    data JSON NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (name, user_id),
    INDEX name_record_idx (name ASC, record ASC)
) ENGINE = InnoDB;
CREATE VIEW ranked_record AS
SELECT name,
    user_id,
    record,
    DENSE_RANK() OVER (
        PARTITION BY name
        ORDER BY record ASC
    ) AS ranking,
    data,
    created_at,
    updated_at
FROM record
ORDER BY record ASC;