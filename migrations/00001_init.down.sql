-- +migrate Down
BEGIN;

DROP TABLE IF EXISTS public."user" CASCADE;
DROP TABLE IF EXISTS public."table" CASCADE;
DROP TABLE IF EXISTS public."note" CASCADE;
DROP TABLE IF EXISTS public."category" CASCADE;

END;