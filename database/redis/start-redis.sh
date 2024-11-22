#!/usr/bin/env bash

# Define cleanup function
stop_redis() {
    echo "==========   Stopping Redis server...   =========="
    redis-cli SHUTDOWN SAVE
    exit 0
}

# Trap SIGINT and SIGTERM signals and run cleanup function
trap stop_redis SIGINT
trap stop_redis SIGTERM

# Start redis
redis-server /pumpkin/redis/redis.conf &
wait %?redis-server
