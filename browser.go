package sws

import (
	"time"
)

type Browser struct {
	Name       string
	UserAgent  string
	LastSeenAt time.Time
	Engine     string
	Count      int
}

func BrowsersFromHits(hits []*Hit) map[string]*Browser {
	out := make(map[string]*Browser)
	for _, h := range hits {
		if h.UserAgentHash != nil {
			b, ok := out[*h.UserAgentHash]
			if !ok {
				b = &Browser{
					// TODO name
					UserAgent:  h.UserAgent.Name,
					LastSeenAt: h.CreatedAt,
				}
			}
			if b.LastSeenAt.Before(h.CreatedAt) {
				b.LastSeenAt = h.CreatedAt
			}
			b.Count++
			out[*h.UserAgentHash] = b
		}
	}
	return out
}
