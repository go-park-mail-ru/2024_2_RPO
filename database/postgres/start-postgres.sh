#!/bin/bash

echo "Waiting 1 second for doxxr to mnt those volumz"
sleep 1

chmod 777 /pg_logs
chmod 777 /tmp/postgres/postgres.sock/
chown postgres /pg_data/
su postgres -c "chmod 700 /pg_data/"
rm /tmp/postgres/postgres.sock/*

chown postgres /pg_setup/server.key
chown postgres /pg_setup/server.crt
su postgres -c "chmod 600 /pg_setup/server.key"
su postgres -c "chmod 600 /pg_setup/server.crt"

if [ ! -f "/pg_data/PG_VERSION" ]; then
    BLANK=1
    echo "Creating PostgreSQL db..."

    su postgres -c "initdb -D /pg_data/"

    # Проверка успешности инициализации
    if [ $? -eq 0 ]; then
        echo "PostgreSQL db created"
    else
        echo "PostgreSQL db NOT CREATED, THERE IS AN ERROR" >&2
        exit 1
    fi
fi

terminate() {
    echo "Gracefully stop PostgreSQL server..."
    kill -TERM "$POSTGRES_PID"
    rm /tmp/postgres/postgres.sock/*
}

trap terminate SIGINT SIGTERM

echo "======> STARTING POSTGRESQL SERVER"

su postgres -c "postgres -c config_file=/pg_setup/postgresql.conf -c hba_file=/pg_setup/pg_hba.conf" &

POSTGRES_PID=$!

echo "PostgreSQL works with PID $POSTGRES_PID"

if [ $BLANK -eq 1 ]; then
    echo "BLANK, so create db pumpkin"
    echo "Sleep 60 seconds to wait server to create and start..."
    sleep 60
    su postgres -c "createdb -h /tmp/postgres/postgres.sock -U postgres pumpkin --locale en_US.utf8"
    if [ $? -eq 0 ]; then
        echo "Service db pumpkin created"
    else
        echo "Service db pumpkin NOT CREATED, THERE IS AN ERROR" >&2
        exit 1
    fi
fi

wait "$POSTGRES_PID"
EXIT_CODE=$?

echo "======> POSTGRESQL EXIT CODE $EXIT_CODE"
