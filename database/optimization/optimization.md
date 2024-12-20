# Домашнее задание №3 по СУБД

## Заполнение базы данных

Заполняться будет одна доска. Будет создано 100 колонок, для каждой из них - 100 карточек.
В каждой второй карточке будет 10 вложений, 10 комментариев и 10 строк чеклиста.

## Попытка открыть доску с фронта

Вылетает ошибка 500. Это из-за того, что выражение отменено по таймауту.

Посмотрим на запрос, который был отменён:

```sql
SELECT
    c.card_id,
    c.col_id,
    c.title,
    c.created_at,
    c.updated_at,
    c.deadline,
    c.is_done,
    (SELECT (NOT COUNT(*)=0) FROM checklist_field AS f WHERE f.card_id=c.card_id),
    (SELECT (NOT COUNT(*)=0) FROM card_attachment AS f WHERE f.card_id=c.card_id),
    (SELECT (NOT COUNT(*)=0) FROM card_user_assignment AS f WHERE f.card_id=c.card_id),
    (SELECT (NOT COUNT(*)=0) FROM card_comment AS f WHERE f.card_id=c.card_id),
    COALESCE(uuf.file_uuid::text, ''),
    COALESCE(uuf.file_extension::text, ''),
    c.card_uuid
FROM card c
JOIN kanban_column kc ON c.col_id = kc.col_id
LEFT JOIN user_uploaded_file uuf ON c.cover_file_id = uuf.file_id
WHERE kc.board_id = $1
ORDER BY c.order_index;
```

Запрос содержит коррелирующие подзапросы. Скорее всего, именно они приводят к таймауту.

## Перепишем без коррелирующих подзапросов

```sql
SELECT
    c.card_id,
    c.col_id,
    c.title,
    c.created_at,
    c.updated_at,
    c.deadline,
    c.is_done,
    COUNT(clf.checklist_field_id)>0,
    COUNT(a.attachment_id)>0,
    COUNT(ua.assignment_id)>0,
    COUNT(com.comment_id)>0,
    COALESCE(uuf.file_uuid::text, ''),
    COALESCE(uuf.file_extension::text, ''),
    c.card_uuid
FROM card c
JOIN kanban_column kc ON c.col_id = kc.col_id
LEFT JOIN user_uploaded_file uuf ON c.cover_file_id = uuf.file_id
LEFT JOIN checklist_field clf ON clf.card_id=c.card_id
LEFT JOIN card_attachment a ON a.card_id=c.card_id
LEFT JOIN card_comment com ON com.card_id=c.card_id
LEFT JOIN card_user_assignment ua ON ua.card_id=c.card_id
WHERE kc.board_id = $1
GROUP BY c.card_id, uuf.file_id
ORDER BY c.order_index;
```
