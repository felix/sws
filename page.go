package sws

import (
	"time"
)

type Page struct {
	ID     int     `json:"id"`
	SiteID int     `json:"site_id"`
	Path   string  `json:"path"`
	Title  *string `json:"title,omitempty"`

	LastVisitedAt time.Time `json:"last_visited_at"`

	// TODO
	Site *Site `json:"-"`
}
