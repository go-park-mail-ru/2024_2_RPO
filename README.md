# Pumpkin (Trello clone)

## by team RPO

- [Стремин Валентин](https://github.com/supchaser)

- [Жуков Георгий](https://github.com/dedxyk594)

- [Константин Сафронов](https://github.com/kosafronov)

Менторы:

- [Тарасов Артём](https://github.com/tarasovxx) - Backend

- [Фикслер Леонид](https://github.com/reddiridabl666) - Backend

- [Алёхин Владислав](https://github.com/3kybika) - СУБД

[Репозиторий фронтенда](https://github.com/frontend-park-mail-ru/2024_2_RPO)

[Ссылка на деплой](https://kanban-pumpkin.ru)

[Ссылка на Swagger](https://dedxyk594.github.io/swagger_ui_RPO/index.html)

## Стандарты разработки

- Все комментарии на русском языке
- Все логи на английском языке
- У аббревиатур все буквы в одном регистре: ~~`sessionId`~~ `sessionID`
- В импортах между нашими пакетами и стандартными должна быть пустая строка

## Запуск сервера

### Docker

Этот вариант подходит для Production, поскольку его легко разворачивать.

Для docker-compose файла, который лежит в репе, надо задать следующие настройки в .env:

```
POSTGRES_HOST = postgres
POSTGRES_PORT = 5432

REDIS_HOST = redis
REDIS_PORT = 6379

SERVER_PORT = 8800

USER_UPLOADS_DIR = /uploads
```

### Локальный

> [!CAUTION]
> Этот раздел может значительно измениться при внедрении микросервисов

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

TEST_DATABASE_URL = postgresql://3kybika:12345678@localhost:5432/migrate_gen_db?sslmode=disable
```

Запуск: `make run`

Миграции: `make migrate-up`

### Установка зависимостей

go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

### Запуск тестов

Надо сделать всё, что нужно для запуска сервера, но вместо `make run` запустить `make test`. Чтобы получить информацию о покрытии, надо запустить `make coverage`

## Схема базы данных

Актуальная модель данных находится в директории `database/schema.sql`

Для создания миграций используется программа [Atlas](https://atlasgo.io). Для применения используется [go-migrate](https://github.com/golang-migrate/migrate)

После изменения модели данных надо сгенерировать миграцию. Чтобы её сгенерировать, понадобится dev-база (и развёрнутый Postgres), потому что [без него Atlas не создаст миграции](https://atlasgo.io/atlas-schema/sql#dev-database).

URL dev-базы надо указать в файлe `.env` по имени `TEST_DATABASE_URL`. Эта база должна быть пустая; после работы atlas очистит всё, что он насоздавал. У пользователя должны быть все права на базу, а также права на создание ролей

Миграцию надо генерировать командой `make make-migrations`. Имя миграции не должно содержать пробелы, должно быть типа `add_tags_table`

Применять миграции надо командой `make migrate-up`. Команда интерактивно попросит логин и пароль от root-пользователя. На проде надо делать миграции пользователем postgres
