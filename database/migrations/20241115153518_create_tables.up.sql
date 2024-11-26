CREATE EXTENSION "uuid-ossp";

-- Create enum type "user_role"
CREATE TYPE "public"."user_role" AS ENUM ('viewer', 'editor', 'editor_chief', 'admin');
-- Create "board" table
CREATE TABLE "public"."board" ("board_id" bigint NOT NULL GENERATED ALWAYS AS IDENTITY, "name" text NOT NULL, "description" text NOT NULL, "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP, "created_by" bigint NULL, "background_image_id" bigint NULL, "updated_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP, PRIMARY KEY ("board_id"));
-- Create "card" table
CREATE TABLE "public"."card" ("card_id" bigint NOT NULL GENERATED ALWAYS AS IDENTITY, "col_id" bigint NOT NULL, "title" text NOT NULL, "order_index" integer NULL, "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP, "updated_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP, "cover_file_id" bigint NULL, PRIMARY KEY ("card_id"));
-- Create "kanban_column" table
CREATE TABLE "public"."kanban_column" ("col_id" bigint NOT NULL GENERATED ALWAYS AS IDENTITY, "board_id" bigint NOT NULL, "title" text NOT NULL, "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP, "updated_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP, "order_index" integer NULL, PRIMARY KEY ("col_id"));
-- Create "user" table
CREATE TABLE "public"."user" ("u_id" bigint NOT NULL GENERATED ALWAYS AS IDENTITY, "nickname" text NOT NULL, "description" text NOT NULL, "joined_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP, "updated_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP, "password_hash" text NOT NULL, "email" text NOT NULL, "avatar_file_id" bigint NULL, PRIMARY KEY ("u_id"), CONSTRAINT "user_email_key" UNIQUE ("email"), CONSTRAINT "user_nickname_key" UNIQUE ("nickname"));
-- Create "user_to_board" table
CREATE TABLE "public"."user_to_board" ("user_to_board_id" bigint NOT NULL GENERATED ALWAYS AS IDENTITY, "u_id" bigint NOT NULL, "board_id" bigint NOT NULL, "added_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP, "updated_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP, "last_visit_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP, "added_by" bigint NULL, "updated_by" bigint NULL, "role" "public"."user_role" NOT NULL DEFAULT 'viewer', PRIMARY KEY ("user_to_board_id"), CONSTRAINT "user_to_board_u_id_board_id_key" UNIQUE ("u_id", "board_id"));
-- Create "user_uploaded_file" table
CREATE TABLE "public"."user_uploaded_file" ("file_id" bigint NOT NULL GENERATED ALWAYS AS IDENTITY, "file_uuid" uuid NOT NULL DEFAULT public.uuid_generate_v4(), "file_extension" text NULL, "size" integer NOT NULL, "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP, "created_by" bigint NULL, PRIMARY KEY ("file_id"));
-- Modify "board" table
ALTER TABLE "public"."board" ADD CONSTRAINT "board_background_image_id_fkey" FOREIGN KEY ("background_image_id") REFERENCES "public"."user_uploaded_file" ("file_id") ON UPDATE CASCADE ON DELETE SET NULL, ADD CONSTRAINT "board_created_by_fkey" FOREIGN KEY ("created_by") REFERENCES "public"."user" ("u_id") ON UPDATE CASCADE ON DELETE SET NULL;
-- Modify "card" table
ALTER TABLE "public"."card" ADD CONSTRAINT "card_col_id_fkey" FOREIGN KEY ("col_id") REFERENCES "public"."kanban_column" ("col_id") ON UPDATE CASCADE ON DELETE CASCADE, ADD CONSTRAINT "card_cover_file_id_fkey" FOREIGN KEY ("cover_file_id") REFERENCES "public"."user_uploaded_file" ("file_id") ON UPDATE CASCADE ON DELETE SET NULL;
-- Modify "kanban_column" table
ALTER TABLE "public"."kanban_column" ADD CONSTRAINT "kanban_column_board_id_fkey" FOREIGN KEY ("board_id") REFERENCES "public"."board" ("board_id") ON UPDATE CASCADE ON DELETE CASCADE;
-- Modify "user" table
ALTER TABLE "public"."user" ADD CONSTRAINT "fk_avatar_file_id" FOREIGN KEY ("avatar_file_id") REFERENCES "public"."user_uploaded_file" ("file_id") ON UPDATE CASCADE ON DELETE SET NULL;
-- Modify "user_to_board" table
ALTER TABLE "public"."user_to_board" ADD CONSTRAINT "user_to_board_added_by_fkey" FOREIGN KEY ("added_by") REFERENCES "public"."user" ("u_id") ON UPDATE CASCADE ON DELETE SET NULL, ADD CONSTRAINT "user_to_board_board_id_fkey" FOREIGN KEY ("board_id") REFERENCES "public"."board" ("board_id") ON UPDATE CASCADE ON DELETE CASCADE, ADD CONSTRAINT "user_to_board_u_id_fkey" FOREIGN KEY ("u_id") REFERENCES "public"."user" ("u_id") ON UPDATE CASCADE ON DELETE CASCADE, ADD CONSTRAINT "user_to_board_updated_by_fkey" FOREIGN KEY ("updated_by") REFERENCES "public"."user" ("u_id") ON UPDATE CASCADE ON DELETE SET NULL;
-- Modify "user_uploaded_file" table
ALTER TABLE "public"."user_uploaded_file" ADD CONSTRAINT "user_uploaded_file_created_by_fkey" FOREIGN KEY ("created_by") REFERENCES "public"."user" ("u_id") ON UPDATE CASCADE ON DELETE SET NULL;
