package store

import (
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/reflectx"
	"src.userspace.com.au/sws"
)

type Sqlite3 struct {
	db *sqlx.DB
}

func NewSqlite3Store(db *sqlx.DB) *Sqlite3 {
	db.Mapper = reflectx.NewMapperFunc("json", strings.ToLower)
	return &Sqlite3{db}
}

func (s *Sqlite3) GetSites() ([]*sws.Site, error) {
	rows, err := s.db.Queryx(stmts["sites"])
	if err != nil {
		return nil, err
	}
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
	pvs := make([]*sws.Hit, 0)

	sql := stmts["hits"]
	filter["site_id"] = *d.ID
	processFilter(&sql, filter)

	rows, err := s.db.NamedQuery(sql, filter)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		pv := &sws.Hit{}
		if err := rows.StructScan(pv); err != nil {
			return pvs, err
		}
		pvs = append(pvs, pv)
	}
	return pvs, nil
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

func (s *Sqlite3) GetPages(d sws.Site, filter map[string]interface{}) ([]*sws.Page, error) {
	pages := make([]*sws.Page, 0)

	sql := stmts["pages"]
	filter["h.site_id"] = *d.ID
	for k, _ := range filter {
		sql += " and"
		switch k {
		case "begin":
			sql += fmt.Sprintf(" l.created_at > :%s", k)
		case "end":
			sql += fmt.Sprintf(" l.created_at < :%s", k)
		default:
			sql += fmt.Sprintf(" %s = :%s", k, k)
		}
	}

	rows, err := s.db.NamedQuery(sql, filter)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		p := &sws.Page{}
		if err := rows.StructScan(p); err != nil {
			return pages, err
		}
		pages = append(pages, p)
	}
	return pages, nil
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
	"sites": `select id, name, description, aliases, enabled,
created_at, updated_at
from sites`,

	"siteByName": `select id, name, description, aliases, enabled,
created_at, updated_at
from sites
where name = $1 limit 1`,

	"siteByID": `select id, name, description, aliases, enabled,
created_at, updated_at
from sites
where id = $1 limit 1`,

	"saveSite": `insert into sites (
name, description, aliases, enabled, created_at, updated_at)
values (:name, :description, :aliases, :enabled, :created_at, :updated_at)`,

	"userAgentByHash": `select id, hash, name, last_seen_at from sites
where hash = $1 limit 1`,

	"saveUserAgent": `insert into user_agents
(hash, name, last_seen_at)
values (:hash, :name, :last_seen_at)
on conflict(hash) do update set last_seen_at = :last_seen_at`,

	"saveHit": `insert into hits (
site_id, addr, scheme, host, path, query, title, referrer, user_agent_hash,
view_port, no_script, created_at)
values (:site_id, :addr, :scheme, :host, :path, :query, :title, :referrer,
:user_agent_hash, :view_port, :no_script, :created_at)`,

	"pages": `with latest as (select site_id, path, max(created_at) as created_at
from hits group by site_id, path)
select h.site_id, h.path, h.created_at as last_visited_at
from hits h, latest l
where l.site_id = h.site_id and l.path = h.path and h.created_at = l.created_at`,

	"hits": `select site_id, addr, scheme, host, path, title,
referrer, user_agent_hash, view_port, no_script, created_at
from hits
where site_id = :site_id`,
}
