CREATE TYPE "range" AS ENUM (
  'daily',
  'weekly',
  'monthly',
  'yearly'
);

CREATE TYPE "role" AS ENUM (
  'user',
  'admin',
  'superadmin'
);

CREATE TABLE "progress" (
  "id" uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  "progress_name" text NOT NULL,
  "progress_value" bigint NOT NULL,
  "range_type" range NOT NULL,
  "range_value" text NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL DEFAULT (now()),
  "user_id" uuid NOT NULL
);

CREATE TABLE "users" (
  "id" uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  "username" text UNIQUE NOT NULL,
  "email" text UNIQUE NOT NULL,
  "password" text NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL DEFAULT (now()),
  "role" role NOT NULL,
  "photo_profile_url" text
);

ALTER TABLE "progress" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");
