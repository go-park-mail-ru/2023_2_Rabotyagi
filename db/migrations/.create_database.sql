drop table if exists public."user" cascade;
drop table if exists public."product" cascade;
drop table if exists public."category" cascade;
drop table if exists public."order" cascade;
drop table if exists public."image" cascade;
drop table if exists public."favourite" cascade;
drop sequence if exists user_id_seq;
drop sequence if exists category_id_seq;
drop sequence if exists product_id_seq;
drop sequence if exists order_id_seq;
drop sequence if exists image_id_seq;
drop sequence if exists favourite_id_seq;

create sequence user_id_seq;
create sequence category_id_seq;
create sequence product_id_seq;
create sequence order_id_seq;
create sequence image_id_seq;
create sequence favourite_id_seq;

create table public."user"
(
    id       bigint                not null primary key default nextval('user_id_seq'::regclass),
    email    character(256) unique not null,
    phone    character(18) unique  not null,
    name     character(256),
    pass     character(256)        not null,
    birthday timestamp with time zone
);

create table public."category"
(
    id        bigint                not null primary key default nextval('category_id_seq'::regclass),
    name      character(256) unique not null,
    parent_id bigint                                     default null references public.category (id)
);

create table public."product"
(
    id              bigint         not null primary key default nextval('product_id_seq'::regclass),
    saler_id        bigint         not null references public."user" (id),
    category_id     bigint         not null references category (id),
    title           character(256) not null,
    description     text           not null,
    price           bigint                              default 0 check (price >= 0),
    create_date     timestamp with time zone            default now(),
    views           int                                 default 0 check (views >= 0),
    in_favourites   int                                 default 0 check (in_favourites >= 0),
    available_count int                                 default 0 check (available_count >= 0),
    city            character(256) not null,
    delivery        boolean                             default false,
    safe_dial       boolean                             default false,
    is_active       boolean                             default false,
    constraint not_null_good_count check (not (available_count = 0 and is_active))
);

create table public."order"
(
    id          bigint                   not null primary key default nextval('order_id_seq'::regclass),
    owner_id    bigint                   not null references public."user" (id),
    product_id  bigint                   not null references product (id),
    count       smallint                 not null             default 1 check (count > 0),
    status      smallint                 not null             default 0,
    create_date timestamp with time zone not null             default now(),
    update_date timestamp with time zone not null             default now(),
    close_date  timestamp with time zone                      default null
);

create table public."image"
(
    id         bigint         not null primary key default nextval('image_id_seq'::regclass),
    url        character(256) not null unique,
    product_id bigint         not null references product (id) on delete cascade
);

create table public."favourite"
(
    id         bigint not null primary key default nextval('favourite_id_seq'::regclass),
    owner_id   bigint not null references public."user" (id),
    product_id bigint not null references product (id) on delete cascade
);

create or replace function update_data_func()
    returns trigger as $$
begin
    new.update_data = now();
    return new;
end;
$$ language plpgsql;

create trigger verify_update_data
    before update
    on public."order"
    for each row
execute procedure update_data_func();