CREATE EXTENSION "uuid-ossp";

CREATE TABLE user_uploaded_file(
    file_id        BIGINT      PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    file_uuid      UUID        NOT NULL DEFAULT uuid_generate_v4(),
    file_hash      TEXT,
    file_extension TEXT,
    "size"         INTEGER     NOT NULL,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE "user" (
    u_id           BIGINT      GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    nickname       TEXT        NOT NULL UNIQUE,
    joined_at      TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    password_hash  TEXT,
    csat_poll_dt   TIMESTAMPTZ NOT NULL,
    email          TEXT UNIQUE NOT NULL,
    avatar_file_id BIGINT,

    FOREIGN KEY (avatar_file_id) REFERENCES user_uploaded_file(file_id) ON UPDATE CASCADE ON DELETE SET NULL
);

CREATE TYPE user_role AS ENUM (
    'viewer',
    'editor',
    'editor_chief',
    'admin'
);

CREATE TABLE board (
    board_id            BIGINT      GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    "name"              TEXT        NOT NULL,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by          BIGINT,
    background_image_id BIGINT,
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (created_by) REFERENCES "user"(u_id) ON UPDATE CASCADE ON DELETE SET NULL,
    FOREIGN KEY (background_image_id) REFERENCES user_uploaded_file(file_id) ON UPDATE CASCADE ON DELETE SET NULL
);

CREATE TABLE user_to_board (
    u_id             BIGINT      NOT NULL,
    board_id         BIGINT      NOT NULL,
    added_at         TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    last_visit_at    TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    added_by         BIGINT,
    updated_by       BIGINT,
    invite_link_uuid UUID,
    "role"           user_role   NOT NULL DEFAULT 'viewer',

    FOREIGN KEY (u_id) REFERENCES "user"(u_id) ON UPDATE CASCADE ON DELETE CASCADE,
    FOREIGN KEY (board_id) REFERENCES board(board_id) ON UPDATE CASCADE ON DELETE CASCADE,
    FOREIGN KEY (added_by) REFERENCES "user"(u_id) ON UPDATE CASCADE ON DELETE SET NULL,
    FOREIGN KEY (updated_by) REFERENCES "user"(u_id) ON UPDATE CASCADE ON DELETE SET NULL,
    PRIMARY KEY(u_id, board_id)
);

CREATE TABLE kanban_column (
    col_id      BIGINT      GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    board_id    BIGINT      NOT NULL,
    title       TEXT        NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    order_index INT         NOT NULL, -- Порядковый номер колонки на доске

    FOREIGN KEY (board_id) REFERENCES board(board_id) ON UPDATE CASCADE ON DELETE CASCADE,
    UNIQUE (board_id, order_index) DEFERRABLE INITIALLY DEFERRED
);

CREATE TABLE "card" (
    card_id       BIGINT      GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    card_uuid     UUID        NOT NULL DEFAULT uuid_generate_v4(), -- UUID для ссылки на карточку
    title         TEXT        NOT NULL,
    col_id        BIGINT      NOT NULL,
    order_index   INTEGER, -- Порядковый номер карточки в колонке
    created_at    TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    cover_file_id BIGINT,
    deadline      TIMESTAMPTZ,
    is_done       BOOLEAN NOT NULL DEFAULT FALSE, -- Видна, когда задан deadline или чеклист

    FOREIGN KEY (col_id) REFERENCES kanban_column(col_id) ON UPDATE CASCADE ON DELETE CASCADE,
    FOREIGN KEY (cover_file_id) REFERENCES user_uploaded_file(file_id) ON UPDATE CASCADE ON DELETE SET NULL
);

CREATE TABLE tag (
    tag_id     BIGINT      GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    title      TEXT        NOT NULL,
    board_id   BIGINT,
    color      TEXT        NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (board_id) REFERENCES board(board_id) ON UPDATE CASCADE ON DELETE CASCADE
);

CREATE TABLE tag_to_card (
    tag_id        BIGINT,
    card_id       BIGINT,

    FOREIGN KEY (card_id) REFERENCES card(card_id) ON UPDATE CASCADE ON DELETE CASCADE,
    FOREIGN KEY (tag_id) REFERENCES tag(tag_id) ON UPDATE CASCADE ON DELETE CASCADE,
    PRIMARY KEY(tag_id, card_id)
);

CREATE TABLE card_attachment (
    attachment_id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    card_id       BIGINT      NOT NULL,
    file_id       BIGINT      NOT NULL,
    original_name TEXT        NOT NULL,
    attached_at   TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    attached_by   BIGINT NOT  NULL,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (card_id) REFERENCES card(card_id) ON UPDATE CASCADE ON DELETE CASCADE,
    FOREIGN KEY (file_id) REFERENCES user_uploaded_file(file_id) ON UPDATE CASCADE ON DELETE CASCADE,
    FOREIGN KEY (attached_by) REFERENCES "user"(u_id) ON UPDATE CASCADE ON DELETE CASCADE
);

CREATE TABLE checklist_field (
    checklist_field_id BIGINT      GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    card_id            BIGINT      NOT NULL,
    title              TEXT        NOT NULL,
    is_done            BOOLEAN     NOT NULL DEFAULT FALSE,
    created_at         TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    order_index        INTEGER,

    FOREIGN KEY (card_id) REFERENCES card(card_id) ON UPDATE CASCADE ON DELETE CASCADE
);

CREATE TABLE card_comment (
    comment_id BIGINT      GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    card_id    BIGINT      NOT NULL,
    title      TEXT        NOT NULL,
    created_by BIGINT      NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    is_edited  BOOLEAN     NOT NULL DEFAULT FALSE,

    FOREIGN KEY (card_id) REFERENCES card(card_id) ON UPDATE CASCADE ON DELETE CASCADE,
    FOREIGN KEY (created_by) REFERENCES "user"(u_id) ON UPDATE CASCADE ON DELETE CASCADE
);

CREATE TABLE card_user_assignment (
    assignment_id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    card_id       BIGINT NOT NULL,
    u_id          BIGINT NOT NULL,

    FOREIGN KEY (card_id) REFERENCES card(card_id) ON UPDATE CASCADE ON DELETE CASCADE,
    FOREIGN KEY (u_id) REFERENCES "user"(u_id) ON UPDATE CASCADE ON DELETE CASCADE
);

CREATE TYPE question_type AS ENUM (
    'answer_text',
    'answer_rating'
);

CREATE TABLE csat_question (
    question_id   BIGINT        GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    is_active     BOOLEAN       NOT NULL DEFAULT TRUE,
    question_text TEXT          NOT NULL,
    "type"        question_type NOT NULL,
    created_at    TIMESTAMPTZ   NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE csat_results (
    result_id   BIGINT      GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    question_id BIGINT      NOT NULL,
    rating      INTEGER,
    comment     TEXT,
    u_id        BIGINT      NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (u_id) REFERENCES "user"(u_id) ON UPDATE CASCADE ON DELETE CASCADE,
    FOREIGN KEY (question_id) REFERENCES csat_question(question_id) ON UPDATE CASCADE ON DELETE CASCADE
);
