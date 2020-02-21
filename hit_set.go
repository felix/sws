package sws

import (
	"sort"
	"time"
)

type HitSet struct {
	duration   time.Duration
	tMin, tMax time.Time
	cMin, cMax int
	hits       []*Hit
	buckets    []*bucket
}

type bucket struct {
	t    time.Time
	hits []*Hit
}

// NewHitSet converts a slice of hits to time buckets, group by duration.
func NewHitSet(hits []*Hit, b, e time.Time, d time.Duration) *HitSet {
	out := &HitSet{
		duration: d,
		buckets:  make([]*bucket, 0),
	}
	for _, h := range hits {
		out.Add(h)
	}
	//out.updateMinMax()
	sort.Sort(out)
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
	SortHits(out)
	return out
}

func (hs *HitSet) Add(h *Hit) {
	k := h.CreatedAt.Truncate(hs.duration)
	if k.Before(hs.tMin) {
		hs.tMin = k
	}
	if k.After(hs.tMax) {
		hs.tMax = k
	}

	var bk *bucket
	for _, b := range hs.buckets {
		if b.t.Equal(k) {
			bk = b
		}
	}
	if bk == nil {
		// Create new bucket
		bk = &bucket{t: k, hits: []*Hit{h}}
		hs.buckets = append(hs.buckets, bk)
	}
	c := len(bk.hits)
	if c < hs.cMin {
		hs.cMin = c
	}
	if c > hs.cMax {
		hs.cMax = c
	}
}

// Implement Hitter interface.
func (hs *HitSet) Begin() time.Time {
	return hs.tMin
}

func (hs *HitSet) End() time.Time {
	return hs.tMax
}

func (hs *HitSet) Duration() time.Duration {
	return hs.duration
}

func (hs HitSet) Count() int {
	out := 0
	for _, b := range hs.buckets {
		out += len(b.hits)
	}
	return out
}

// Implement sort.Interface
func (hs HitSet) Len() int           { return len(hs.buckets) }
func (hs HitSet) Less(i, j int) bool { return hs.buckets[i].t.Before(hs.buckets[i].t) }
func (hs HitSet) Swap(i, j int)      { hs.buckets[i], hs.buckets[j] = hs.buckets[j], hs.buckets[i] }

func (hs *HitSet) Fill(b, e *time.Time) {
	begin := hs.tMin
	if b != nil {
		begin = *b
	}
	end := hs.tMax
	if e != nil {
		end = *e
	}

	total := diffDurations(begin, end, hs.duration)

	newBuckets := make([]*bucket, total)

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
	hs.tMin = begin
	hs.tMax = end
	//hs.updateMinMax()
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
	return hs.cMax
}
func (hs HitSet) XSeries() []*bucket {
	return hs.buckets
}

func (b bucket) Label() string {
	return b.t.Format("15:04 Jan 2")
}

func (b bucket) YValue() int {
	return len(b.hits)
}

func (b bucket) Time() time.Time {
	return b.t
}

/*
func (hs *HitSet) updateMinMax() {
	if len(hs.buckets) < 1 {
		return
	}
	hs.cMin = len(hs.buckets[0].hits)
	hs.cMax = len(hs.buckets[0].hits)
	hs.tMin = hs.buckets[0].t
	hs.tMax = hs.buckets[0].t
	for _, b := range hs.buckets {
		c := len(b.hits)
		if c < hs.cMin {
			hs.cMin = c
		}
		if c > hs.cMax {
			hs.cMax = c
		}
		if b.t.Before(hs.tMin) {
			hs.tMin = b.t
		}
		if b.t.After(hs.tMax) {
			hs.tMax = b.t
		}
	}
}
*/
