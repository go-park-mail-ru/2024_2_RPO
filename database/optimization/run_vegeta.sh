URL=https://kanban-pumpkin.ru/api/cards/board_193/allContent
METHOD=GET
DURATION=60s
CSRF_TOKEN=12345
SESSION_ID=ed0626543ef422de1a5296e8f97b8750d5b6ac43f70384ea726ac0eb5b7b23e4
MAX_WORKERS=16

echo =====
echo Start stress test
echo URL: $URL
echo Method: $METHOD
echo Duration: $DURATION
echo Max workers: $MAX_WORKERS
echo =====

echo $METHOD $URL | vegeta attack -duration=$DURATION \
    -header "X-Csrf-Token: $CSRF_TOKEN" \
    -header "Cookie: session_id=$(echo $SESSION_ID); csrf_token=$(echo $CSRF_TOKEN)" \
    -max-workers $MAX_WORKERS \
    -timeout $DURATION \
    >results.vegeta.bin

echo =====
echo Test finished! Creating report...

vegeta report results.vegeta.bin
