ALTER TABLE public."product"
    ADD COLUMN premium_begin TIMESTAMP WITH TIME ZONE DEFAULT NULL,
ADD COLUMN premium_expire TIMESTAMP WITH TIME ZONE DEFAULT NULL;

CREATE OR REPLACE FUNCTION update_premium_expire_check()
    RETURNS TRIGGER AS $$
BEGIN
    IF NEW.premium_expire < CURRENT_TIMESTAMP THEN
        UPDATE public."product"
        SET premium = FALSE, premium_expire = NULL, premium_begin = NULL
        WHERE id = NEW.id;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS premium_expire_check ON public."product";
CREATE TRIGGER premium_expire_check
    BEFORE INSERT OR UPDATE
    ON public."product"
    FOR EACH ROW
EXECUTE FUNCTION update_premium_expire_check();