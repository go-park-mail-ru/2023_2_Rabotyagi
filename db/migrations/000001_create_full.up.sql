CREATE SEQUENCE IF NOT EXISTS user_id_seq;
CREATE SEQUENCE IF NOT EXISTS product_id_seq;
CREATE SEQUENCE IF NOT EXISTS city_id_seq;
CREATE SEQUENCE IF NOT EXISTS category_id_seq;
CREATE SEQUENCE IF NOT EXISTS order_id_seq;
CREATE SEQUENCE IF NOT EXISTS image_id_seq;
CREATE SEQUENCE IF NOT EXISTS favourite_id_seq;

CREATE TABLE IF NOT EXISTS public."user"
(
    id         BIGINT                   DEFAULT NEXTVAL('user_id_seq'::regclass) NOT NULL PRIMARY KEY,
    email      TEXT UNIQUE                                                       NOT NULL CHECK (email <> '')
    CONSTRAINT max_len_email CHECK (LENGTH(email) <= 256),
    phone      TEXT UNIQUE DEFAULT NULL
    CONSTRAINT max_len_phone CHECK (LENGTH(phone) <= 18),
    name       TEXT UNIQUE DEFAULT NULL
    CONSTRAINT max_len_name CHECK (LENGTH(name) <= 256),
    password   TEXT                                                              NOT NULL CHECK (password <> '')
    CONSTRAINT max_len_password CHECK (LENGTH(password) <= 256),
    birthday   TIMESTAMP WITH TIME ZONE,
                             avatar     TEXT UNIQUE
                             CONSTRAINT max_len_avatar CHECK (LENGTH(avatar) <= 256),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()                            NOT NULL
);

CREATE TABLE IF NOT EXISTS public."category"
(
    id        BIGINT DEFAULT NEXTVAL('category_id_seq'::regclass) NOT NULL PRIMARY KEY,
    name      TEXT UNIQUE                                         NOT NULL CHECK (name <> '')
        CONSTRAINT max_len_name CHECK (LENGTH(name) <= 256),
    parent_id BIGINT DEFAULT NULL REFERENCES public."category" (id)
);

CREATE TABLE IF NOT EXISTS public."city"
(
    id              BIGINT            DEFAULT NEXTVAL('city_id_seq'::regclass) NOT NULL PRIMARY KEY,
    name            TEXT                                                       NOT NULL CHECK (name <> '')
    CONSTRAINT max_len_name CHECK (LENGTH(name) <= 256)
);

CREATE TABLE IF NOT EXISTS public."product"
(
    id              BIGINT                   DEFAULT NEXTVAL('product_id_seq'::regclass) NOT NULL PRIMARY KEY,
    saler_id        BIGINT                                                               NOT NULL REFERENCES public."user" (id),
    category_id     BIGINT                                                               NOT NULL REFERENCES public."category" (id),
    city_id         BIGINT                                                               NOT NULL REFERENCES public."city" (id),
    title           TEXT                                                                 NOT NULL CHECK (title <> '')
        CONSTRAINT max_len_title CHECK (LENGTH(title) <= 256),
    description     TEXT                                                                 NOT NULL CHECK (description <> '')
        CONSTRAINT max_len_description CHECK (LENGTH(description) <= 4000),
    price           BIGINT                   DEFAULT 0                                   NOT NULL
        CONSTRAINT not_negative_price CHECK (price >= 0),
    created_at      TIMESTAMP WITH TIME ZONE DEFAULT NOW()                               NOT NULL,
    views           INT                      DEFAULT 0                                   NOT NULL
        CONSTRAINT not_negative_views CHECK (views >= 0),
    available_count INT                      DEFAULT 0                                   NOT NULL
        CONSTRAINT max_len_available_count CHECK (available_count >= 0),
    delivery        BOOLEAN                  DEFAULT FALSE                               NOT NULL,
    safe_deal       BOOLEAN                  DEFAULT FALSE                               NOT NULL,
    is_active       BOOLEAN                  DEFAULT TRUE                                NOT NULL,
    CONSTRAINT not_zero_count_with_active CHECK (not (available_count = 0 and is_active))
);

CREATE TABLE IF NOT EXISTS public."order"
(
    id         BIGINT                   DEFAULT NEXTVAL('order_id_seq'::regclass) NOT NULL PRIMARY KEY,
    owner_id   BIGINT                                                             NOT NULL REFERENCES public."user" (id),
    product_id BIGINT                                                             NOT NULL REFERENCES public."product" (id),
    count      SMALLINT                                                           NOT NULL DEFAULT 1
        CONSTRAINT positive_count CHECK (count > 0),
    status     SMALLINT                                                           NOT NULL DEFAULT 0
        CONSTRAINT status_contract CHECK ( status BETWEEN 0 AND 3),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()                             NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()                             NOT NULL,
    closed_at  TIMESTAMP WITH TIME ZONE DEFAULT NULL
);

CREATE TABLE IF NOT EXISTS public."image"
(
    id         BIGINT DEFAULT NEXTVAL('image_id_seq'::regclass) NOT NULL PRIMARY KEY,
    url        TEXT                                             NOT NULL CHECK (url <> '')
        CONSTRAINT max_len_url CHECK (LENGTH(url) <= 256),
    product_id BIGINT                                           NOT NULL REFERENCES public."product" (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS public."favourite"
(
    id         BIGINT DEFAULT NEXTVAL('favourite_id_seq'::regclass) NOT NULL PRIMARY KEY,
    owner_id   BIGINT                                               NOT NULL REFERENCES public."user" (id),
    product_id BIGINT                                               NOT NULL REFERENCES public."product" (id) ON DELETE CASCADE,
    CONSTRAINT uniq_together_product_id_owner_id unique (owner_id, product_id)
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

CREATE OR REPLACE FUNCTION not_zero_count_with_active_product()
    RETURNS TRIGGER AS
$$
BEGIN
    IF NEW.available_count = 0 THEN
        NEW.is_active = false;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;


DROP TRIGGER IF EXISTS check_not_zero_count_with_active_product ON public."product";
CREATE TRIGGER check_not_zero_count_with_active_product
    BEFORE UPDATE
    ON public."product"
    FOR EACH ROW
EXECUTE PROCEDURE not_zero_count_with_active_product();