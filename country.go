package sws

type Country struct {
	Name   string `json:"name"`
	hitSet *HitSet
}

type CountrySet []*Country

func NewCountrySet(hs *HitSet) CountrySet {
	tmp := make(map[string]*Country)
	for _, h := range hs.Hits() {
		if h.CountryCode == nil {
			continue
		}
		cc := *h.CountryCode

		if _, ok := tmp[cc]; ok {
			// Already captured this country
			continue
		}
		b := &Country{
			Name: cc,
			//LastSeenAt: h.CreatedAt,
			hitSet: hs.Filter(func(t *Hit) bool {
				return t.CountryCode != nil && *t.CountryCode == cc
			}),
		}
		tmp[cc] = b
	}
	if len(tmp) < 1 {
		return nil
	}
	out := make([]*Country, len(tmp))
	i := 0
	for _, b := range tmp {
		out[i] = b
		i++
	}
	return CountrySet(out)
}

func (c Country) Label() string {
	return c.Name
}

func (c Country) Count() int {
	return c.hitSet.Count()
}

func (cs CountrySet) Count() int {
	return len(cs)
}

func (cs CountrySet) YMax() int {
	max := 0
	for _, c := range cs {
		if c.hitSet.Count() > max {
			max = c.hitSet.Count()
		}
	}
	return max
}
func (cs CountrySet) XSeries() []*Country {
	return cs
}
