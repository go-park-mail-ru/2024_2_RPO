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

```
