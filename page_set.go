package sws

import (
	"sort"
)

type PageSet []*Page

func NewPageSet(hs *HitSet) (PageSet, error) {
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
	return PageSet(out), nil
}

func (ps PageSet) Hits() []*Hit {
	out := make([]*Hit, 0)
	for _, p := range ps {
		out = append(out, p.hitSet.Hits()...)
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

func (ps PageSet) Page(s string) *Page {
	for _, p := range ps {
		if p.Path == s {
			return p
		}
	}
	return nil
}

func (ps PageSet) YMax() int {
	max := 0
	for _, p := range ps {
		if p.Count() > max {
			max = p.Count()
		}
	}
	return max
}
func (ps PageSet) XSeries() []*Page {
	return ps
}
