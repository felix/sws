package sws

import (
	"sort"
	"time"
)

type HitSet struct {
	duration   time.Duration
	location   *time.Location
	tMin, tMax *time.Time
	cMin, cMax *int
	hits       []*Hit
	filter     FilterFunc
	buckets    []*bucket
}

type bucket struct {
	t    time.Time
	hits []*Hit
}

type HitSetOption func(*HitSet) error

func TimeZone(s string) HitSetOption {
	return func(hs *HitSet) error {
		var err error
		hs.location, err = time.LoadLocation(s)
		return err
	}
}

func Duration(d time.Duration) HitSetOption {
	return func(hs *HitSet) error {
		hs.duration = d
		return nil
	}
}
func DurationString(s string) HitSetOption {
	return func(hs *HitSet) error {
		d, err := time.ParseDuration(s)
		Duration(d)
		return err
	}
}
func FromHits(hits []*Hit) HitSetOption {
	return func(hs *HitSet) error {
		hs.hits = hits
		return nil
	}
}

func WithFilter(f FilterFunc) HitSetOption {
	return func(hs *HitSet) error {
		hs.filter = f
		return nil
	}
}

// NewHitSet converts a slice of hits to time buckets, group by duration.
func NewHitSet(opts ...HitSetOption) (*HitSet, error) {
	out := &HitSet{
		duration: time.Hour,
		location: time.UTC,
		buckets:  make([]*bucket, 0),
	}
	for _, o := range opts {
		if err := o(out); err != nil {
			return nil, err
		}
	}
	if out.hits != nil {
		for _, h := range out.hits {
			if out.filter == nil || out.filter(h) {
				out.Add(h)
			}
		}
	}
	return out, nil
}

type FilterFunc func(*Hit) bool

func (hs *HitSet) Filter(f FilterFunc) *HitSet {
	out := &HitSet{
		duration: hs.duration,
		location: hs.location,
		tMin:     hs.tMin,
		tMax:     hs.tMax,
		buckets:  make([]*bucket, len(hs.buckets)),
	}
	for i, b := range hs.buckets {
		nb := &bucket{
			t:    b.t,
			hits: make([]*Hit, 0),
		}
		for _, h := range b.hits {
			if f(h) {
				nb.hits = append(nb.hits, h)
			}
		}
		out.buckets[i] = nb
	}
	return out
}

// Hits returns all hits in the set.
func (hs *HitSet) Hits() []*Hit {
	out := make([]*Hit, 0)
	for _, b := range hs.buckets {
		for _, h := range b.hits {
			out = append(out, h)
		}
	}
	return out
}

func (hs *HitSet) Add(h *Hit) {
	hs.updateMinMax()
	k := h.CreatedAt.In(hs.location).Truncate(hs.duration)
	h.CreatedAt = h.CreatedAt.In(hs.location)

	if k.Before(*hs.tMin) {
		hs.tMin = &k
	}
	if k.After(*hs.tMax) {
		hs.tMax = &k
	}

	var bk *bucket
	for _, b := range hs.buckets {
		if b.t.Equal(k) {
			b.hits = append(b.hits, h)
			bk = b
		}
	}
	if bk == nil {
		bk = &bucket{t: k, hits: []*Hit{h}}
		hs.buckets = append(hs.buckets, bk)
	}
	c := len(bk.hits)
	if c < *hs.cMin {
		hs.cMin = &c
	}
	if c > *hs.cMax {
		hs.cMax = &c
	}
}

// Implement Hitter interface.
func (hs *HitSet) Begin() time.Time {
	hs.updateMinMax()
	return *hs.tMin
}

func (hs *HitSet) End() time.Time {
	hs.updateMinMax()
	return *hs.tMax
}

func (hs *HitSet) Duration() time.Duration {
	return hs.duration
}

func (hs *HitSet) Location() *time.Location {
	return hs.location
}

func (hs HitSet) Count() int {
	out := 0
	for _, b := range hs.buckets {
		out += len(b.hits)
	}
	return out
}

func (hs *HitSet) SortByDate() {
	sort.Slice(hs.buckets, func(i, j int) bool {
		return hs.buckets[i].t.Before(hs.buckets[j].t)
	})
}
func (hs *HitSet) SortByHits() {
	sort.Slice(hs.buckets, func(i, j int) bool {
		return len(hs.buckets[i].hits) > len(hs.buckets[j].hits)
	})
}

func (hs *HitSet) Fill(b, e *time.Time) {
	hs.updateMinMax()
	if hs.Count() < 1 {
		return
	}
	begin := *hs.tMin
	if b != nil {
		begin = b.In(hs.location)
	}
	end := *hs.tMax
	if e != nil {
		end = e.In(hs.location)
	}
	begin = begin.Truncate(hs.duration)
	end = end.Truncate(hs.duration)

	total := diffDurations(begin, end, hs.duration)

	newBuckets := make([]*bucket, total)

	hs.SortByDate()

	var existing int
	var idx int
	for n := begin; idx < total && !n.After(end); n = n.Add(hs.duration) {
		switch {
		case existing >= len(hs.buckets):
			newBuckets[idx] = &bucket{t: n, hits: []*Hit{}}

		case n.Before(hs.buckets[existing].t):
			newBuckets[idx] = &bucket{t: n, hits: []*Hit{}}

		default:
			newBuckets[idx] = hs.buckets[existing]
			existing++
		}
		idx++
	}
	hs.tMin = &begin
	hs.tMax = &end
	hs.buckets = newBuckets
}

// XYValues splits the buckets into two data series, one with the times
// and the other with the values.
func (hs HitSet) XYValues() ([]time.Time, []float64) {
	x := make([]time.Time, len(hs.buckets))
	y := make([]float64, len(hs.buckets))
	for i, b := range hs.buckets {
		x[i] = b.t
		y[i] = float64(len(b.hits))
	}
	return x, y
}

func (hs HitSet) YMax() int {
	hs.updateMinMax()
	return *hs.cMax
}
func (hs HitSet) XSeries() []*bucket {
	return hs.buckets
}

func (b bucket) Label() string {
	return b.t.Format("15:04 Jan 2")
}

func (b bucket) Count() int {
	return len(b.hits)
}

func (b bucket) Time() time.Time {
	return b.t
}

func (hs *HitSet) updateMinMax() {
	if hs.tMin != nil && hs.tMax != nil && hs.cMax != nil && hs.cMin != nil {
		return
	}
	if len(hs.buckets) < 1 {
		now := time.Now().Truncate(hs.duration)
		zero := 0
		hs.tMin, hs.tMax = &now, &now
		hs.cMin, hs.cMax = &zero, &zero
		return
	}
	lHits := len(hs.buckets[0].hits)
	hs.cMin = &lHits
	hs.cMax = &lHits
	hs.tMin = &hs.buckets[0].t
	hs.tMax = &hs.buckets[0].t

	for _, b := range hs.buckets {
		c := len(b.hits)
		if c < *hs.cMin {
			hs.cMin = &c
		}
		if c > *hs.cMax {
			hs.cMax = &c
		}
		if b.t.Before(*hs.tMin) {
			hs.tMin = &b.t
		}
		if b.t.After(*hs.tMax) {
			hs.tMax = &b.t
		}
	}
}
