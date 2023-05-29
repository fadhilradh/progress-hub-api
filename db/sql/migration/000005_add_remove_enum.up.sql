ALTER TABLE charts ALTER COLUMN range_type TYPE text;
ALTER TABLE users ALTER COLUMN role TYPE text;

DROP TYPE "range";
DROP TYPE "role";