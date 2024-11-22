#!/bin/bash

chmod 777 /pg_logs
chmod 777 /tmp/postgres/postgres.sock/

chown postgres /pg_setup/server.key
chown postgres /pg_setup/server.crt
su postgres -c "chmod 600 /pg_setup/server.key"
su postgres -c "chmod 600 /pg_setup/server.crt"

terminate() {
    echo "Gracefully stop PostgreSQL server..."
    kill -TERM "$POSTGRES_PID"
}

trap terminate SIGINT SIGTERM

echo "======> STARTING POSTGRESQL SERVER"

su postgres -c "postgres -c config_file=/pg_setup/postgresql.conf -c hba_file=/pg_setup/pg_hba.conf" &

POSTGRES_PID=$!

echo "PostgreSQL works with PID $POSTGRES_PID"

wait "$POSTGRES_PID"
EXIT_CODE=$?

echo "======> POSTGRESQL EXIT CODE $EXIT_CODE"
