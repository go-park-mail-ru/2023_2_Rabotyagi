-- Table: public.order

drop table if exists public."order" cascade;
drop sequence if exists order_id_seq;

create sequence order_id_seq;
create table public."order"
(
    id          bigint                   default nextval('order_id_seq'::regclass) not null primary key,
    owner_id    bigint                                                             not null references public."user" (id),
    product_id  bigint                                                             not null references public."product" (id),
    count       smallint                                                           not null default 1 check (count > 0),
    status      smallint                                                           not null default 0,
    create_date timestamp with time zone default now()                             not null,
    update_date timestamp with time zone default now()                             not null,
    close_date  timestamp with time zone default null
);

create or replace function update_data()
    returns trigger as
$$
begin
    new.update_date = now();
    return new;
end;
$$ language plpgsql;

create trigger verify_update_data
    before update
    on public."order"
    for each row
execute procedure update_data();