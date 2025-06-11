#!/bin/sh

# Exit on error
set -e

# Variables
BACKUP_DIR=/backups
TIMESTAMP=$(date +"%Y%m%d_%H%M%S")
BACKUP_FILE="${BACKUP_DIR}/mydatabase_backup_${TIMESTAMP}.sql"
RETENTION_COUNT=${BACKUP_RETENTION_COUNT:-5}

# Ensure backup directory exists
mkdir -p ${BACKUP_DIR}

# Perform backup
pg_dump -h ${PGHOST} -p ${PGPORT} -U ${PGUSER} -d ${PGDB} -F p -f ${BACKUP_FILE}

# Retain only the last RETENTION_COUNT backups
ls -t ${BACKUP_DIR}/*.sql | tail -n +$((RETENTION_COUNT + 1)) | xargs --no-run-if-empty rm -f

echo "Backup completed: ${BACKUP_FILE}"