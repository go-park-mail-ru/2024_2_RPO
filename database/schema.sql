CREATE EXTENSION "uuid-ossp";

CREATE TABLE user_uploaded_file(
    file_id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    file_uuid UUID NOT NULL DEFAULT uuid_generate_v4(),
    file_hash TEXT,
    file_extension TEXT,
    "size" INTEGER NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE "user" (
    u_id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    nickname TEXT NOT NULL UNIQUE,
    "description" TEXT NOT NULL,
    joined_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    password_hash TEXT NOT NULL,
    email TEXT UNIQUE NOT NULL,
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
    board_id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    "name" TEXT NOT NULL,
    "description" TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by BIGINT,
    background_image_id BIGINT,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (created_by) REFERENCES "user"(u_id) ON UPDATE CASCADE ON DELETE SET NULL,
    FOREIGN KEY (background_image_id) REFERENCES user_uploaded_file(file_id) ON UPDATE CASCADE ON DELETE SET NULL
);

CREATE TABLE user_to_board (
    u_id BIGINT NOT NULL,
    board_id BIGINT NOT NULL,
    added_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    last_visit_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    added_by BIGINT,
    updated_by BIGINT,
    invite_link_uuid UUID,
    "role" user_role NOT NULL DEFAULT 'viewer',
    FOREIGN KEY (u_id) REFERENCES "user"(u_id) ON UPDATE CASCADE ON DELETE CASCADE,
    FOREIGN KEY (board_id) REFERENCES board(board_id) ON UPDATE CASCADE ON DELETE CASCADE,
    FOREIGN KEY (added_by) REFERENCES "user"(u_id) ON UPDATE CASCADE ON DELETE SET NULL,
    FOREIGN KEY (updated_by) REFERENCES "user"(u_id) ON UPDATE CASCADE ON DELETE SET NULL,
    PRIMARY KEY(u_id, board_id)
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
    card_uuid UUID NOT NULL DEFAULT uuid_generate_v4(), -- UUID для ссылки на карточку
    col_id BIGINT NOT NULL,
    title TEXT NOT NULL,
    order_index INTEGER, -- Порядковый номер карточки в колонке
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    cover_file_id BIGINT,
    deadline TIMESTAMPTZ,
    is_done BOOLEAN DEFAULT FALSE, -- Задана, когда задан deadline
    FOREIGN KEY (col_id) REFERENCES kanban_column(col_id) ON UPDATE CASCADE ON DELETE CASCADE,
    FOREIGN KEY (cover_file_id) REFERENCES user_uploaded_file(file_id) ON UPDATE CASCADE ON DELETE SET NULL
);

CREATE TABLE card_attachment (
    attachment_id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    card_id BIGINT NOT NULL,
    file_id BIGINT NOT NULL,
    original_name TEXT NOT NULL,
    attached_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    attached_by BIGINT NOT NULL,
    FOREIGN KEY (card_id) REFERENCES card(card_id) ON UPDATE CASCADE ON DELETE CASCADE,
    FOREIGN KEY (file_id) REFERENCES user_uploaded_file(file_id) ON UPDATE CASCADE ON DELETE CASCADE,
    FOREIGN KEY (attached_by) REFERENCES "user"(u_id) ON UPDATE CASCADE ON DELETE CASCADE
);

CREATE TABLE checklist_field (
    checklist_field_id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    card_id BIGINT NOT NULL,
    title TEXT NOT NULL,
    is_done BOOLEAN NOT NULL DEFAULT FALSE,
    order_index INTEGER,
    FOREIGN KEY (card_id) REFERENCES card(card_id) ON UPDATE CASCADE ON DELETE CASCADE
);

CREATE TABLE card_comment (
    comment_id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    card_id BIGINT NOT NULL,
    title TEXT NOT NULL,
    created_by BIGINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (card_id) REFERENCES card(card_id) ON UPDATE CASCADE ON DELETE CASCADE,
    FOREIGN KEY (created_by) REFERENCES "user"(u_id) ON UPDATE CASCADE ON DELETE CASCADE
);

CREATE TABLE card_user_assignment (
    assignment_id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    card_id BIGINT NOT NULL,
    u_id BIGINT NOT NULL,
    FOREIGN KEY (card_id) REFERENCES card(card_id) ON UPDATE CASCADE ON DELETE CASCADE,
    FOREIGN KEY (u_id) REFERENCES "user"(u_id) ON UPDATE CASCADE ON DELETE CASCADE
);

-- CREATE TYPE card_update_type AS ENUM (
--     'created',
--     'moved',
--     'edited_title',
--     'attached_file',
--     'removed_file',
--     'replaced_file',
--     'assigned_member',
--     'deassigned_member',
--     'add_checklist_field',
--     'remove_checklist_field',
--     'edit_checklist_field',
--     'mark_checklist_field_done',
--     'mark_checklist_field_undone',
--     'left_comment',
--     'set_deadline',
--     'remove_deadline'
-- );

-- CREATE TABLE card_update (
--     card_update_id BIGINT PRIMARY KEY ALWAYS GENERATED AS IDENTITY,
--     card_id BIGINT NOT NULL,
--     update_type card_update_type NOT NULL,
--     updated_by BIGINT NOT NULL,
--     updated_at TIMESTAMPTZ,
--     assigned_member BIGINT,
--     checklist_field_prev_name TEXT,
--     checklist_field_current_name TEXT,
--     checklist_prev_state BOOLEAN,
--     file_name TEXT,
--     prev_column TEXT,
--     curr_column TEXT,
--     prev_title TEXT,
--     curr_title TEXT,
--     deadline TIMESTAMPTZ,
--     FOREIGN KEY (card_id) REFERENCES "card"(card_id) ON UPDATE CASCADE ON DELETE CASCADE,
--     FOREIGN KEY (updated_by) REFERENCES "user"(u_id) ON UPDATE CASCADE ON DELETE CASCADE,
--     FOREIGN KEY (assigned_member) REFERENCES "user"(u_id) ON UPDATE CASCADE ON DELETE CASCADE
-- );

CREATE TABLE tarasovxx(
    tarasovxx_id BIGINT GENERATED ALWAYS AS IDENTITY
);
