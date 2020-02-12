package sws

import (
	"fmt"
	"time"
)

const slugSalt = "saltyslugs"

type Domain struct {
	ID          *int       `json:"id,omitempty"`
	Name        *string    `json:"name,omitempty"`
	Description string     `json:"description,omitempty"`
	Aliases     *string    `json:"aliases,omitempty"`
	Enabled     bool       `json:"enabled"`
	CreatedAt   *time.Time `json:"created_at,omitempty"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
}

func (d *Domain) Validate() []error {
	var out []error
	if d.Name == nil {
		out = append(out, fmt.Errorf("missing name"))
	}
	return out
}
