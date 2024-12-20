FROM ubuntu:24.10
RUN apt update
RUN apt install wget make -y

WORKDIR /tmp
RUN wget https://github.com/golang-migrate/migrate/releases/download/v4.18.1/migrate.linux-amd64.tar.gz -O migrate.tar.gz
RUN tar -xzf migrate.tar.gz
RUN cp ./migrate /usr/bin/migrate
RUN chmod +x /usr/bin/migrate
RUN rm -rf migrate migrate.tar.gz

WORKDIR /app
COPY ./Makefile .
COPY ./database/elastic/mapping.json .
COPY ./database/migrations/ ./migrations/
COPY ./build/auth_app ./build/utils ./
RUN chmod +x ./auth_app && chmod +x utils

CMD ["/app/auth_app"]
