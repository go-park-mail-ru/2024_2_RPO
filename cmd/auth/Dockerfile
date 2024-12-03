FROM ubuntu
RUN apt update
RUN apt install wget make -y

WORKDIR /tmp
RUN wget https://github.com/golang-migrate/migrate/releases/download/v4.18.1/migrate.linux-amd64.tar.gz -O migrate.tar.gz
RUN tar -xzf migrate.tar.gz
RUN cp ./migrate /usr/bin/migrate
RUN rm -rf migrate migrate.tar.gz

WORKDIR /app
COPY ./Makefile .
COPY ./database/migrations/ ./migrations/
COPY ./build/auth_app .
RUN chmod +x ./auth_app
RUN chmod +x /usr/bin/migrate

ENTRYPOINT /app/auth_app
