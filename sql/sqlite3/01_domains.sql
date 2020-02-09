create table domains (
	id integer primary key autoincrement,
	name varchar not null check(length(name) >= 4 and length(name) <= 255),
	description varchar null,
	enabled integer not null default 0,
	created_at timestamp not null check(created_at = strftime('%Y-%m-%d %H:%M:%S', created_at)),
	updated_at timestamp not null check(updated_at = strftime('%Y-%m-%d %H:%M:%S', updated_at))
);
