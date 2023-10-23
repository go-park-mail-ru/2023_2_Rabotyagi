-- Table: public.favourite

drop table if exists public."favourite" cascade;
drop sequence if exists favourite_id_seq;

create sequence favourite_id_seq;
create table public."favourite"
(
    id         bigint default nextval('favourite_id_seq'::regclass) not null primary key,
    owner_id   bigint                                               not null references public."user" (id),
    product_id bigint                                               not null references public."product" (id) on delete cascade
);