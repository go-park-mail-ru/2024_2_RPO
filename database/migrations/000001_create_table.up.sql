CREATE EXTENSION IF NOT EXISTS "citext";
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE "user" (
    u_id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    nickname TEXT NOT NULL UNIQUE,
    "description" TEXT NOT NULL,
    joined_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    password_hash TEXT NOT NULL,
    email CITEXT UNIQUE NOT NULL
    -- avatar_file_uuid UUID,
    -- FOREIGN KEY (avatar_file_uuid) REFERENCES user_uploaded_file(file_uuid)
    -- ON UPDATE CASCADE ON DELETE SET NULL
);

CREATE TYPE FILE_TYPE AS ENUM
(
    'avatar',
    'background_image',
    'card_cover_image',
    'file'
);

CREATE TABLE user_uploaded_file(
    file_uuid UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    file_extension TEXT,
    "size" INTEGER NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by BIGINT,
    "type" FILE_TYPE NOT NULL,
    FOREIGN KEY (created_by) REFERENCES "user"(u_id) ON UPDATE CASCADE ON DELETE SET NULL
);

ALTER TABLE "user" ADD COLUMN avatar_file_uuid UUID;

ALTER TABLE "user" ADD CONSTRAINT fk_avatar_file_uuid FOREIGN KEY (avatar_file_uuid) REFERENCES user_uploaded_file(file_uuid)
    ON UPDATE CASCADE ON DELETE SET NULL;

CREATE TABLE board (
    board_id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by BIGINT,
    background_image_uuid UUID,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (created_by) REFERENCES "user"(u_id) ON UPDATE CASCADE ON DELETE
    SET NULL,
    FOREIGN KEY (background_image_uuid) REFERENCES user_uploaded_file(file_uuid)
    ON UPDATE CASCADE ON DELETE SET NULL
);

CREATE TABLE user_to_board (
    user_to_board_id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    u_id BIGINT NOT NULL,
    board_id BIGINT NOT NULL,
    added_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    last_visit_at TIMESTAMPTZ,
    added_by BIGINT,
    updated_by BIGINT,
    can_edit BOOLEAN DEFAULT FALSE,
    can_share BOOLEAN DEFAULT FALSE,
    can_invite_members BOOLEAN DEFAULT FALSE,
    is_admin BOOLEAN DEFAULT FALSE,
    notification_level INT,
    FOREIGN KEY (u_id) REFERENCES "user"(u_id) ON UPDATE CASCADE ON DELETE CASCADE,
    FOREIGN KEY (board_id) REFERENCES board(board_id) ON UPDATE CASCADE ON DELETE CASCADE,
    FOREIGN KEY (added_by) REFERENCES "user"(u_id) ON UPDATE CASCADE ON DELETE
    SET NULL,
        FOREIGN KEY (updated_by) REFERENCES "user"(u_id) ON UPDATE CASCADE ON DELETE
    SET NULL
);

CREATE TABLE kanban_column (
    col_id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    board_id BIGINT NOT NULL,
    title TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    order_index INT, -- Порядковый номер колонки на доске
    FOREIGN KEY (board_id) REFERENCES board(board_id) ON UPDATE CASCADE ON DELETE CASCADE
);

CREATE TABLE "card" (
    card_id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    col_id BIGINT NOT NULL,
    title TEXT NOT NULL,
    order_index INT, -- Порядковый номер карточки в колонке
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    cover_file_uuid UUID,
    FOREIGN KEY (col_id) REFERENCES kanban_column(col_id) ON UPDATE CASCADE ON DELETE RESTRICT,
    FOREIGN KEY (cover_file_uuid) REFERENCES user_uploaded_file(file_uuid)
    ON UPDATE CASCADE ON DELETE SET NULL
);

CREATE TABLE tag (
    tag_id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    board_id BIGINT NOT NULL,
    title TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (board_id) REFERENCES board(board_id) ON UPDATE CASCADE ON DELETE CASCADE
);

CREATE TABLE tag_to_card (
    tag_to_card_id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    card_id BIGINT NOT NULL,
    tag_id BIGINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (card_id) REFERENCES Card(card_id) ON UPDATE CASCADE ON DELETE CASCADE,
    FOREIGN KEY (tag_id) REFERENCES tag(tag_id) ON UPDATE CASCADE ON DELETE CASCADE
);

CREATE TABLE card_update (
    card_update_id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    card_id BIGINT NOT NULL,
    is_visible BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by BIGINT,
    assigned_to BIGINT,
    "type" TEXT,
    "text" TEXT,
    attached_file_uuid UUID,
    FOREIGN KEY (card_id) REFERENCES Card(card_id) ON UPDATE CASCADE ON DELETE CASCADE,
    FOREIGN KEY (created_by) REFERENCES "user"(u_id) ON UPDATE CASCADE ON DELETE SET NULL,
    FOREIGN KEY (attached_file_uuid) REFERENCES user_uploaded_file(file_uuid)
    ON UPDATE CASCADE ON DELETE SET NULL
);

CREATE TABLE "notification" (
    notification_id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    u_id BIGINT NOT NULL,
    board_id BIGINT NOT NULL,
    notification_type TEXT,
    content TEXT,
    is_dismissed BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (u_id) REFERENCES "user"(u_id) ON UPDATE CASCADE ON DELETE CASCADE,
    FOREIGN KEY (board_id) REFERENCES board(board_id) ON UPDATE CASCADE ON DELETE CASCADE
);

CREATE TABLE checklist_field (
    checklist_field_id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    card_id BIGINT NOT NULL,
    order_index INT NOT NULL, -- Порядковый номер в чеклисте
    title TEXT NOT NULL,
    is_done BOOLEAN DEFAULT FALSE,
    FOREIGN KEY (card_id) REFERENCES "card"(card_id) ON UPDATE CASCADE ON DELETE CASCADE
);
