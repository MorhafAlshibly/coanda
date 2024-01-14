CREATE TABLE record (
    name VARCHAR(255) NOT NULL,
    user_id BIGINT UNSIGNED NOT NULL,
    record BIGINT UNSIGNED NOT NULL,
    data JSON,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (name, user_id)
);
CREATE INDEX idx_name_record ON record (name, record);
CREATE INDEX idx_user_id_record ON record (user_id, record);
CREATE INDEX idx_record ON record (record);
CREATE VIEW ranked_record AS
SELECT name,
    user_id,
    record,
    data,
    DENSE_RANK() OVER (
        PARTITION BY name
        ORDER BY record ASC
    ) AS ranking,
    created_at,
    updated_at
FROM record;