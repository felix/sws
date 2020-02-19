package sws

import (
	"time"
)

type Page struct {
	SiteID        int       `json:"site_id"`
	Path          string    `json:"path"`
	Title         *string   `json:"title,omitempty"`
	Count         int       `json:"count"`
	LastVisitedAt time.Time `json:"last_visited_at"`

	// TODO
	Site *Site `json:"-"`
}

func PagesFromHits(hits []*Hit) map[string]*Page {
	out := make(map[string]*Page)
	for _, h := range hits {
		p, ok := out[h.Path]
		if !ok {
			p = &Page{
				Path:          h.Path,
				SiteID:        *h.SiteID,
				Title:         h.Title,
				LastVisitedAt: h.CreatedAt,
			}
		}
		if p.LastVisitedAt.Before(h.CreatedAt) {
			p.LastVisitedAt = h.CreatedAt
		}
		p.Count++
		out[h.Path] = p
	}
	return out
}
