package sws

type PageSet map[string]*Page

func NewPageSet(hitter Hitter) PageSet {
	out := make(map[string]*Page)
	for _, h := range hitter.Hits() {
		p, ok := out[h.Path]
		if !ok {
			p = &Page{
				Path:          h.Path,
				SiteID:        *h.SiteID,
				Title:         h.Title,
				LastVisitedAt: h.CreatedAt,
				hitSet: &HitSet{
					duration: hitter.Duration(),
				},
			}
		}
		if p.LastVisitedAt.Before(h.CreatedAt) {
			p.LastVisitedAt = h.CreatedAt
		}
		p.hitSet.Add(h)
		out[h.Path] = p
	}
	b := hitter.Begin()
	e := hitter.End()
	for _, p := range out {
		p.hitSet.Fill(&b, &e)
	}
	return PageSet(out)
}

func (ps PageSet) Page(p string) *Page {
	pg, _ := ps[p]
	return pg
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
	out := make([]*Page, len(ps))
	i := 0
	for _, v := range ps {
		out[i] = v
		i++
	}
	return out
}
