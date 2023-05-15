-- +migrate Up
BEGIN;
SET statement_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = ON;
SET check_function_bodies = FALSE;
SET client_min_messages = WARNING;
SET search_path = public, extensions;
SET default_tablespace = '';
SET default_with_oids = FALSE;

CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE public."user"
(
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name           TEXT NOT NULL,
    email          TEXT NOT NULL,
    pwd            TEXT NOT NULL,
    payment_status BOOLEAN DEFAULT false
);
CREATE TABLE public."category"
(
    id   SERIAL PRIMARY KEY,
    name TEXT,
    back_color TEXT,
    word_color TEXT
);

CREATE TABLE public."table"
(
    id            SERIAL PRIMARY KEY,
    user_id       UUID REFERENCES public."user" (id),
    capacity      INT DEFAULT 1
);

CREATE TABLE public."note"
(
    id            SERIAL PRIMARY KEY,
    table_id      INT REFERENCES public."table" (id),
    category_id   INT REFERENCES public."category" (id),
    deadline      TIMESTAMPTZ,
    title         TEXT,
    description   TEXT

);

INSERT INTO public."category" (name, back_color, word_color)
VALUES ('работа', '#F14158', '#2F32FA');
INSERT INTO public."category" (name, back_color, word_color)
VALUES ('приколы', '#8375DC', '#B5AD91');
INSERT INTO public."category" (name, back_color, word_color)
VALUES ('учеба', '#338495', '#F4C636');

COMMIT;