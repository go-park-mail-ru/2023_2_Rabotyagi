-- premium_status = 0 not premium
-- premium_status = 1 pending
-- premium_status = 2 waiting_for_capture
-- premium_status = 3 succeeded
-- premium_status = 4 canceled

ALTER TABLE public."product"
    ADD COLUMN premium_status INT DEFAULT 0 NOT NULL,
    ADD CONSTRAINT correctness_status_premium CHECK (premium_status >= 0 AND premium_status <= 4);

UPDATE public."product" SET premium_status=0 WHERE premium=FALSE;
UPDATE public."product" SET premium_status=3 WHERE premium=TRUE;

ALTER TABLE public."product"
    DROP COLUMN premium;

