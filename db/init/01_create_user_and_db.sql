DO $$
BEGIN
  IF NOT EXISTS (SELECT FROM pg_roles WHERE rolname = 'appuser') THEN
    CREATE USER appuser WITH PASSWORD 'appuserpass' CREATEDB CREATEROLE;
  END IF;
END $$;

DO $$
BEGIN
  IF NOT EXISTS (SELECT FROM pg_roles WHERE rolname = 'appuser' AND rolcreaterole = true) THEN
    ALTER USER appuser CREATEROLE;
  END IF;
END $$;

SELECT 'CREATE DATABASE postgres OWNER appuser'
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'postgres')\gexec

-- Grant privileges to appuser on the public schema
GRANT ALL PRIVILEGES ON DATABASE postgres TO appuser;
GRANT USAGE, CREATE ON SCHEMA public TO appuser;
