-- Table: public.product
DROP TABLE IF EXISTS public."product" CASCADE;
DROP SEQUENCE IF EXISTS product_id_seq;

CREATE SEQUENCE product_id_seq;

CREATE TABLE public."product"
(
    id              BIGINT                   DEFAULT NEXTVAL('product_id_seq'::regclass) NOT NULL PRIMARY KEY,
    saler_id        BIGINT                                                               NOT NULL REFERENCES public."user" (id),
    category_id     BIGINT                                                               NOT NULL REFERENCES public."category" (id),
    title           VARCHAR(256)                                                         NOT NULL CHECK (title <> ''),
    description     TEXT                                                                 NOT NULL CHECK (description <> '')
        CONSTRAINT max_len_description CHECK (LENGTH(description) <= 4000),
    price           BIGINT                   DEFAULT 0                                   NOT NULL CHECK (price >= 0),
    create_date     TIMESTAMP WITH TIME ZONE DEFAULT NOW()                               NOT NULL,
    views           INT                      DEFAULT 0                                   NOT NULL CHECK (views >= 0),
    in_favourites   INT                      DEFAULT 0                                   NOT NULL CHECK (in_favourites >= 0),
    available_count INT                      DEFAULT 0                                   NOT NULL CHECK (available_count >= 0),
    city            VARCHAR(256)                                                         NOT NULL CHECK (city <> ''),
    delivery        BOOLEAN                  DEFAULT FALSE                               NOT NULL,
    safe_dial       BOOLEAN                  DEFAULT FALSE                               NOT NULL,
    is_active       BOOLEAN                  DEFAULT FALSE                               NOT NULL,
    CONSTRAINT not_zero_count_with_active CHECK (not (available_count = 0 and is_active))
);