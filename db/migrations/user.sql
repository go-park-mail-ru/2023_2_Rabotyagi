-- Table: public.user

drop table if exists public."user" cascade;
drop sequence if exists user_id_seq;

create sequence user_id_seq;
create table public."user"
(
    id bigint not null primary key default nextval('user_id_seq'::regclass),
    email character(256) unique not null,
    phone character(18) unique not null,
    name character(256),
    pass character(256) not null,
    birthday timestamp with time zone
);