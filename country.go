package sws

import (
	"net"
	"sort"

	maxminddb "github.com/oschwald/maxminddb-golang"
)

type Country struct {
	Name   string `json:"name"`
	hitSet *HitSet
}

type CountrySet []*Country

func NewCountrySet(hs *HitSet) *CountrySet {
	tmp := make(map[string]*Country)
	for _, h := range hs.Hits() {
		if h.CountryCode == nil {
			continue
		}
		cc := *h.CountryCode

		if _, ok := tmp[cc]; ok {
			// Already captured this country
			continue
		}
		b := &Country{
			Name: cc,
			//LastSeenAt: h.CreatedAt,
			hitSet: hs.Filter(func(t *Hit) bool {
				return t.CountryCode != nil && *t.CountryCode == cc
			}),
		}
		tmp[cc] = b
	}
	if len(tmp) < 1 {
		return nil
	}
	out := make([]*Country, len(tmp))
	i := 0
	for _, b := range tmp {
		out[i] = b
		i++
	}
	cs := CountrySet(out)
	return &cs
}

func (cs *CountrySet) SortByHits() {
	sort.Slice(*cs, func(i, j int) bool {
		return (*cs)[i].hitSet.Count() > (*cs)[j].hitSet.Count()
	})
}

func (c Country) Label() string {
	return c.Name
}

func (c Country) Count() int {
	return c.hitSet.Count()
}

func (cs CountrySet) Count() int {
	return len(cs)
}

func (cs CountrySet) YMax() int {
	max := 0
	for _, c := range cs {
		if c.hitSet.Count() > max {
			max = c.hitSet.Count()
		}
	}
	return max
}
func (cs CountrySet) YSum() int {
	sum := 0
	for _, c := range cs {
		sum += c.hitSet.Count()
	}
	return sum
}
func (cs CountrySet) XSeries() []*Country {
	return cs
}

func FetchCountryCode(path, host string) (*string, error) {
	db, err := maxminddb.Open(path)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	ip := net.ParseIP(host)
	var r struct {
		Country struct {
			ISOCode string `maxminddb:"iso_code"`
		} `maxminddb:"country"`
	}
	if err := db.Lookup(ip, &r); err != nil {
		return nil, err
	}
	return &r.Country.ISOCode, nil
}
