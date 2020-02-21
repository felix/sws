package sws

import (
	detector "github.com/mssola/user_agent"
)

type UserAgentSet map[string]*UserAgent

// NewUserAgentSet collects the browsers from provided hits.
func NewUserAgentSet(hitter Hitter) UserAgentSet {
	out := make(map[string]*UserAgent)
	for _, h := range hitter.Hits() {
		if h.UserAgentHash == nil {
			// TODO
			continue
		}
		d := detector.New(h.UserAgent.Name)
		browser, _ := d.Browser()
		b, ok := out[browser]
		if !ok {
			b = &UserAgent{
				Name:       h.UserAgent.Name,
				LastSeenAt: h.CreatedAt,
				hitSet:     &HitSet{},
				ua:         d,
			}
		}
		if b.LastSeenAt.Before(h.CreatedAt) {
			b.LastSeenAt = h.CreatedAt
		}
		b.hitSet.Add(h)
		out[browser] = b
	}
	return UserAgentSet(out)
}

func (uas UserAgentSet) YMax() int {
	max := 0
	for _, ua := range uas {
		if ua.Count() > max {
			max = ua.Count()
		}
	}
	return max
}
func (uas UserAgentSet) XSeries() []*UserAgent {
	out := make([]*UserAgent, len(uas))
	i := 0
	for _, v := range uas {
		out[i] = v
		i++
	}
	return out
}
