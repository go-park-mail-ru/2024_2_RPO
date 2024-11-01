# Домашнее задание №1 по курсу СУБД

## Описание отношений

#### Отношение "Пользователь"

`{u_id} -> {nickname, description, joined_at, updated_at, password_hash, email, avatar_file_uuid}`

`{email} -> {u_id, nickname, description, joined_at, updated_at, password_hash, avatar_file_uuid}`

`{nickname} -> {u_id, email, description, joined_at, updated_at, password_hash, avatar_file_uuid}`

`{nickname, email} -> {u_id, description, joined_at, updated_at, password_hash, avatar_file_uuid}`

`{nickname, u_id} -> {email, description, joined_at, updated_at, password_hash, avatar_file_uuid}`

`{email, u_id} -> {nickname, description, joined_at, updated_at, password_hash, avatar_file_uuid}`

`{email, u_id, nickname} -> {description, joined_at, updated_at, password_hash, avatar_file_uuid}`

#### Отношение "Загруженный пользователем файл"

`{file_uuid} -> {file_extension, size, created_at, created_by, type}`

#### Отношение "Доска"

`{board_id} -> {name, description, created_at, created_by, background_image_uuid, updated_at}`

#### Отношение "Пользователь на доске"

`{u_id} -> {board_id, added_at, updated_at, last_visit_at, added_by, updated_by, role}`

`{board_id} -> {u_id, added_at, updated_at, last_visit_at, added_by, updated_by, role}`

`{u_id, board_id} -> {added_at, updated_at, last_visit_at, added_by, updated_by, role}`

#### Отношение "Колонка канбана"

`{col_id} -> {board_id, title, created_at, updated_at, order_index}`

#### Отношение "Карточка"

`{card_id} -> {col_id, title, order_index, created_at, updated_at, cover_file_uuid}`

#### Транзитивные отношения (справедливые не для любых кортежей)

`{user.avatar_file_uuid} -> {user_uploaded_file.type}` (только для некоторых файлов)

`{board.background_file_uuid} -> {user_uploaded_file.type}` (только для некоторых файлов)

`{user_to_board.added_by, user_to_board.u_id, user_to_board.created_by -> {board.created_by}` (только для кортежа user_to_board, где пользователь сам создал доску и на ней присутствует)

## Нормализация модели

### 1 Н.Ф.

1. _Нет упорядочивания строк сверху вниз (другими словами, порядок строк не несет в себе никакой информации)._

Везде, где порядок имеет значение, мы порядковый индекс использовали как атрибут кортежа. Наша модель данных не зависит от порядка кортежей.

2. _Нет упорядочивания столбцов слева направо (другими словами, порядок столбцов не несет в себе никакой информации)._

У нас каждый атрибут имеет уникальное имя, по которому мы его достаём. Мы не ориентируемся на порядок столбцов.

3. _Нет повторяющихся строк._

Все отношения имеют первичный ключ, который не может повторяться между строк, поэтому в любых двух кортежах одного отношения будет отличаться хотя бы первичный ключ.

4. _Каждое пересечение строки и столбца содержит ровно одно значение из соответствующего домена (и больше ничего)._

У нас нет массивов, поэтому это требование соблюдается.

5. _Все столбцы являются обычными_

Это свойство Postgres'а как реляционной СУБД. А, например, для Cassandra это утверждение не было бы справедливо

### 2 Н.Ф.

_Переменная отношения находится во второй нормальной форме тогда и только тогда, когда она находится в первой нормальной форме и каждый неключевой атрибут неприводимо зависит от (каждого) её потенциального ключа_

Наша модель находится в 2НФ, потому что каждый неключевой атрибут зависит только от первичного ключа

### 3 Н.Ф.

_Переменная отношения R находится в третьей нормальной форме тогда и только тогда, когда неключевые атрибуты непосредственно (нетранзитивно) функционально зависят от ключей_

Наша модель находится в третьей нормальной форме, потому что транзитивные зависимости справедливы не для всех кортежей, и каждый неключевой атрибут непосредственно зависит от первичного ключа

### НФБК

_Переменная отношения находится в НФБК тогда и только тогда, когда каждый её детерминант является потенциальным ключом_

Наша модель находится в НФБК, потому что во всех вышеприведённых отношениях атрибуты, находящиеся в левой части, позволяют однозначно определить кортеж в таблице (то есть, являются потенциальным ключом)

## ERD-диаграмма

```mermaid
erDiagram
        USER {
        BIGINT u_id PK
        TEXT nickname
        TEXT description
        TIMESTAMPTZ joined_at
        TIMESTAMPTZ updated_at
        TEXT password_hash
        CITEXT email
        UUID avatar_file_uuid
    }

    USER_UPLOADED_FILE {
        UUID file_uuid PK
        TEXT file_extension
        INTEGER size
        TIMESTAMPTZ created_at
        BIGINT created_by
        FILE_TYPE type
    }

    BOARD {
        BIGINT board_id PK
        TEXT name
        TEXT description
        TIMESTAMPTZ created_at
        BIGINT created_by
        UUID background_image_uuid
        TIMESTAMPTZ updated_at
    }

    USER_TO_BOARD {
        BIGINT u_id PK
        BIGINT board_id PK
        TIMESTAMPTZ added_at
        TIMESTAMPTZ updated_at
        TIMESTAMPTZ last_visit_at
        BIGINT added_by FK
        BIGINT updated_by FK
        USER_ROLE role
    }

    KANBAN_COLUMN {
        BIGINT col_id PK
        BIGINT board_id FK
        TEXT title
        TIMESTAMPTZ created_at
        TIMESTAMPTZ updated_at
        INT order_index
    }

    CARD {
        BIGINT card_id PK
        BIGINT col_id FK
        TEXT title
        INT order_index
        TIMESTAMPTZ created_at
        TIMESTAMPTZ updated_at
        UUID cover_file_uuid FK
    }

    %% Типы ENUM
    FILE_TYPE_ENUM {
        FILE_TYPE avatar
        FILE_TYPE background_image
        FILE_TYPE card_cover_image
        FILE_TYPE pinned_file
    }

    USER_ROLE_ENUM {
        USER_ROLE viewer
        USER_ROLE editor
        USER_ROLE editor_chief
        USER_ROLE admin
    }

    %% Связи
    USER ||--o{ USER_UPLOADED_FILE : "created files"
    USER_UPLOADED_FILE ||--|| USER : "avatar_file_uuid"
    USER ||--o{ BOARD : "created boards"
    BOARD ||--o{ USER_UPLOADED_FILE : "background_image"
    USER ||--o{ USER_TO_BOARD : "boards membership"
    BOARD ||--o{ USER_TO_BOARD : "members"
    BOARD ||--o{ KANBAN_COLUMN : "columns"
    KANBAN_COLUMN ||--o{ CARD : "cards"
    USER_TO_BOARD }|..|{ USER : "added by"
    USER_TO_BOARD }|..|{ USER : "updated by"
    CARD ||--|| USER_UPLOADED_FILE : "cover file"
```
