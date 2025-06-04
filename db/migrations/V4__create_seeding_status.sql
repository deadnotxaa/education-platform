CREATE TABLE IF NOT EXISTS seeding_status (
    id SERIAL PRIMARY KEY,
    seeded BOOLEAN NOT NULL DEFAULT FALSE,
    seeded_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Insert a single row to initialize the status
INSERT INTO seeding_status (seeded) VALUES (FALSE) ON CONFLICT DO NOTHING;
