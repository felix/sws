package sws

import (
	"time"
)

const slugSalt = "saltyslugs"

type Site struct {
	ID               *int       `json:"id,omitempty"`
	Name             string     `json:"name,omitempty"`
	Description      string     `json:"description,omitempty"`
	AcceptSubdomains bool       `json:"subdomains"`
	Aliases          string     `json:"aliases,omitempty"`
	IgnoreIPs        string     `json:"ignore_ips"`
	Enabled          bool       `json:"enabled"`
	CreatedAt        *time.Time `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt        *time.Time `json:"updated_at,omitempty" db:"updated_at"`
}

func (d *Site) Validate() []string {
	var out []string
	if d.Name == "" {
		out = append(out, "missing name")
	}
	return out
}
