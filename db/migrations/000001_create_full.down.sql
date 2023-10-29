DROP TRIGGER IF EXISTS verify_updated_at ON public."order";
DROP FUNCTION IF EXISTS updated_at_now;

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