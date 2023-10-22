-- Table: public.image

drop table if exists public."image" cascade;
drop sequence if exists image_id_seq;

create sequence image_id_seq;
create table public."image"
(
    id         bigint         not null primary key default nextval('image_id_seq'::regclass),
    url        character(256) not null unique,
    product_id bigint         not null references product (id) on delete cascade
)