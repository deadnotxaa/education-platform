#!/bin/sh

# Exit on error
set -e

# Check if ANALYST_NAMES is set
if [ -z "$ANALYST_NAMES" ]; then
  echo "ANALYST_NAMES environment variable is not set or empty, no users created"
  exit 0
fi

# Wait for PostgreSQL to be ready
until pg_isready -h "$PGHOST" -p "$PGPORT" -U "$PGUSER"; do
  echo "Waiting for PostgreSQL to be ready..."
  sleep 1
done

# Split ANALYST_NAMES by comma
IFS=','

for name in $ANALYST_NAMES; do
  # Trim whitespace
  name=$(echo "$name" | tr -d '[:space:]')
  
  # Skip empty names
  if [ -z "$name" ]; then
    continue
  fi

  echo "> Creating user $name"
  psql -h "$PGHOST" -p "$PGPORT" -U "$PGUSER" -d "$PGDB" -v ON_ERROR_STOP=1 <<EOF
DO \$\$ BEGIN
  CREATE USER "$name" WITH PASSWORD '${name}_123';
EXCEPTION WHEN duplicate_object THEN
  RAISE NOTICE 'User % already exists, skipping', '$name';
END \$\$;

GRANT analytic TO "$name";
EOF
done
