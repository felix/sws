package sws

import (
	"sort"
)

type PageSet []*Page

func NewPageSet(hs *HitSet) (*PageSet, error) {
	tmp := make(map[string]*Page)
	for _, h := range hs.Hits() {
		if _, ok := tmp[h.Path]; ok {
			// Already captured this path
			continue
		}
		p := &Page{
			Path:          h.Path,
			SiteID:        *h.SiteID,
			Title:         h.Title,
			LastVisitedAt: h.CreatedAt,
			hitSet: hs.Filter(func(t *Hit) bool {
				return t.Path == h.Path
			}),
		}
		// if p.LastVisitedAt.Before(h.CreatedAt) {
		// 	p.LastVisitedAt = h.CreatedAt
		// }
		//p.hitSet.Add(h)
		tmp[h.Path] = p
	}
	if len(tmp) < 1 {
		return nil, nil
	}
	out := make([]*Page, len(tmp))
	i := 0
	for _, p := range tmp {
		out[i] = p
		i++
	}
	ps := PageSet(out)
	return &ps, nil
}

func (ps *PageSet) Count() int {
	return len(*ps)
}

func (ps PageSet) Hits() []*Hit {
	out := make([]*Hit, 0)
	for i := 0; i < len(ps); i++ {
		out = append(out, ps[i].hitSet.Hits()...)
	}
	return out
}

func (ps *PageSet) SortByPath() {
	sort.Slice(*ps, func(i, j int) bool {
		return (*ps)[i].Path < (*ps)[j].Path
	})
}

func (ps *PageSet) SortByHits() {
	sort.Slice(*ps, func(i, j int) bool {
		return (*ps)[i].hitSet.Count() > (*ps)[j].hitSet.Count()
	})
}

func (ps PageSet) GetPage(s string) *Page {
	for i := 0; i < len(ps); i++ {
		if ps[i].Path == s {
			return ps[i]
		}
	}
	return nil
}

func (ps PageSet) YMax() int {
	max := 0
	for i := 0; i < len(ps); i++ {
		if ps[i].Count() > max {
			max = ps[i].Count()
		}
	}
	return max
}
func (ps PageSet) YSum() int {
	sum := 0
	for _, p := range ps {
		sum += p.hitSet.Count()
	}
	return sum
}
func (ps PageSet) XSeries() []*Page {
	max := 10
	if len(ps) < 10 {
		max = len(ps)
	}
	return ps[0:max]
}
