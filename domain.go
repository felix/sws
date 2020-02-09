package sws

import (
	"time"

	"github.com/speps/go-hashids"
)

const slugSalt = "saltyslugs"

type Domain struct {
	ID          *int       `json:"id"`
	Name        *string    `json:"name"`
	Description *string    `json:"description"`
	Enabled     bool       `json:"enabled"`
	CreatedAt   *time.Time `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
}

func (d Domain) Slug() string {
	hd := hashids.NewData()
	hd.Salt = slugSalt
	h, _ := hashids.NewWithData(hd)
	out, _ := h.Encode([]int{*d.ID})
	return out
}
