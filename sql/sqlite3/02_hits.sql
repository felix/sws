create table hits (
	id integer primary key autoincrement,
	domain_id integer check(domain_id >0),
	addr varchar not null,
	scheme varchar not null,
	host varchar not null,
	path varchar not null,
	query varchar null,
	title varchar null,
	referrer varchar null,
	user_agent varchar null,
	view_port varchar null,
	created_at timestamp not null
);
create index "hits#domain_id#created" on hits(domain_id, created_at);
