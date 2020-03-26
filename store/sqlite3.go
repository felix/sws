package store

import (
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	//"github.com/jmoiron/sqlx/reflectx"
	"src.userspace.com.au/sws"
)

type Sqlite3 struct {
	db *sqlx.DB
}

func NewSqlite3Store(db *sqlx.DB) *Sqlite3 {
	//db.Mapper = reflectx.NewMapperFunc("json", strings.ToLower)
	out := &Sqlite3{
		db: db,
	}
	return out
}

func (s *Sqlite3) GetSites() ([]*sws.Site, error) {
	rows, err := s.db.Queryx(stmts["sites"])
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []*sws.Site

	for rows.Next() {
		var s sws.Site
		if err := rows.StructScan(&s); err != nil {
			return nil, err
		}
		out = append(out, &s)
	}
	return out, nil
}

func (s *Sqlite3) GetSiteByID(id int) (*sws.Site, error) {
	var d sws.Site
	if err := s.db.QueryRowx(stmts["siteByID"], id).StructScan(&d); err != nil {
		return nil, err
	}
	return &d, nil
}

func (s *Sqlite3) GetSiteByName(name string) (*sws.Site, error) {
	var d sws.Site
	// Ensure port is split off
	name = strings.Split(name, ":")[0]
	if err := s.db.QueryRowx(stmts["siteByName"], name).StructScan(&d); err != nil {
		return nil, err
	}
	return &d, nil
}

func (s *Sqlite3) SaveSite(d *sws.Site) error {
	if _, err := s.db.NamedExec(stmts["saveSite"], d); err != nil {
		return err
	}
	return nil
}

func (s *Sqlite3) GetHits(d sws.Site, filter map[string]interface{}) ([]*sws.Hit, error) {
	hits := make([]*sws.Hit, 0)

	if filter == nil {
		filter = make(map[string]interface{})
	}

	sql := stmts["hits"] + ` where h.site_id = :site_id`
	filter["site_id"] = *d.ID
	processFilter(&sql, filter)

	rows, err := s.db.NamedQuery(sql, filter)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		h := &sws.Hit{}
		if err := rows.StructScan(h); err != nil {
			return hits, err
		}
		hits = append(hits, h)
	}
	return hits, nil
}

func (s *Sqlite3) HitCursor(f func(h *sws.Hit) error) error {
	sql := stmts["hits"]

	rows, err := s.db.Queryx(sql)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		h := &sws.Hit{}
		if err := rows.StructScan(h); err != nil {
			return err
		}
		if err := f(h); err != nil {
			return err
		}
	}
	return nil
}

func (s *Sqlite3) SaveHit(h *sws.Hit) error {
	if h.UserAgent != nil {
		if _, err := s.db.NamedExec(stmts["saveUserAgent"], *h.UserAgent); err != nil {
			return err
		}
	}
	if _, err := s.db.NamedExec(stmts["saveHit"], h); err != nil {
		return err
	}
	return nil
}

func (s *Sqlite3) GetUserByEmail(email string) (*sws.User, error) {
	var u sws.User
	if err := s.db.QueryRowx(stmts["userByEmail"], email).StructScan(&u); err != nil {
		return nil, err
	}
	return &u, nil
}

func (s *Sqlite3) GetUserByID(id int) (*sws.User, error) {
	var u sws.User
	if err := s.db.QueryRowx(stmts["userByID"], id).StructScan(&u); err != nil {
		return nil, err
	}
	return &u, nil
}

func (s *Sqlite3) SaveUser(u *sws.User) error {
	if _, err := s.db.NamedExec(stmts["saveUser"], u); err != nil {
		return err
	}
	return nil
}

func processFilter(sql *string, filter map[string]interface{}) {
	if sql == nil {
		panic("empty sql")
	}
	for k, _ := range filter {
		*sql += " and"
		switch k {
		case "begin":
			*sql += fmt.Sprintf(" created_at > :%s", k)
		case "end":
			*sql += fmt.Sprintf(" created_at < :%s", k)
		default:
			*sql += fmt.Sprintf(" %s = :%s", k, k)
		}
	}
}

var stmts = map[string]string{
	"sites": `select * from sites`,

	"siteByName": `select * from sites where name = $1 limit 1`,

	"siteByID": `select * from sites where id = $1 limit 1`,

	"saveSite": `insert into sites
(id, name, description, aliases, enabled, subdomains, ignore_ips, created_at, updated_at)
values (:id, :name, :description, :aliases, :enabled, :subdomains, :ignore_ips, date('now'), date('now'))
on conflict(id) do update set
name = :name,
description = :description,
aliases = :aliases,
subdomains = :subdomains,
updated_at = date('now')`,

	"userAgentByHash": `select id, hash, name, last_seen_at from sites
where hash = $1 limit 1`,

	"saveUserAgent": `insert into user_agents
(hash, name, last_seen_at, browser, platform, version, bot, mobile)
values (:hash, :name, :last_seen_at, :browser, :platform, :version, :bot, :mobile)
on conflict(hash) do update set
last_seen_at = :last_seen_at,
browser = :browser,
platform = :platform,
version = :version,
bot = :bot,
mobile = :mobile`,

	"saveHit": `insert into hits
(id, site_id, addr, scheme, host, path, query, title, referrer, user_agent_hash,
view_port, country_code, no_script, created_at)
values (:id, :site_id, :addr, :scheme, :host, :path, :query, :title, :referrer,
:user_agent_hash, :view_port, :country_code, :no_script, :created_at)
on conflict(id) do nothing`,

	// The explicit useragent stuff is to work around sqlx filling nested structs
	"hits": `select h.*,
ua.hash as "ua.hash", ua.name as "ua.name", ua.last_seen_at as "ua.last_seen_at",
ua.browser as "ua.browser", ua.platform as "ua.platform", ua.bot as "ua.bot",
ua.mobile as "ua.mobile"
from hits h
join user_agents ua on h.user_agent_hash = ua.hash`,

	"userByEmail": `select id, email, first_name, last_name, pw_hash, pw_salt, enabled,
created_at, updated_at, last_login_at
from users
where email = $1`,

	"userByID": `select id, email, first_name, last_name, pw_hash, pw_salt, enabled,
created_at, updated_at, last_login_at
from users
where id = $1`,

	"saveUser": `insert into users
(id, first_name, last_name, email, pw_hash, pw_salt, created_at, updated_at)
values (:id, :first_name, :last_name, :email, :pw_hash, :pw_salt, date('now'), date('now'))
on conflict(id) do update set
first_name = :first_name,
last_name = :last_name,
email = :email,
pw_hash = :pw_hash,
pw_salt = :pw_salt,
updated_at = date('now')`,
}
