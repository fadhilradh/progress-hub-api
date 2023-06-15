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

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS "progress" (
  "id" uuid PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "chart_id" uuid,
  "progress_value" bigint NOT NULL ,
  "range_value" text NOT NULL,
  "created_at" timestamptz,
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE IF NOT EXISTS "charts" (
  "id" uuid PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "user_id" uuid,
  "created_at" timestamptz,
  "updated_at" timestamptz NOT NULL DEFAULT (now()),
  "range_type" range NOT NULL,
  "progress_name" text NOT NULL
);

CREATE TABLE IF NOT EXISTS "users" (
  "id" uuid PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "username" text UNIQUE NOT NULL,
  "email" text UNIQUE NOT NULL,
  "password" text NOT NULL,
  "created_at" timestamptz,
  "updated_at" timestamptz NOT NULL DEFAULT (now()),
  "role" role NOT NULL,
  "photo_profile_url" text
);

ALTER TABLE "progress" ADD FOREIGN KEY ("chart_id") REFERENCES "charts" ("id");

ALTER TABLE "charts" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");
