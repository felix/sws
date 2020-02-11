package sws

import (
	"fmt"
	"time"

	"github.com/speps/go-hashids"
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

func (d Domain) Slug() string {
	hd := hashids.NewData()
	hd.Salt = slugSalt
	h, _ := hashids.NewWithData(hd)
	out, _ := h.Encode([]int{*d.ID})
	return out
}
