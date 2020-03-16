package sws

import (
	"net/url"
	"sort"
	"strings"
	"time"
)

type Referrer struct {
	Name       string    `json:"name"`
	LastSeenAt time.Time `json:"last_seen_at" db:"last_seen_at"`
	hitSet     *HitSet
}

type ReferrerSet []*Referrer

func NewReferrerSet(hs *HitSet) ReferrerSet {
	tmp := make(map[string]*Referrer)
	for _, h := range hs.Hits() {
		if h.Referrer == nil {
			continue
		}

		u, err := url.Parse(*h.Referrer)
		if err != nil || h.Host == u.Host {
			continue
		}
		host := u.Host
		if u.Host == "" {
			host = "direct"
		}
		r := &Referrer{
			Name:       host,
			LastSeenAt: h.CreatedAt,
			hitSet: hs.Filter(func(t *Hit) bool {
				if t.Referrer == nil {
					return false
				}
				return strings.Contains(*t.Referrer, u.Host)
			}),
		}
		// if b.LastSeenAt.Before(h.CreatedAt) {
		// 	b.LastSeenAt = h.CreatedAt
		// }
		//b.hitSet.Add(h)
		tmp[u.Host] = r
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
	return ReferrerSet(out)
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
