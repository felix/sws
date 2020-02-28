package sws

import (
	"time"

	detector "github.com/mssola/user_agent"
)

type Browser struct {
	Name       string    `json:"name"`
	LastSeenAt time.Time `json:"last_seen_at" db:"last_seen_at"`
	hitSet     *HitSet
}

type BrowserSet []*Browser

func NewBrowserSet(hs *HitSet) BrowserSet {
	tmp := make(map[string]*Browser)
	for _, h := range hs.Hits() {
		browser := ""
		if h.UserAgentHash != nil {
			d := detector.New(h.UserAgent.Name)
			browser, _ = d.Browser()
		}
		if _, ok := tmp[browser]; ok {
			// Already captured this UA
			continue
		}
		b := &Browser{
			Name:       browser,
			LastSeenAt: h.CreatedAt,
			hitSet: hs.Filter(func(t *Hit) bool {
				if t.UserAgentHash == nil {
					return browser == ""
				}
				test, _ := detector.New(t.UserAgent.Name).Browser()
				return browser == test
			}),
		}
		// if b.LastSeenAt.Before(h.CreatedAt) {
		// 	b.LastSeenAt = h.CreatedAt
		// }
		//b.hitSet.Add(h)
		tmp[browser] = b
	}
	out := make([]*Browser, len(tmp))
	i := 0
	for _, b := range tmp {
		out[i] = b
		i++
	}
	return BrowserSet(out)
}

func (b Browser) Label() string {
	return b.Name
}

func (b Browser) Count() int {
	return b.hitSet.Count()
}

func (b Browser) YValue() int {
	return b.hitSet.Count()
}

func (bs BrowserSet) YMax() int {
	max := 0
	for _, b := range bs {
		if b.hitSet.Count() > max {
			max = b.hitSet.Count()
		}
	}
	return max
}
func (bs BrowserSet) XSeries() []*Browser {
	return bs
}
