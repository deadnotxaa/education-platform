DO $$
BEGIN
  CREATE ROLE analytic NOLOGIN;
EXCEPTION WHEN duplicate_object THEN
  RAISE NOTICE 'Role analytic already exists!';
END
$$;

-- Grant USAGE on the public schema
GRANT USAGE ON SCHEMA public TO analytic;

-- Grant SELECT on all existing tables
GRANT SELECT ON ALL TABLES IN SCHEMA public TO analytic;

-- Revoke write permissions explicitly
REVOKE INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public FROM analytic;

-- Set default privileges for future tables
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT SELECT ON TABLES TO analytic;
ALTER DEFAULT PRIVILEGES IN SCHEMA public REVOKE INSERT, UPDATE, DELETE ON TABLES FROM analytic;
