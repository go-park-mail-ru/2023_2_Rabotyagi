DROP TRIGGER IF EXISTS premium_expire_check ON public."product";

ALTER TABLE public."product"
DROP COLUMN premium_expire,
DROP COLUMN premium_begin;

