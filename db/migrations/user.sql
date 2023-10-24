-- Table: public.user

DROP TABLE IF EXISTS public."user" CASCADE;
DROP SEQUENCE IF EXISTS user_id_seq;

CREATE SEQUENCE user_id_seq;
CREATE TABLE public."user"
(
    id       BIGINT DEFAULT NEXTVAL('user_id_seq'::regclass) NOT NULL PRIMARY KEY,
    email    VARCHAR(256) UNIQUE                             NOT NULL CHECK (email <> ''),
    phone    VARCHAR(18) UNIQUE                              NOT NULL CHECK (phone <> ''),
    name     VARCHAR(256)                                    NOT NULL CHECK (name <> ''),
    pass     VARCHAR(256)                                    NOT NULL CHECK (pass <> ''),
    birthday TIMESTAMP WITH TIME ZONE
);