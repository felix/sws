package sws

import "time"

type Store interface {
	SiteStore
	GetSiteByName(string) (*Site, error)
	GetHits(Site, time.Time, time.Time, map[string]interface{}) ([]*Hit, error)
	SaveHit(*Hit) error
}
type SimpleSiteStore interface {
	GetSiteByName(string) (*Site, error)
}
type SiteStore interface {
	GetSites() ([]*Site, error)
	GetSiteByID(int) (*Site, error)
	GetPages(Site, time.Time, time.Time) ([]*Page, error)
	//SaveSite(*Site) error
}

type HitStore interface {
	SimpleSiteStore
	GetHits(Site, time.Time, time.Time, map[string]interface{}) ([]*Hit, error)
}
type CounterStore interface {
	SimpleSiteStore
	SaveHit(*Hit) error
}
