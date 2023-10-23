-- Table: public.favourite

DROP TABLE IF EXISTS public."favourite" CASCADE;
DROP SEQUENCE IF EXISTS favourite_id_seq;

CREATE SEQUENCE favourite_id_seq;
CREATE TABLE public."favourite"
(
    id         BIGINT DEFAULT NEXTVAL('favourite_id_seq'::regclass) NOT NULL PRIMARY KEY,
    owner_id   BIGINT                                               NOT NULL REFERENCES public."user" (id),
    product_id BIGINT                                               NOT NULL REFERENCES public."product" (id) ON DELETE CASCADE
);