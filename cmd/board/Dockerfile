FROM ubuntu

WORKDIR /app
COPY ./build/board_app .
RUN chmod +x ./board_app

ENTRYPOINT /app/board_app
