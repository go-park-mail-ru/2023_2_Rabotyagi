-- Table: public.category

DROP TABLE IF EXISTS public."category" CASCADE;
DROP SEQUENCE IF EXISTS category_id_seq;

CREATE SEQUENCE category_id_seq;
CREATE TABLE public."category"
(
    id        BIGINT DEFAULT NEXTVAL('category_id_seq'::regclass) NOT NULL PRIMARY KEY,
    name      CHARACTER(256) UNIQUE                               NOT NULL,
    parent_id BIGINT DEFAULT NULL REFERENCES public.category (id)
);