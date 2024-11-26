FROM alpine

WORKDIR /app
COPY ./build/user_app .
RUN chmod +x ./user_app

ENTRYPOINT ./user_app
