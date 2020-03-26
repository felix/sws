alter table user_agents add column browser varchar not null default '';
alter table user_agents add column platform varchar not null default '';
alter table user_agents add column version varchar not null default '';
alter table user_agents add column bot integer not null default 0;
alter table user_agents add column mobile integer not null default 0;
