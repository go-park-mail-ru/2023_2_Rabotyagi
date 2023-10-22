drop table if exists public."order";
drop sequence if exists order_id_seq;

create sequence order_id_seq;
create table public."order"
(
    id          bigint                   not null primary key default nextval('order_id_seq'::regclass),
    owner_id    bigint                   not null,                              -- foreign key
    product_id  bigint                   not null,                              -- foreign key
    count       smallint                 not null             default 1,
    status      smallint                 not null             default 0,
    create_date timestamp with time zone not null             default time.now,
    update_date timestamp with time zone not null             default time.now, -- triger on update
    close_date  timestamp with time zone
)