ALTER TABLE public."product"
    ADD COLUMN premium BOOLEAN DEFAULT FALSE NOT NULL;

UPDATE public."product" SET premium=FALSE WHERE premium_status=0;
UPDATE public."product" SET premium=TRUE WHERE premium_status=3;

ALTER TABLE public."product"
    DROP CONSTRAINT correctness_status_premium,
    DROP COLUMN premium_status;