create table users (
	id integer primary key autoincrement,
	email varchar not null,
	first_name varchar not null,
	last_name varchar not null,
	enabled integer not null default 0,
	pw_hash varchar not null,
	pw_salt varchar not null,
	last_login_at timestamp null,
	created_at timestamp not null,
	updated_at timestamp not null
);
create index "users#email" on users(email);

insert into users (email, first_name, last_name, enabled, pw_hash, pw_salt, created_at, updated_at)
values ('admin@example.com', 'Admin', 'User', 1, 'OD+kBtFSc+HUMSzCsJL/HL7TjOViMTZ3jMDuhUil/ys', 'UvQgb6w0RjCOaM9L', date('now'), date('now'));
