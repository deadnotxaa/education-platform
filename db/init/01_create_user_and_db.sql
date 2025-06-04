-- Create application user with CREATEDB and CREATEROLE privileges
DO $$
BEGIN
  IF NOT EXISTS (SELECT FROM pg_roles WHERE rolname = 'appuser') THEN
    CREATE USER appuser WITH PASSWORD 'appuserpass' CREATEDB CREATEROLE;
  END IF;
END $$;

-- Grant CREATEROLE if not already granted
DO $$
BEGIN
  IF NOT EXISTS (SELECT FROM pg_roles WHERE rolname = 'appuser' AND rolcreaterole = true) THEN
    ALTER USER appuser CREATEROLE;
  END IF;
END $$;

-- Create database if it doesn't exist
SELECT 'CREATE DATABASE mydatabase OWNER appuser'
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'mydatabase')\gexec

-- Grant privileges to appuser
GRANT ALL PRIVILEGES ON DATABASE mydatabase TO appuser;
