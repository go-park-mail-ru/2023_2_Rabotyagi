DROP TABLE IF EXISTS public."image";
DROP SEQUENCE IF EXISTS image_id_seq;

CREATE SEQUENCE image_id_seq;
CREATE TABLE public."image"
(
    id         BIGINT DEFAULT NEXTVAL('image_id_seq'::regclass) NOT NULL PRIMARY KEY,
    url        VARCHAR(256)                                     NOT NULL CHECK (url <> ''),
    product_id BIGINT                                           NOT NULL REFERENCES public."product" (id) ON DELETE CASCADE
);