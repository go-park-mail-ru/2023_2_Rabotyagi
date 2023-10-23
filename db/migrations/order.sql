-- Table: public.order

DROP TABLE IF EXISTS public."order" CASCADE;
DROP SEQUENCE IF EXISTS order_id_seq;

CREATE SEQUENCE order_id_seq;
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