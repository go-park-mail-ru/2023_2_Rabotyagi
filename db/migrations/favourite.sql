-- Table: public.favourite

drop table if exists public."favourite" cascade;
drop sequence if exists favourite_id_seq;

create sequence favourite_id_seq;
create table public."favourite"
(
    id bigint not null primary key default nextval('favourite_id_seq'::regclass),
	owner_id bigint not null references public."user" (id),
    product_id bigint not null references product (id) on delete cascade
)