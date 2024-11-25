-- Modify "board" table
ALTER TABLE "public"."board" DROP COLUMN "description";
-- Modify "card" table
ALTER TABLE "public"."card" DROP COLUMN "title";
-- Modify "user" table
ALTER TABLE "public"."user" ALTER COLUMN "password_hash" DROP NOT NULL;
