-- create extension if not exists hstore;

-- set role felix;

create table if not exists domains (
    id serial primary key,
    name varchar(255) not null,
    created_at timestamp without time zone not null,
    updated_at timestamp without time zone not null
);
create unique index if not exists domains_name_idx on domains (name);

create table if not exists page_views (
    id bigserial primary key,
    domain_id integer not null references domains (id),
    address varchar(50),
    scheme varchar(20) not null,
    host varchar(255) not null,
    page varchar(1000),
    title varchar(255),
    referrer varchar(2000),
    user_agent varchar(255),
    view_port varchar(20),
    attributes hstore,
    created_at timestamp without time zone not null
);
