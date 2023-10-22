-- Table: public.product
drop table if exists public."product";
drop sequence if exists product_id_seq;

create sequence product_id_seq;
create table public."product"
(
  id bigint not null primary key default nextval('product_id_seq'::regclass),
	saler_id bigint not null,
	category_id bigint not null,
	title character(256) not null,
	description text not null,
	price bigint default 0,
	create_date timestamp with time zone default now(),
	views int default 0,
	in_favourites int default 0,
	available_count int default 0,
	city character(256) not null,
	delivery boolean default false,
	safe_dial boolean default false,
	is_active boolean default false
)