FROM ubuntu

WORKDIR /app
COPY ./build/user_app .
RUN chmod +x ./user_app

ENTRYPOINT /app/user_app
