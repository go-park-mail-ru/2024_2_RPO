## Trello

### by team RPO

[Стремин Валентин](https://github.com/supchaser)
[Жуков Георгий](https://github.com/dedxyk594)
[Константин Сафронов](https://github.com/kosafronov)

Ментор: [Тарасов Артём](https://github.com/tarasovxx)

[Репозиторий фронтенда](https://github.com/frontend-park-mail-ru/2024_2_RPO)

[Ссылка на деплой](http://109.120.180.70:8002)

[Ссылка на Swagger и `CREATE TABLE`](https://github.com/go-park-mail-ru/2024_2_RPO/tree/swagger_approved)

### Стандарты разработки

* Все комментарии на русском языке
* Все логи на английском языке

### Запуск сервера

Надо развернуть PostgreSQL и Redis

Создать в PostgresQL базу данных pumpkin

Запустить `CREATE TABLE`, который лежит в ветке `swagger_approved`

Затем надо оформить файл `.env`. Пример:

```
DB_URL = postgres://tarasovxx:my_secure_password@localhost:5432/pumpkin

SERVER_PORT = 8800

REDIS_URL = redis://:my_secure_password@localhost:6379

CORS_ORIGIN = https://mysite.com

LOGS_FILE = tarasovxx.json
```

Запуск: `make run`

### Запуск тестов

Надо сделать всё, что нужно для запуска сервера, но вместо `make run` запустить `make test`. Чтобы получить информацию о покрытии, надо запустить `make coverage`
