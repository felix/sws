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
		b, ok := out[*h.UserAgentHash]
		if !ok {
			b = &UserAgent{
				Name:       h.UserAgent.Name,
				LastSeenAt: h.CreatedAt,
				hitSet:     &HitSet{},
				ua:         detector.New(h.UserAgent.Name),
			}
		}
		if b.LastSeenAt.Before(h.CreatedAt) {
			b.LastSeenAt = h.CreatedAt
		}
		b.hitSet.Add(h)
		out[*h.UserAgentHash] = b
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
func (s UserAgentSet) XSeries() []*UserAgent {
	out := make([]*UserAgent, len(s))
	i := 0
	for _, v := range s {
		out[i] = v
		i++
	}
	return out
}
