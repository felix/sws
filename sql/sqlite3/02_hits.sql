pragma foreign_keys = on;

create table user_agents (
	hash varchar not null primary key,
	name varchar not null,
	last_seen_at timestamp not null
);

create table hits (
	id integer primary key autoincrement,
	site_id integer not null,
	addr varchar not null,
	scheme varchar not null,
	host varchar not null,
	path varchar not null,
	query varchar null,
	title varchar null,
	referrer varchar null,
	user_agent_hash varchar null,
	view_port varchar null,
	no_script integer not null default 0,
	created_at timestamp not null,
	foreign key(user_agent_hash) references user_agents(hash)
);
create index "hits#site_id#created" on hits(site_id, created_at);
create index "hits#site_id#path" on hits(site_id, path);
