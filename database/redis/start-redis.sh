chmod 777 /tmp/redis/redis.sock/
rm /tmp/redis/redis.sock/*

# Define cleanup function
terminate() {
    echo "==========   Stopping Redis server...   =========="
    redis-cli SHUTDOWN SAVE
    rm /tmp/redis/redis.sock/*
    exit 0
}

trap terminate TERM
trap terminate INT

echo "==========   Starting Redis server...   =========="

# Start redis
redis-server /pumpkin/redis/redis.conf &

REDIS_PID=$!
echo "==========   Redis PID: $REDIS_PID  =========="

wait $REDIS_PID

EXIT_CODE=$?

echo "==========   Redis exit code: $EXIT_CODE  =========="
