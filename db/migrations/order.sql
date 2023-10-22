-- Table: public.order

drop table if exists public."order" cascade;
drop sequence if exists order_id_seq;

create sequence order_id_seq;
create table public."order"
(
    id          bigint not null primary key default nextval('order_id_seq'::regclass),
    owner_id    bigint not null references public."user" (id),
    product_id  bigint not null references product (id),
    count       smallint not null default 1,
    status      smallint not null default 0,
    create_date timestamp with time zone not null default now(),
    update_date timestamp with time zone not null default now(),
    close_date  timestamp with time zone default null
);