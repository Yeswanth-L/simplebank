#!/bin/sh

set -e  # Exit immediately if a command fails

echo "Running DB migration"

# Use `.` instead of `source` for POSIX compliance
. /app/app.env

# Ensure DB_SOURCE is not empty
if [ -z "$DB_SOURCE" ]; then
    echo "Error: DB_SOURCE is empty. Check your environment variables."
    exit 1
fi

# Run database migration
/app/migrate -path /app/db/migration -database "$DB_SOURCE" -verbose up

echo "Starting the app"
exec "$@"
