package sws

type Store interface {
	DomainStore
	HitStore
	CounterStore
}
type SimpleDomainStore interface {
	GetDomainByName(string) (*Domain, error)
}
type DomainStore interface {
	//GetDomainFromID(int) (*Domain, error)
	//GetDomainFromCode(string) (*Domain, error)
	//SaveDomain(*Domain) error
}

type HitStore interface {
	//GetHitsByDomain(d Domain) ([]*Hit, error)
}
type CounterStore interface {
	SimpleDomainStore
	GetHits(map[string]interface{}) ([]*Hit, error)
	SaveHit(*Hit) error
}
