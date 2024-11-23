-- Create enum type "question_type"
CREATE TYPE "public"."question_type" AS ENUM ('answer_text', 'answer_rating');
-- Modify "card" table
ALTER TABLE "public"."card" DROP CONSTRAINT "card_col_id_order_index_key";
-- Modify "checklist_field" table
ALTER TABLE "public"."checklist_field" DROP CONSTRAINT "checklist_field_card_id_order_index_key";
-- Create "csat_question" table
CREATE TABLE "public"."csat_question" ("question_id" bigint NOT NULL GENERATED ALWAYS AS IDENTITY, "is_active" boolean NOT NULL DEFAULT true, "question_text" text NOT NULL, "type" "public"."question_type" NOT NULL, "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP, PRIMARY KEY ("question_id"));
-- Modify "user" table
ALTER TABLE "public"."user" ADD COLUMN "csat_poll_dt" timestamptz NOT NULL;
-- Create "csat_results" table
CREATE TABLE "public"."csat_results" ("result_id" bigint NOT NULL GENERATED ALWAYS AS IDENTITY, "question_id" bigint NOT NULL, "rating" integer NULL, "comment" text NULL, "u_id" bigint NOT NULL, "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP, PRIMARY KEY ("result_id"), CONSTRAINT "csat_results_question_id_fkey" FOREIGN KEY ("question_id") REFERENCES "public"."csat_question" ("question_id") ON UPDATE CASCADE ON DELETE CASCADE, CONSTRAINT "csat_results_u_id_fkey" FOREIGN KEY ("u_id") REFERENCES "public"."user" ("u_id") ON UPDATE CASCADE ON DELETE CASCADE);
