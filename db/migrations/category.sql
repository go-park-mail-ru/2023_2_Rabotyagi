-- Table: public.category

drop table if exists public."category";
drop sequence if exists category_id_seq;

create sequence category_id_seq;
create table public."category"
(
    id bigint not null primary key default nextval('category_id_seq'::regclass),
    name character(256) unique not null,
    parent_id bigint default null
)