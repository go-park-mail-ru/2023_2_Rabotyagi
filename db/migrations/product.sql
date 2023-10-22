-- Table: public.product
drop table if exists public."product" cascade;
drop sequence if exists product_id_seq;

create sequence product_id_seq;
create table public."product"
(
  id bigint not null primary key default nextval('product_id_seq'::regclass),
	saler_id bigint not null references public."user" (id),
	category_id bigint not null references category (id),
	title character(256) not null,
	description text not null,
	price bigint default 0 check (price >= 0),
	create_date timestamp with time zone default now(),
	views int default 0 check (views >= 0),
	in_favourites int default 0 check (in_favourites >= 0),
	available_count int default 0 check (available_count >= 0),
	city character(256) not null,
	delivery boolean default false,
	safe_dial boolean default false,
	is_active boolean default false,
	constraint not_null_good_count check (not (available_count = 0 and is_active))
);