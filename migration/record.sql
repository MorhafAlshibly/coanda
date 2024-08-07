CREATE TABLE record (
    id BIGINT UNSIGNED PRIMARY KEY NOT NULL AUTO_INCREMENT,
    name VARCHAR(255) NOT NULL,
    user_id BIGINT UNSIGNED NOT NULL,
    record BIGINT UNSIGNED NOT NULL,
    data JSON NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE INDEX name_user_id_record_idx (name ASC, user_id ASC),
    INDEX name_record_idx (name ASC, record ASC),
    INDEX user_id_record_idx (user_id ASC, record ASC)
) ENGINE = InnoDB;
CREATE VIEW ranked_record AS
SELECT id,
    name,
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
ORDER BY name ASC,
    record ASC;