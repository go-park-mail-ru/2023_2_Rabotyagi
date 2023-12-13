ALTER TABLE public."product"
ADD COLUMN premium_expire TIMESTAMP WITH TIMEZONE DEFAULT NULL;

DROP TRIGGER IF EXISTS premium_expire_check ON public."product";
CREATE TRIGGER premium_expire_check
    BEFORE READ
    ON public."product"
    FOR EACH ROW
BEGIN
  IF NEW.premium_expire < CURRENT_TIMESTAMP THEN
UPDATE product
SET premium = FALSE, premium_expire = NULL
WHERE product_id = NEW.product_id;
END IF;
END;