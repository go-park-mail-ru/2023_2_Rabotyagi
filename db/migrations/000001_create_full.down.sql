DROP TRIGGER IF EXISTS verify_updated_at ON public."order";
DROP FUNCTION IF EXISTS updated_at_now;
DROP TRIGGER IF EXISTS check_not_zero_count_with_active_product ON public."product";
DROP FUNCTION IF EXISTS not_zero_count_with_active_product;

DROP TABLE IF EXISTS public."user" CASCADE;
DROP TABLE IF EXISTS public."product" CASCADE;
DROP TABLE IF EXISTS public."view" CASCADE;
DROP TABLE IF EXISTS public."city" CASCADE;
DROP TABLE IF EXISTS public."category" CASCADE;
DROP TABLE IF EXISTS public."order" CASCADE;
DROP TABLE IF EXISTS public."image" CASCADE;
DROP TABLE IF EXISTS public."favourite" CASCADE;

DROP INDEX IF EXISTS product_description_search_idx;
DROP INDEX IF EXISTS product_title_search_idx;

DROP SEQUENCE IF EXISTS user_id_seq;
DROP SEQUENCE IF EXISTS product_id_seq;
DROP SEQUENCE IF EXISTS view_id_seq;
DROP SEQUENCE IF EXISTS city_id_seq;
DROP SEQUENCE IF EXISTS category_id_seq;
DROP SEQUENCE IF EXISTS order_id_seq;
DROP SEQUENCE IF EXISTS image_id_seq;
DROP SEQUENCE IF EXISTS favourite_id_seq;