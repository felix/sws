package sws

import "time"

type Store interface {
	DomainStore
	GetDomainByName(string) (*Domain, error)
	GetHits(Domain, time.Time, time.Time, map[string]interface{}) ([]*Hit, error)
	SaveHit(*Hit) error
}
type SimpleDomainStore interface {
	GetDomainByName(string) (*Domain, error)
}
type DomainStore interface {
	//GetDomainFromID(int) (*Domain, error)
	GetDomainByID(int) (*Domain, error)
	//SaveDomain(*Domain) error
}

type HitStore interface {
	SimpleDomainStore
	GetHits(Domain, time.Time, time.Time, map[string]interface{}) ([]*Hit, error)
}
type CounterStore interface {
	SimpleDomainStore
	SaveHit(*Hit) error
}
