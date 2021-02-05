alter table sites add column subdomains integer not null default 0;
alter table sites add column ignore_ips varchar not null;
