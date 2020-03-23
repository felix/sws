package sws

import (
	"sort"
	"time"

	detector "github.com/mssola/user_agent"
)

type Browser struct {
	Name       string    `json:"name"`
	LastSeenAt time.Time `json:"last_seen_at" db:"last_seen_at"`
	hitSet     *HitSet
}

type BrowserSet []*Browser

func NewBrowserSet(hs *HitSet) *BrowserSet {
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
	if len(tmp) < 1 {
		return nil
	}
	out := make([]*Browser, len(tmp))
	i := 0
	for _, b := range tmp {
		out[i] = b
		i++
	}
	bs := BrowserSet(out)
	return &bs
}
func (bs *BrowserSet) SortByName() {
	sort.Slice(*bs, func(i, j int) bool {
		return (*bs)[i].Label() < (*bs)[j].Label()
	})
}

func (bs *BrowserSet) SortByHits() {
	sort.Slice(*bs, func(i, j int) bool {
		return (*bs)[i].hitSet.Count() > (*bs)[j].hitSet.Count()
	})
}

func (b Browser) Label() string {
	return b.Name
}

func (b Browser) Count() int {
	return b.hitSet.Count()
}

func (bs BrowserSet) Count() int {
	return len(bs)
}
func (bs BrowserSet) Labels() []string {
	out := make([]string, len(bs))
	for i := 0; i < len(bs); i++ {
		out[i] = bs[i].Label()
	}
	return out
}
func (bs BrowserSet) Counts() []int {
	out := make([]int, len(bs))
	for i := 0; i < len(bs); i++ {
		out[i] = bs[i].Count()
	}
	return out
}

/*
func (bs BrowserSet) Ratios() []float64 {
	out := make([]float64, len(bs))
	max := 0.0
	for i := 0; i < len(bs); i++ {
		out[i] = float64(bs[i].Count())
		if out[i] > max {
			max = out[i]
		}
	}
	for i := 0; i < len(out); i++ {
		out[i] = out[i] / max
	}
	return out
}
*/
func (bs BrowserSet) YMax() int {
	max := 0
	for _, b := range bs {
		if b.hitSet.Count() > max {
			max = b.hitSet.Count()
		}
	}
	return max
}
func (bs BrowserSet) YSum() int {
	sum := 0
	for _, b := range bs {
		sum += b.hitSet.Count()
	}
	return sum
}
func (bs BrowserSet) XSeries() []*Browser {
	return bs
}
