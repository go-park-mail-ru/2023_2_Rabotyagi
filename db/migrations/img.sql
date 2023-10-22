drop table if exists public."img";
drop sequence if exists img_id_seq;

create sequence img_id_seq;
create table public."img"
(
    id         bigint         not null primary key default nextval('img_id_seq'::regclass),
    url        character(256) not null,
    product_id bigint         not null -- foreign key
)