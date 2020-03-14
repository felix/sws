package sws

type Store interface {
	SiteStore
	UserStore
	GetSiteByName(string) (*Site, error)
	HitSaver
}
type HitSaver interface {
	SaveHit(*Hit) error
}
type SimpleSiteStore interface {
	GetSiteByName(string) (*Site, error)
}
type SiteStore interface {
	GetSites() ([]*Site, error)
	GetSiteByID(int) (*Site, error)
	GetHits(Site, map[string]interface{}) ([]*Hit, error)
	SaveSite(*Site) error
}
type HitStore interface {
	SimpleSiteStore
	GetHits(Site, map[string]interface{}) ([]*Hit, error)
}
type CounterStore interface {
	SimpleSiteStore
	HitSaver
}
type UserStore interface {
	GetUserByID(int) (*User, error)
	GetUserByEmail(string) (*User, error)
	SaveUser(*User) error
}
