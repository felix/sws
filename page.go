package sws

import (
	"time"
)

type Page struct {
	SiteID        int       `json:"site_id"`
	Path          string    `json:"path"`
	Title         *string   `json:"title,omitempty"`
	LastVisitedAt time.Time `json:"last_visited_at"`
	hitSet        *HitSet
	// TODO
	Site *Site `json:"-"`
}

func (p Page) YMax() int {
	return p.hitSet.YMax()
}
func (p Page) XSeries() []*bucket {
	//p.hitSet.Fill(nil, nil)
	//p.hitSet.SortByDate()
	return p.hitSet.XSeries()
}

func (p Page) Count() int {
	return p.hitSet.Count()
}

func (p Page) Label() string {
	return p.Path
}
