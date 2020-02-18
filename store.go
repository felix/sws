package sws

type Store interface {
	SiteStore
	GetSiteByName(string) (*Site, error)
	SaveHit(*Hit) error
}
type SimpleSiteStore interface {
	GetSiteByName(string) (*Site, error)
}
type SiteStore interface {
	GetSites() ([]*Site, error)
	GetSiteByID(int) (*Site, error)
	GetPages(Site, map[string]interface{}) ([]*Page, error)
	GetHits(Site, map[string]interface{}) ([]*Hit, error)
	//SaveSite(*Site) error
}

type HitStore interface {
	SimpleSiteStore
	GetHits(Site, map[string]interface{}) ([]*Hit, error)
}
type CounterStore interface {
	SimpleSiteStore
	SaveHit(*Hit) error
}
