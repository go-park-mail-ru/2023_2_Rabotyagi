CREATE SEQUENCE IF NOT EXISTS user_id_seq;
CREATE SEQUENCE IF NOT EXISTS product_id_seq;
CREATE SEQUENCE IF NOT EXISTS category_id_seq;
CREATE SEQUENCE IF NOT EXISTS order_id_seq;
CREATE SEQUENCE IF NOT EXISTS image_id_seq;
CREATE SEQUENCE IF NOT EXISTS favourite_id_seq;

CREATE TABLE IF NOT EXISTS public."user"
(
    id       BIGINT DEFAULT NEXTVAL('user_id_seq'::regclass) NOT NULL PRIMARY KEY,
    email    TEXT UNIQUE                                     NOT NULL CHECK (email <> '') CHECK (length(email) <= 256),
    phone    TEXT UNIQUE                                     NOT NULL CHECK (phone <> '') CHECK (length(phone) <= 18),
    name     TEXT                                            NOT NULL CHECK (name <> '') CHECK (length(name) <= 256),
    password TEXT                                            NOT NULL CHECK (password <> '') CHECK (length(password) <= 256),
    birthday TIMESTAMP WITH TIME ZONE
);

CREATE TABLE IF NOT EXISTS public."category"
(
    id        BIGINT DEFAULT NEXTVAL('category_id_seq'::regclass) NOT NULL PRIMARY KEY,
    name      TEXT UNIQUE                                         NOT NULL CHECK (name <> '') CHECK (length(name) <= 256),
    parent_id BIGINT DEFAULT NULL REFERENCES public."category" (id)
);

CREATE TABLE IF NOT EXISTS public."product"
(
    id              BIGINT                   DEFAULT NEXTVAL('product_id_seq'::regclass) NOT NULL PRIMARY KEY,
    saler_id        BIGINT                                                               NOT NULL REFERENCES public."user" (id),
    category_id     BIGINT                                                               NOT NULL REFERENCES public."category" (id),
    title           TEXT                                                                 NOT NULL CHECK (title <> '') CHECK (length(title) <= 256),
    description     TEXT                                                                 NOT NULL CHECK (description <> '')
        CONSTRAINT max_len_description CHECK (LENGTH(description) <= 4000),
    price           BIGINT                   DEFAULT 0                                   NOT NULL CHECK (price >= 0),
    created_at      TIMESTAMP WITH TIME ZONE DEFAULT NOW()                               NOT NULL,
    views           INT                      DEFAULT 0                                   NOT NULL CHECK (views >= 0),
    available_count INT                      DEFAULT 0                                   NOT NULL CHECK (available_count >= 0),
city            TEXT                                                                     NOT NULL CHECK (city <> '') CHECK (length(city) <= 256),
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
    url        TEXT                                             NOT NULL CHECK (url <> '') CHECK (length(url) <= 256),
    product_id BIGINT                                           NOT NULL REFERENCES public."product" (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS public."favourite"
(
    id         BIGINT DEFAULT NEXTVAL('favourite_id_seq'::regclass) NOT NULL PRIMARY KEY,
    owner_id   BIGINT                                               NOT NULL REFERENCES public."user" (id),
    product_id BIGINT                                               NOT NULL REFERENCES public."product" (id) ON DELETE CASCADE
);

CREATE OR REPLACE FUNCTION updated_at_now()
    RETURNS TRIGGER AS
$$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS verify_updated_at ON public."order";
CREATE TRIGGER verify_updated_at
    BEFORE UPDATE
    ON public."order"
    FOR EACH ROW
EXECUTE PROCEDURE updated_at_now();