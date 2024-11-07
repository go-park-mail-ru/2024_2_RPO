## Trello

### by team RPO

[Стремин Валентин](https://github.com/supchaser)
[Жуков Георгий](https://github.com/dedxyk594)
[Константин Сафронов](https://github.com/kosafronov)

Менторы:

- [Тарасов Артём](https://github.com/tarasovxx)
- [Фикслер Леонид](https://github.com/reddiridabl666)
- [Алёхин Владислав](https://github.com/3kybika)

[Репозиторий фронтенда](https://github.com/frontend-park-mail-ru/2024_2_RPO)

[Ссылка на деплой](http://109.120.180.70:8002)

[Ссылка на Swagger](https://dedxyk594.github.io/swagger_ui_RPO/index.html)

### Стандарты разработки

- Все комментарии на русском языке
- Все логи на английском языке
- У аббревиатур все буквы в одном регистре: ~~`sessionId`~~ `sessionID`

### Запуск сервера

Надо развернуть PostgreSQL и Redis

Создать в PostgresQL базу данных pumpkin

Запустить `CREATE TABLE`, который лежит в ветке `swagger_approved`

Затем надо оформить файл `.env`. Пример:

```
POSTGRES_HOST = localhost
POSTGRES_PORT = 5432
POSTGRES_USER = tarasovxx
POSTGRES_PASSWORD = my_secure_password
POSTGRES_DB = pumpkin
POSTGRES_SSLMODE = require

SERVER_PORT = 8800

REDIS_HOST = localhost
REDIS_PORT = 6379
REDIS_PASSWORD = my_secure_password

CORS_ORIGIN = https://mysite.com

LOGS_FILE = log.json

USER_UPLOADS_DIR = /opt/uploads
USER_UPLOADS_URL = /files
```

Запуск: `make run`

Миграции: `make migrate-up`

### Установка зависимостей

go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

### Запуск тестов

Надо сделать всё, что нужно для запуска сервера, но вместо `make run` запустить `make test`. Чтобы получить информацию о покрытии, надо запустить `make coverage`
