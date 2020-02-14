create table sites (
	id integer primary key autoincrement,
	name varchar not null check(length(name) >= 4 and length(name) <= 255),
	description varchar null,
	aliases varchar null,
	enabled integer not null default 0,
	created_at timestamp not null,
	updated_at timestamp not null
);

insert into sites (name, description, enabled, created_at, updated_at)
values ('localhost', 'Example site', 1, date('now'), date('now'));
