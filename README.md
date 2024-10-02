## Trello

### by team RPO

[Стремин Валентин](https://github.com/supchaser)
[Жуков Георгий](https://github.com/dedxyk594)
[Константин Сафронов](https://github.com/kosafronov)

Ментор: [Тарасов Артём](https://github.com/tarasovxx)

[Репозиторий фронтенда](https://github.com/frontend-park-mail-ru/2024_2_RPO)

[Ссылка на деплой](http://109.120.180.70:8002)

[Ссылка на Swagger и `CREATE TABLE`](https://github.com/go-park-mail-ru/2024_2_RPO/tree/swagger_approved)

### Запуск сервера

Надо развернуть PostgreSQL и Redis

Создать в PostgresQL базу данных pumpkin

Запустить `CREATE TABLE`, который лежит в ветке `swagger_approved`

Затем надо оформить файл `.env`. Пример:

```
DB_PASSWORD=my_secure_password
DB_USER=tarasovxx
DB_PORT=5432

SERVER_PORT=8800

REDIS_PORT=6379
REDIS_USER=tarasovxx
REDIS_PASSWORD=my_secure_password
```

Запуск: `go run main.go`

### Запуск тестов

Надо сделать всё, что нужно для запуска сервера, но вместо `go run main.go` запустить `./get_coverage_data.sh`
