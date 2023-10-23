-- Table: public.user

DROP TABLE IF EXISTS public."user" CASCADE;
DROP SEQUENCE IF EXISTS user_id_seq;

CREATE SEQUENCE user_id_seq;
CREATE TABLE public."user"
(
    id       BIGINT DEFAULT NEXTVAL('user_id_seq'::regclass) NOT NULL PRIMARY KEY,
    email    CHARACTER(256) UNIQUE                           NOT NULL,
    phone    CHARACTER(18) UNIQUE                            NOT NULL,
    name     CHARACTER(256),
    pass     CHARACTER(256)                                  NOT NULL,
    birthday TIMESTAMP WITH TIME ZONE
);