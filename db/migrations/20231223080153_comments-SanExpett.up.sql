CREATE SEQUENCE IF NOT EXISTS comment_id_seq;

CREATE TABLE IF NOT EXISTS public."comment"
(
    id              BIGINT            DEFAULT NEXTVAL('comment_id_seq'::regclass) NOT NULL PRIMARY KEY,
    sender_id       BIGINT                                                        NOT NULL REFERENCES public."user" (id) ON DELETE CASCADE,
    recipient_id    BIGINT                                                        NOT NULL REFERENCES public."user" (id) ON DELETE CASCADE,
    text            TEXT                                                          NOT NULL CHECK (text <> '')
    CONSTRAINT max_len_text CHECK (LENGTH(text) <= 4000),
    rating          INT                                                           NOT NULL
    CONSTRAINT rating_range CHECK (rating >= 1 AND rating <= 5),
    created_at      TIMESTAMP WITH TIME ZONE DEFAULT NOW()                        NOT NULL
);