package sws

type Store interface {
	SiteStore
	UserStore
	HitSaver
	HitCursor(func(*Hit) error) error
}
type HitSaver interface {
	SaveHit(*Hit) error
}
type SiteStore interface {
	SiteGetter
	GetHits(Site, map[string]interface{}) ([]*Hit, error)
	GetSites() ([]*Site, error)
	SaveSite(*Site) error
}
type SiteGetter interface {
	GetSiteByID(int) (*Site, error)
}
type HitStore interface {
	HitSaver
	GetHits(Site, map[string]interface{}) ([]*Hit, error)
	HitCursor(func(*Hit) error) error
}
type CounterStore interface {
	SiteGetter
	HitSaver
}
type UserStore interface {
	GetUserByID(int) (*User, error)
	GetUserByEmail(string) (*User, error)
	SaveUser(*User) error
}
