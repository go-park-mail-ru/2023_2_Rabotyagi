DROP FUNCTION IF EXISTS update_data;
DROP TRIGGER IF EXISTS verify_update_date ON public."order";

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
