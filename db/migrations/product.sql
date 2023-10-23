-- Table: public.product
DROP TABLE IF EXISTS public."product" CASCADE;
DROP SEQUENCE IF EXISTS product_id_seq;

CREATE SEQUENCE product_id_seq;
CREATE TABLE public."product"
(
    id              BIGINT                   DEFAULT NEXTVAL('product_id_seq'::regclass) NOT NULL PRIMARY KEY,
    saler_id        BIGINT                                                               NOT NULL REFERENCES public."user" (id),
    category_id     BIGINT                                                               NOT NULL REFERENCES public."category" (id),
    title           CHARACTER(256)                                                       NOT NULL,
    description     TEXT                                                                 NOT NULL,
    price           BIGINT                   DEFAULT 0 CHECK (price >= 0),
    create_date     TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    views           INT                      DEFAULT 0 CHECK (views >= 0),
    in_favourites   INT                      DEFAULT 0 CHECK (in_favourites >= 0),
    available_count INT                      DEFAULT 0 CHECK (available_count >= 0),
    city            CHARACTER(256)                                                       NOT NULL,
    delivery        BOOLEAN                  DEFAULT FALSE,
    safe_dial       BOOLEAN                  DEFAULT FALSE,
    is_active       BOOLEAN                  DEFAULT FALSE,
    constraint not_null_good_count CHECK (not (available_count = 0 and is_active))
);