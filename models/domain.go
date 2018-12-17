package models

import (
	"strings"
	"time"
)

type Domain struct {
	ID        *int
	Name      *string
	CreatedAt *time.Time
	UpdatedAt *time.Time
}

func GetDomainByName(db Queryer, name string) (*Domain, error) {
	d := Domain{}
	name = strings.Split(name, ":")[0]
	if err := db.QueryRow(sqlDomainByName, name).Scan(
		&d.ID,
		&d.Name,
		&d.CreatedAt,
		&d.UpdatedAt,
	); err != nil {
		return nil, err
	}
	return &d, nil
}

const (
	sqlDomainByName = `select
id, name, created_at, updated_at
from domains
where $1 = name
or right($1, length(name)) = name
limit 1`
)
