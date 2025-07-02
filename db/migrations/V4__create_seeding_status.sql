CREATE TABLE IF NOT EXISTS seeding_status (
    id SERIAL PRIMARY KEY,
    migration_version INTEGER NOT NULL DEFAULT 0,
    table_name VARCHAR(255),
    success BOOLEAN NOT NULL DEFAULT FALSE,
    seeded_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Main tracking record (with null table_name)
INSERT INTO seeding_status (migration_version, success, table_name)
VALUES (0, true, NULL)
ON CONFLICT DO NOTHING;