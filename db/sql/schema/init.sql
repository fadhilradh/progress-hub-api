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
  "id" uuid PRIMARY KEY,
  "progress_name" text,
  "progress_value" bigint,
  "range_type" range,
  "range_value" text,
  "created_at" timestamptz,
  "updated_at" timestamptz,
  "user_id" uuid
);

CREATE TABLE "user" (
  "id" uuid PRIMARY KEY,
  "username" text UNIQUE NOT NULL,
  "email" text UNIQUE NOT NULL,
  "password" text NOT NULL,
  "created_at" timestamptz,
  "updated_at" timestamptz DEFAULT (now()),
  "role" role,
  "photo_profile_url" text
);

ALTER TABLE "progress" ADD FOREIGN KEY ("user_id") REFERENCES "user" ("id");
