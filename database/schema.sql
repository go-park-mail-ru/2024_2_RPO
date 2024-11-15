CREATE EXTENSION "uuid-ossp";
CREATE TABLE "user" (
    u_id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    nickname TEXT NOT NULL UNIQUE,
    "description" TEXT NOT NULL,
    joined_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    password_hash TEXT NOT NULL,
    email TEXT UNIQUE NOT NULL,
    avatar_file_id BIGINT
);
CREATE TYPE USER_ROLE AS ENUM(
    'viewer',
    'editor',
    'editor_chief',
    'admin'
);
CREATE TABLE user_uploaded_file(
    file_id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    file_uuid UUID NOT NULL DEFAULT uuid_generate_v4(),
    file_extension TEXT,
    "size" INTEGER NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by BIGINT,
    FOREIGN KEY (created_by) REFERENCES "user"(u_id) ON UPDATE CASCADE ON DELETE
    SET NULL
);
ALTER TABLE "user"
ADD CONSTRAINT fk_avatar_file_id FOREIGN KEY (avatar_file_id) REFERENCES user_uploaded_file(file_id) ON UPDATE CASCADE ON DELETE
SET NULL;
CREATE TABLE board (
    board_id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    "name" TEXT NOT NULL,
    "description" TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by BIGINT,
    background_image_id BIGINT,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (created_by) REFERENCES "user"(u_id) ON UPDATE CASCADE ON DELETE
    SET NULL,
        FOREIGN KEY (background_image_id) REFERENCES user_uploaded_file(file_id) ON UPDATE CASCADE ON DELETE
    SET NULL
);
CREATE TABLE user_to_board (
    user_to_board_id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    u_id BIGINT NOT NULL,
    board_id BIGINT NOT NULL,
    added_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    last_visit_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    added_by BIGINT,
    updated_by BIGINT,
    "role" USER_ROLE NOT NULL DEFAULT 'viewer',
    FOREIGN KEY (u_id) REFERENCES "user"(u_id) ON UPDATE CASCADE ON DELETE CASCADE,
    FOREIGN KEY (board_id) REFERENCES board(board_id) ON UPDATE CASCADE ON DELETE CASCADE,
    FOREIGN KEY (added_by) REFERENCES "user"(u_id) ON UPDATE CASCADE ON DELETE
    SET NULL,
        FOREIGN KEY (updated_by) REFERENCES "user"(u_id) ON UPDATE CASCADE ON DELETE
    SET NULL,
        UNIQUE(u_id, board_id)
);
CREATE TABLE kanban_column (
    col_id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    board_id BIGINT NOT NULL,
    title TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    order_index INT,
    -- Порядковый номер колонки на доске
    FOREIGN KEY (board_id) REFERENCES board(board_id) ON UPDATE CASCADE ON DELETE CASCADE
);
CREATE TABLE "card" (
    card_id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    col_id BIGINT NOT NULL,
    title TEXT NOT NULL,
    order_index INT,
    -- Порядковый номер карточки в колонке
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    cover_file_id BIGINT,
    FOREIGN KEY (col_id) REFERENCES kanban_column(col_id) ON UPDATE CASCADE ON DELETE CASCADE,
    FOREIGN KEY (cover_file_id) REFERENCES user_uploaded_file(file_id) ON UPDATE CASCADE ON DELETE
    SET NULL
);
