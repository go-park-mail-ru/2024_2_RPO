FROM alpine

WORKDIR /app
COPY ./build/board_app .
RUN chmod +x ./board_app

ENTRYPOINT ./board_app
