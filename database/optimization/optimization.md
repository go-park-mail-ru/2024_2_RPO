# Домашнее задание №3 по СУБД

## Заполнение базы данных

Заполняться будет одна доска. Будет создано 100 колонок, для каждой из них - 100 карточек.
В каждой второй карточке будет 10 вложений, 10 комментариев и 10 строк чеклиста.

Была создана доска с board_id=209 и тестовый пользователь с u_id=39

## Получение содержимого доски

### Откроем большую доску

Оно не вернуло список карточек. Это из-за того, что выражение отменено по таймауту. Почему не было 500-ки? Потому что pgx вернул ошибку `pgx.ErrNoRows`

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

### Перепишем без коррелирующих подзапросов

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

Теперь оно загружает доску за 10 секунд.

Теперь попробуем сделать stress-тестирование. Для этого был создан файл `run_vegeta.sh`:

```
=====
Start stress test
URL: https://kanban-pumpkin.ru/api/cards/board_209/allContent
Method: GET
Duration: 60s
Max workers: 2
=====
=====
Test finished! Creating report...
Requests      [total, rate, throughput]         13, 0.18, 0.17
Duration      [total, attack, wait]             1m18s, 1m12s, 5.755s
Latencies     [min, mean, 50, 90, 95, 99, max]  5.755s, 11.541s, 10.515s, 20.115s, 20.145s, 20.152s, 20.152s
Bytes In      [total, mean]                     34086823, 2622063.31
Bytes Out     [total, mean]                     0, 0.00
Success       [ratio]                           100.00%
Status Codes  [code:count]                      200:13
Error Set:
```

Мы видим, что запрос делается катастрофически долго. Пришлось даже ограничить количество воркеров, чтобы запросы не заканчивались по таймауту.

### Рассмотрим EXPLAIN и разработаем стратегию оптимизации

![](./screens/1.svg)
