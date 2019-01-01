package models

import (
	"fmt"
	"net/url"
	"time"
)

type PageView struct {
	ID        *int              `json:"id"`
	Address   *string           `json:"address,omitempty"`
	Scheme    *string           `json:"scheme,omitempty"`
	Host      *string           `json:"host,omitempty"`
	Page      *string           `json:"page,omitempty"`
	Title     *string           `json:"title,omitempty"`
	Referrer  *string           `json:"referrer,omitempty"`
	UserAgent *string           `json:"user_agent,omitempty"`
	ViewPort  *string           `json:"view_port,omitempty"`
	Features  map[string]string `json:"features,omitempty"`
	CreatedAt *time.Time        `json:"created_at,omitempty"`
	DomainID  *int              `json:"domain_id,omitempty"`

	// TODO
	Domain *Domain `json:"-"`
}

func ptrString(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func ptrTime(t time.Time) *time.Time {
	return &t
}

func PageViewFromURL(u url.URL) (*PageView, error) {
	q := u.Query()
	host := q.Get("h")
	if host == "" {
		return nil, fmt.Errorf("missing host")
	}
	pv := PageView{
		Scheme:    ptrString(q.Get("s")),
		Host:      &host,
		Page:      ptrString(q.Get("p")),
		Title:     ptrString(q.Get("t")),
		Referrer:  ptrString(q.Get("r")),
		UserAgent: ptrString(q.Get("u")),
		ViewPort:  ptrString(q.Get("v")),
		CreatedAt: ptrTime(time.Now()),
	}
	return &pv, nil
}

func GetPageViews(f map[string]interface{}) ([]*PageView, error) {
	pvs := make([]*PageView, 0)
	qa := new(queryArgs)

	sql := sqlFilterPageViews
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
		pv := &PageView{}
		if err := rows.Scan(pv); err != nil {
			return pvs, err
		}
		pvs = append(pvs, pv)
	}
	return pvs, nil
}

func (s *PageView) Save(db Queryer) error {
	var id int
	if err := db.QueryRow(
		sqlSavePageView,
		s.DomainID,
		s.Address,
		s.Scheme,
		s.Host,
		s.Page,
		s.Title,
		s.Referrer,
		s.UserAgent,
		s.ViewPort,
		s.CreatedAt,
	).Scan(&id); err != nil {
		return err
	}
	s.ID = &id
	return nil
}

const (
	sqlSavePageView = `insert into page_views (
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
) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
returning id
`
	sqlFilterPageViews = `select
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
from page_views
`
)
