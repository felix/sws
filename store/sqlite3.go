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

func (s *Sqlite3) GetDomainByName(name string) (*sws.Domain, error) {
	var d sws.Domain
	// Ensure port is split off
	name = strings.Split(name, ":")[0]
	if err := s.db.QueryRowx(stmts["domainByName"], name).StructScan(&d); err != nil {
		return nil, err
	}
	return &d, nil
}

func (s *Sqlite3) GetHits(f map[string]interface{}) ([]*sws.Hit, error) {
	pvs := make([]*sws.Hit, 0)
	qa := queryArgs{}

	sql := stmts["filterHits"]
	if len(f) > 0 {
		sql = sql + " where "
		for k, v := range f {
			sql += fmt.Sprintf(" %s = %s", k, qa.Append(v))
		}
	}

	rows, err := s.db.Query(sql, qa...)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		pv := &sws.Hit{}
		if err := rows.Scan(pv); err != nil {
			return pvs, err
		}
		pvs = append(pvs, pv)
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

	"saveHit": `insert into hits (
domain_id, addr, scheme, host, path, query, title, referrer, user_agent,
view_port, created_at) values (:domain_id, :addr, :scheme, :host, :path, :query,
:title, :referrer, :user_agent, :view_port, :created_at)`,

	"filterHits": `select domain_id, address, scheme, host, page, title,
referrer, user_agent, view_port, created_at from hits`,
}
