FROM alpine

WORKDIR /app
COPY ./build/poll_app .
RUN chmod +x ./poll_app

ENTRYPOINT ./poll_app
