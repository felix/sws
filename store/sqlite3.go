package store

import (
	"fmt"
	"strings"
	"time"

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

func (s *Sqlite3) GetDomainByID(id int) (*sws.Domain, error) {
	var d sws.Domain
	if err := s.db.QueryRowx(stmts["domainByID"], id).StructScan(&d); err != nil {
		return nil, err
	}
	return &d, nil
}

func (s *Sqlite3) GetDomainByName(name string) (*sws.Domain, error) {
	var d sws.Domain
	// Ensure port is split off
	name = strings.Split(name, ":")[0]
	if err := s.db.QueryRowx(stmts["domainByName"], name).StructScan(&d); err != nil {
		return nil, err
	}
	return &d, nil
}

func (s *Sqlite3) SaveDomain(d *sws.Domain) error {
	if _, err := s.db.NamedExec(stmts["saveDomain"], d); err != nil {
		return err
	}
	return nil
}

func (s *Sqlite3) GetHits(d sws.Domain, start, end time.Time, f map[string]interface{}) ([]sws.Hit, error) {
	pvs := make([]sws.Hit, 0)

	filter := map[string]interface{}{
		"start": start,
		"end":   end,
	}

	sql := stmts["filterHits"]
	for k, v := range f {
		filter[k] = v
		sql += fmt.Sprintf("%s = :%s", k, k)
	}

	rows, err := s.db.NamedQuery(sql, filter)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		pv := &sws.Hit{}
		if err := rows.StructScan(pv); err != nil {
			return pvs, err
		}
		pvs = append(pvs, *pv)
	}
	return pvs, nil
}

func (s *Sqlite3) SaveHit(h *sws.Hit) error {
	if _, err := s.db.NamedExec(stmts["saveHit"], h); err != nil {
		return err
	}
	return nil
}

var stmts = map[string]string{
	"domainByName": `select id, name, description, aliases, enabled,
created_at, updated_at from domains where name = $1 limit 1`,

	"domainByID": `select id, name, description, aliases, enabled,
created_at, updated_at from domains where id = $1 limit 1`,

	"saveDomain": `insert into domains (
name, description, aliases, enabled, created_at, updated_at) values (:name,
:description, :aliases, :enabled, :created_at, :updated_at)`,

	"saveHit": `insert into hits (
domain_id, addr, scheme, host, path, query, title, referrer, user_agent,
view_port, created_at) values (:domain_id, :addr, :scheme, :host, :path, :query,
:title, :referrer, :user_agent, :view_port, :created_at)`,

	"filterHits": `select domain_id, addr, scheme, host, path, title,
referrer, user_agent, view_port, created_at from hits where created_at > :start
and created_at < :end`,
}
