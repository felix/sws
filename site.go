package sws

import (
	"fmt"
	"time"
)

const slugSalt = "saltyslugs"

type Site struct {
	ID          *int    `json:"id,omitempty"`
	Name        *string `json:"name,omitempty"`
	Description string  `json:"description,omitempty"`
	Aliases     *string `json:"aliases,omitempty"`
	Enabled     bool    `json:"enabled"`
	//ExcludePaths []string
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}

func (d *Site) Validate() []error {
	var out []error
	if d.Name == nil {
		out = append(out, fmt.Errorf("missing name"))
	}
	return out
}
