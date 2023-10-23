drop table if exists public."image";
drop sequence if exists image_id_seq;

create sequence image_id_seq;
create table public."image"
(
    id         bigint default nextval('image_id_seq'::regclass) not null primary key,
    url        character(256)                                   not null unique,
    product_id bigint                                           not null references public."product" (id) on delete cascade
);