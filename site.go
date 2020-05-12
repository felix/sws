package sws

import (
	"strings"
	"time"
)

const slugSalt = "saltyslugs"

type Site struct {
	ID               *int       `json:"id,omitempty"`
	Name             string     `json:"name,omitempty"`
	Description      string     `json:"description,omitempty"`
	AcceptSubdomains bool       `json:"subdomains" db:"subdomains"`
	Aliases          string     `json:"aliases,omitempty"`
	IgnoreIPs        string     `json:"ignore_ips" db:"ignore_ips"`
	Enabled          bool       `json:"enabled"`
	CreatedAt        *time.Time `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt        *time.Time `json:"updated_at,omitempty" db:"updated_at"`
}

func (s *Site) Validate() []string {
	var out []string
	if s.Name == "" {
		out = append(out, "missing name")
	}
	return out
}

func (s *Site) IncludesDomain(fqdn string) bool {
	if fqdn == s.Name {
		return true
	}
	for _, a := range strings.Split(s.Aliases, ",") {
		if a == fqdn {
			return true
		}
	}
	if s.AcceptSubdomains && strings.HasSuffix(fqdn, s.Name) {
		return true
	}
	return false
}
