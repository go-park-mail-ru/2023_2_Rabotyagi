-- Table: public.favourite

drop table if exists public."favourite";
drop sequence if exists favourite_id_seq;

create sequence favourite_id_seq;
create table public."favourite"
(
    id bigint not null primary key default nextval('favourite_id_seq'::regclass),
    product_id bigint not null,
	owner_id bigint not null
)