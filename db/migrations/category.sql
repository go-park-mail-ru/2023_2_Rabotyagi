-- Table: public.category

drop table if exists public."category" cascade;
drop sequence if exists category_id_seq;

create sequence category_id_seq;
create table public."category"
(
    id        bigint default nextval('category_id_seq'::regclass) not null primary key,
    name      character(256) unique                               not null,
    parent_id bigint default null references public.category (id)
);