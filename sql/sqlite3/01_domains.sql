create table domains (
	id integer primary key autoincrement,
	name varchar not null check(length(name) >= 4 and length(name) <= 255),
	description varchar null,
	aliases varchar null,
	enabled integer not null default 0,
	created_at timestamp not null,
	updated_at timestamp not null
);

insert into domains (name, description, enabled, created_at, updated_at)
values ('localhost', 'Example domain', 1, date('now'), date('now'));
