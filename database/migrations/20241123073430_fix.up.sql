-- Modify "card" table
ALTER TABLE "public"."card" ALTER COLUMN "is_done" SET NOT NULL, ADD COLUMN "title" text NOT NULL, ADD CONSTRAINT "card_col_id_order_index_key" UNIQUE ("col_id", "order_index") DEFERRABLE INITIALLY DEFERRED;
-- Modify "card_attachment" table
ALTER TABLE "public"."card_attachment" ADD COLUMN "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP;
-- Modify "card_comment" table
ALTER TABLE "public"."card_comment" ADD COLUMN "is_edited" boolean NOT NULL DEFAULT false;
-- Modify "checklist_field" table
ALTER TABLE "public"."checklist_field" ADD COLUMN "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP, ADD CONSTRAINT "checklist_field_card_id_order_index_key" UNIQUE ("card_id", "order_index") DEFERRABLE INITIALLY DEFERRED;
-- Modify "kanban_column" table
ALTER TABLE "public"."kanban_column" ALTER COLUMN "order_index" SET NOT NULL, ADD CONSTRAINT "kanban_column_board_id_order_index_key" UNIQUE ("board_id", "order_index") DEFERRABLE INITIALLY DEFERRED;
-- Modify "user" table
ALTER TABLE "public"."user" DROP COLUMN "description";
