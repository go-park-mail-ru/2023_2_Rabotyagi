DROP TABLE IF EXISTS public."user" CASCADE;
DROP TABLE IF EXISTS public."product" CASCADE;
DROP TABLE IF EXISTS public."category" CASCADE;
DROP TABLE IF EXISTS public."order" CASCADE;
DROP TABLE IF EXISTS public."image" CASCADE;
DROP TABLE IF EXISTS public."favourite" CASCADE;
DROP SEQUENCE IF EXISTS user_id_seq;
DROP SEQUENCE IF EXISTS product_id_seq;
DROP SEQUENCE IF EXISTS category_id_seq;
DROP SEQUENCE IF EXISTS order_id_seq;
DROP SEQUENCE IF EXISTS image_id_seq;
DROP SEQUENCE IF EXISTS favourite_id_seq;
CREATE SEQUENCE user_id_seq;
CREATE SEQUENCE product_id_seq;
CREATE SEQUENCE category_id_seq;
CREATE SEQUENCE order_id_seq;
CREATE SEQUENCE image_id_seq;
CREATE SEQUENCE favourite_id_seq;

CREATE TABLE public."user"
(
    id       BIGINT DEFAULT NEXTVAL('user_id_seq'::regclass) NOT NULL PRIMARY KEY,
    email    CHARACTER(256) UNIQUE                           NOT NULL,
    phone    CHARACTER(18) UNIQUE                            NOT NULL,
    name     CHARACTER(256),
    pass     CHARACTER(256)                                  NOT NULL,
    birthday TIMESTAMP WITH TIME ZONE
);

CREATE TABLE public."category"
(
    id        BIGINT DEFAULT NEXTVAL('category_id_seq'::regclass) NOT NULL PRIMARY KEY,
    name      CHARACTER(256) UNIQUE                               NOT NULL,
    parent_id BIGINT DEFAULT NULL REFERENCES public.category (id)
);

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

CREATE TABLE public."order"
(
    id          BIGINT                   DEFAULT NEXTVAL('order_id_seq'::regclass) NOT NULL PRIMARY KEY,
    owner_id    BIGINT                                                             NOT NULL REFERENCES public."user" (id),
    product_id  BIGINT                                                             NOT NULL REFERENCES public."product" (id),
    count       SMALLINT                                                           NOT NULL DEFAULT 1 CHECK (count > 0),
    status      SMALLINT                                                           NOT NULL DEFAULT 0,
    create_date TIMESTAMP WITH TIME ZONE DEFAULT NOW()                             NOT NULL,
    update_date TIMESTAMP WITH TIME ZONE DEFAULT NOW()                             NOT NULL,
    close_date  TIMESTAMP WITH TIME ZONE DEFAULT NULL
);

CREATE TABLE public."image"
(
    id         BIGINT DEFAULT NEXTVAL('image_id_seq'::regclass) NOT NULL PRIMARY KEY,
    url        CHARACTER(256)                                   NOT NULL UNIQUE,
    product_id BIGINT                                           NOT NULL REFERENCES public."product" (id) ON DELETE CASCADE
);

CREATE TABLE public."favourite"
(
    id         BIGINT DEFAULT NEXTVAL('favourite_id_seq'::regclass) NOT NULL PRIMARY KEY,
    owner_id   BIGINT                                               NOT NULL REFERENCES public."user" (id),
    product_id BIGINT                                               NOT NULL REFERENCES public."product" (id) ON DELETE CASCADE
);

CREATE OR REPLACE FUNCTION update_date()
    RETURNS TRIGGER AS
$$
BEGIN
    NEW.update_date = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER verify_update_date
    BEFORE UPDATE
    ON public."order"
    FOR EACH ROW
EXECUTE PROCEDURE update_date();