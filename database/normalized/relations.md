# Домашнее задание №1 по курсу СУБД

## Описание отношений

Отношение "Пользователь"

`{u_id} -> {nickname, description, joined_at, updated_at, password_hash, email, avatar_file_uuid}`

Отношение "Загруженный пользователем файл"

`{file_uuid} -> {file_extension, size, created_at, created_by, type}`

Отношение "Доска"

`{board_id} -> {name, description, created_at, created_by, background_image_uuid, updated_at}`

Отношение "Пользователь на доске"

`{user_to_board_id} -> {u_id, board_id, added_at, updated_at, last_visit_at, added_by, updated_by, can_edit, can_share, can_invite_members, is_admin, notification_level}`

Отношение "Колонка канбана"

`{col_id} -> {board_id, title, created_at, updated_at, order_index}`

Отношение "Карточка"

`{card_id} -> {col_id, title, order_index, created_at, updated_at, cover_file_uuid}`

Отношение "Тег"

`{tag_id} -> {board_id, title, created_at, updated_at}`

Отношение "Тег к карточке"

`{tag_to_card_id} -> {card_id, tag_id, created_at}`

Отношение "Обновление карточки"

`{card_update_id} -> {card_id, is_visible, created_at, created_by, assigned_to, text, type, attached_file_uuid}`

Отношение "Уведомление"

`{notification_id} -> {u_id, board_id, notification_type, content, is_dismissed, created_at}`

Отношение "Поле чек-листа"

`{checklist_field_id} -> {card_id, order_index, title, is_done}`

```mermaidjs
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
        BIGINT user_to_board_id PK
        BIGINT u_id
        BIGINT board_id
        TIMESTAMPTZ added_at
        TIMESTAMPTZ updated_at
        TIMESTAMPTZ last_visit_at
        BIGINT added_by
        BIGINT updated_by
        BOOLEAN can_edit
        BOOLEAN can_share
        BOOLEAN can_invite_members
        BOOLEAN is_admin
        INT notification_level
    }

    KANBAN_COLUMN {
        BIGINT col_id PK
        BIGINT board_id
        TEXT title
        TIMESTAMPTZ created_at
        TIMESTAMPTZ updated_at
        INT order_index
    }

    CARD {
        BIGINT card_id PK
        BIGINT col_id
        TEXT title
        INT order_index
        TIMESTAMPTZ created_at
        TIMESTAMPTZ updated_at
        UUID cover_file_uuid
    }

    TAG {
        BIGINT tag_id PK
        BIGINT board_id
        TEXT title
        TIMESTAMPTZ created_at
        TIMESTAMPTZ updated_at
    }

    TAG_TO_CARD {
        BIGINT tag_to_card_id PK
        BIGINT card_id
        BIGINT tag_id
        TIMESTAMPTZ created_at
    }

    CARD_UPDATE {
        BIGINT card_update_id PK
        BIGINT card_id
        BOOLEAN is_visible
        TIMESTAMPTZ created_at
        BIGINT created_by
        BIGINT assigned_to
        TEXT type
        TEXT text
        UUID attached_file_uuid
    }

    NOTIFICATION {
        BIGINT notification_id PK
        BIGINT u_id
        BIGINT board_id
        TEXT notification_type
        TEXT content
        BOOLEAN is_dismissed
        TIMESTAMPTZ created_at
    }

    CHECKLIST_FIELD {
        BIGINT checklist_field_id PK
        BIGINT card_id
        INT order_index
        TEXT title
        BOOLEAN is_done
    }

    %% Связи

    USER ||--o{ USER_UPLOADED_FILE : "создал"
    USER ||--o{ BOARD : "создал"
    USER ||--o{ USER_TO_BOARD : "участвует в"
    USER_UPLOADED_FILE ||--o{ USER : "использован в"
    BOARD ||--o{ USER_TO_BOARD : "содержит"
    BOARD ||--o{ KANBAN_COLUMN : "содержит"
    BOARD ||--o{ TAG : "содержит"
    BOARD ||--o{ NOTIFICATION : "относятся к"
    USER_UPLOADED_FILE ||--o{ BOARD : "фон"
    USER_UPLOADED_FILE ||--o{ CARD : "обложка"
    KANBAN_COLUMN ||--o{ CARD : "содержит"
    CARD ||--o{ TAG_TO_CARD : "имеет"
    TAG ||--o{ TAG_TO_CARD : "присоединен к"
    CARD ||--o{ CARD_UPDATE : "обновления"
    CARD ||--o{ CHECKLIST_FIELD : "содержит"
    USER ||--o{ CARD_UPDATE : "создал"
    USER ||--o{ NOTIFICATION : "получает"
```
