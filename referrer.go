package sws

import (
	"net/url"
	"sort"
	"time"
)

type Referrer struct {
	Name       string    `json:"name"`
	URL        string    `json:"url"`
	LastSeenAt time.Time `json:"last_seen_at" db:"last_seen_at"`
	hitSet     *HitSet
}

type ReferrerSet []*Referrer

func NewReferrerSet(hs *HitSet, site Site) *ReferrerSet {
	tmp := make(map[string]*Referrer)
	for _, h := range hs.Hits() {
		host := "direct"
		u := ""

		if h.Referrer != nil {
			if r, err := url.Parse(*h.Referrer); err == nil {
				host = r.Host
			}
			u = *h.Referrer
		}
		// Check for internal referrer
		if site.IncludesDomain(host) {
			//host = "internal"
			continue
		}
		tmp[host] = &Referrer{
			Name: host,
			URL:  u,
			hitSet: hs.Filter(func(t *Hit) bool {
				if h.Referrer == nil && t.Referrer == nil {
					return true
				}
				if h.Referrer == nil && t.Referrer != nil {
					return false
				}
				if h.Referrer != nil && t.Referrer == nil {
					return false
				}
				return *t.Referrer == *t.Referrer
			}),
		}
	}
	if len(tmp) < 1 {
		return nil
	}
	out := make([]*Referrer, len(tmp))
	i := 0
	for _, b := range tmp {
		out[i] = b
		i++
	}
	rs := ReferrerSet(out)
	return &rs
}

func (rs *ReferrerSet) SortByName() {
	sort.Slice(*rs, func(i, j int) bool {
		return (*rs)[i].Name < (*rs)[j].Name
	})
}

func (rs *ReferrerSet) SortByHits() {
	sort.Slice(*rs, func(i, j int) bool {
		return (*rs)[i].hitSet.Count() > (*rs)[j].hitSet.Count()
	})
}

func (rs ReferrerSet) Count() int {
	return len(rs)
}

func (rs ReferrerSet) GetReferrer(s string) *Referrer {
	for _, r := range rs {
		if r.Name == s {
			return r
		}
	}
	return nil
}

func (r Referrer) Label() string {
	return r.Name
}

func (r Referrer) Count() int {
	return r.hitSet.Count()
}

func (r Referrer) YValue() int {
	return r.hitSet.Count()
}

func (rs ReferrerSet) YMax() int {
	max := 0
	for _, r := range rs {
		if r.hitSet.Count() > max {
			max = r.hitSet.Count()
		}
	}
	return max
}
func (rs ReferrerSet) YSum() int {
	sum := 0
	for _, r := range rs {
		sum += r.hitSet.Count()
	}
	return sum
}
func (rs ReferrerSet) XSeries() []*Referrer {
	return rs
}
