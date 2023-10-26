CREATE SEQUENCE IF NOT EXISTS user_id_seq;
CREATE SEQUENCE IF NOT EXISTS product_id_seq;
CREATE SEQUENCE IF NOT EXISTS category_id_seq;
CREATE SEQUENCE IF NOT EXISTS order_id_seq;
CREATE SEQUENCE IF NOT EXISTS image_id_seq;
CREATE SEQUENCE IF NOT EXISTS favourite_id_seq;

CREATE TABLE IF NOT EXISTS public."user"
(
    id       BIGINT DEFAULT NEXTVAL('user_id_seq'::regclass) NOT NULL PRIMARY KEY,
    email    VARCHAR(256) UNIQUE                             NOT NULL CHECK (email <> ''),
    phone    VARCHAR(18) UNIQUE                              NOT NULL CHECK (phone <> ''),
    name     VARCHAR(256)                                    NOT NULL CHECK (name <> ''),
    pass     VARCHAR(256)                                    NOT NULL CHECK (pass <> ''),
    birthday TIMESTAMP WITH TIME ZONE
);

CREATE TABLE IF NOT EXISTS public."category"
(
    id        BIGINT DEFAULT NEXTVAL('category_id_seq'::regclass) NOT NULL PRIMARY KEY,
    name      VARCHAR(256) UNIQUE                                 NOT NULL CHECK (name <> ''),
    parent_id BIGINT DEFAULT NULL REFERENCES public.category (id)
);

CREATE TABLE IF NOT EXISTS public."product"
(
    id              BIGINT                   DEFAULT NEXTVAL('product_id_seq'::regclass) NOT NULL PRIMARY KEY,
    saler_id        BIGINT                                                               NOT NULL REFERENCES public."user" (id),
    category_id     BIGINT                                                               NOT NULL REFERENCES public."category" (id),
    title           VARCHAR(256)                                                         NOT NULL CHECK (title <> ''),
    description     TEXT                                                                 NOT NULL CHECK (description <> '')
        CONSTRAINT max_len_description CHECK (LENGTH(description) <= 4000),
    price           BIGINT                   DEFAULT 0                                   NOT NULL CHECK (price >= 0),
    created_at      TIMESTAMP WITH TIME ZONE DEFAULT NOW()                               NOT NULL,
    views           INT                      DEFAULT 0                                   NOT NULL CHECK (views >= 0),
    available_count INT                      DEFAULT 0                                   NOT NULL CHECK (available_count >= 0),
    city            VARCHAR(256)                                                         NOT NULL CHECK (city <> ''),
    delivery        BOOLEAN                  DEFAULT FALSE                               NOT NULL,
    safe_deal       BOOLEAN                  DEFAULT FALSE                               NOT NULL,
    is_active       BOOLEAN                  DEFAULT FALSE                               NOT NULL,
    CONSTRAINT not_zero_count_with_active CHECK (not (available_count = 0 and is_active))
);

CREATE TABLE IF NOT EXISTS public."order"
(
    id          BIGINT                   DEFAULT NEXTVAL('order_id_seq'::regclass) NOT NULL PRIMARY KEY,
    owner_id    BIGINT                                                             NOT NULL REFERENCES public."user" (id),
    product_id  BIGINT                                                             NOT NULL REFERENCES public."product" (id),
    count       SMALLINT                                                           NOT NULL DEFAULT 1 CHECK (count > 0),
    status      SMALLINT                                                           NOT NULL DEFAULT 0
        CONSTRAINT status_contract CHECK ( status BETWEEN 0 AND 3),
    created_at  TIMESTAMP WITH TIME ZONE DEFAULT NOW()                             NOT NULL,
    updated_at  TIMESTAMP WITH TIME ZONE DEFAULT NOW()                             NOT NULL,
    closed_at   TIMESTAMP WITH TIME ZONE DEFAULT NULL
);

CREATE TABLE IF NOT EXISTS public."image"
(
    id         BIGINT DEFAULT NEXTVAL('image_id_seq'::regclass) NOT NULL PRIMARY KEY,
    url        VARCHAR(256)                                     NOT NULL CHECK (url <> ''),
    product_id BIGINT                                           NOT NULL REFERENCES public."product" (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS public."favourite"
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

DROP TRIGGER IF EXISTS verify_update_date ON public."order";
CREATE TRIGGER verify_update_date
    BEFORE UPDATE
    ON public."order"
    FOR EACH ROW
EXECUTE PROCEDURE update_date();