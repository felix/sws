package sws

import (
	"fmt"
)

func GetHits(db Queryer, f map[string]interface{}) ([]*Hit, error) {
	pvs := make([]*Hit, 0)
	qa := queryArgs{}

	sql := sqlFilterHits
	if len(f) > 0 {
		sql = sql + " where "
		for k, v := range f {
			sql += fmt.Sprintf(" %s = %s", k, qa.Append(v))
		}
	}

	rows, err := db.Query(sql, qa...)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		pv := &Hit{}
		if err := rows.Scan(pv); err != nil {
			return pvs, err
		}
		pvs = append(pvs, pv)
	}
	return pvs, nil
}

func (s *Hit) Save(db Queryer) error {
	if _, err := db.Exec(
		sqlSaveHit,
		s.DomainID,
		s.Addr,
		s.Scheme,
		s.Host,
		s.Path,
		s.Query,
		s.Title,
		s.Referrer,
		s.UserAgent,
		s.ViewPort,
		s.CreatedAt,
	); err != nil {
		return err
	}
	return nil
}

const (
	sqlSaveHit = `insert into hits (
domain_id,
addr,
scheme,
host,
path,
query,
title,
referrer,
user_agent,
view_port,
created_at
) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
`
	sqlFilterHits = `select
domain_id,
address,
scheme,
host,
page,
title,
referrer,
user_agent,
view_port,
created_at
from hits
`
)
