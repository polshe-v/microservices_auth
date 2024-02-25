#!/bin/bash
source .env

export MIGRATION_DSN="host=$POSTGRES_HOST port=5432 dbname=$POSTGRES_DB user=$POSTGRES_USER password=$POSTGRES_PASSWORD sslmode=disable"

goose -dir "${MIGRATION_DIR}" postgres "${MIGRATION_DSN}" up -v
