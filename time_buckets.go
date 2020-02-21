// +build ignore

package sws

import (
	"sort"
	"time"
)

type TimeBuckets struct {
	duration           time.Duration
	TimeMin, TimeMax   time.Time
	CountMin, CountMax int
	buckets            []Bucket
}

type Bucket struct {
	t    time.Time
	hits []*Hit
}

func (tb TimeBuckets) Hits() []*Hit {
	out := make([]*Hit, 0)
	for _, b := range tb.buckets {
		for _, h := range b.hits {
			out = append(out, h)
		}
	}
	SortHits(out)
	return out
}

// Implement sort.Interface
func (tb TimeBuckets) Len() int           { return len(tb.buckets) }
func (tb TimeBuckets) Less(i, j int) bool { return tb.buckets[i].t.Before(tb.buckets[i].t) }
func (tb TimeBuckets) Swap(i, j int)      { tb.buckets[i], tb.buckets[j] = tb.buckets[j], tb.buckets[i] }

// XYValues splits the buckets into two data series, one with the times
// and the other with the values.
func (tb TimeBuckets) XYValues() ([]time.Time, []float64) {
	x := make([]time.Time, len(tb.buckets))
	y := make([]float64, len(tb.buckets))
	for i, b := range tb.buckets {
		x[i] = b.t
		y[i] = float64(len(b.hits))
	}
	return x, y
}

// HitsToTimeBuckets converts a slice of hits to time buckets, group by durtation.
func HitsToTimeBuckets(hits []*Hit, d time.Duration) TimeBuckets {
	out := TimeBuckets{
		duration: d,
		buckets:  make([]Bucket, 0),
	}
	SortHits(hits)
	for j, h := range hits {
		k := h.CreatedAt.Truncate(d)
		if j == 0 || k.Before(out.TimeMin) {
			out.TimeMin = k
		}
		if j == 0 || k.After(out.TimeMax) {
			out.TimeMax = k
		}
		var found bool
		for i, tb := range out.buckets {
			if tb.t.Equal(k) {
				out.buckets[i].hits = append(out.buckets[i].hits, h)
				found = true
			}
		}
		if !found {
			out.buckets = append(out.buckets, Bucket{t: k, hits: []*Hit{h}})
		}
	}
	out.updateMinMax()
	sort.Sort(out)
	return out
}

// Fill adds extra buckets so each duration segment has a bucket.
// If no begin or end times are provided it uses the existing min and max times.
func (tb *TimeBuckets) Fill(b, e *time.Time) {
	begin := tb.TimeMin
	if b != nil {
		begin = *b
	}
	end := tb.TimeMax
	if e != nil {
		end = *e
	}

	total := diffDurations(begin, end, tb.duration)

	newBuckets := make([]Bucket, total)

	var existing int
	var idx int
	for n := begin; idx < total && !n.After(end); n = n.Add(tb.duration) {
		switch {
		case existing >= len(tb.buckets):
			newBuckets[idx] = Bucket{t: n, hits: []*Hit{}}

		case n.Before(tb.buckets[existing].t):
			newBuckets[idx] = Bucket{t: n, hits: []*Hit{}}

		default:
			newBuckets[idx] = tb.buckets[existing]
			existing++
		}
		idx++
	}
	tb.updateMinMax()
	tb.buckets = newBuckets
}

func (tb TimeBuckets) YMax() int {
	return tb.CountMax
}
func (tb TimeBuckets) XSeries() []Bucket {
	return tb.buckets
}

func (b Bucket) Label() string {
	return b.t.Format("15:04 Jan 2")
}

func (b Bucket) Time() string {
	return b.t.Format("15:04 Jan 2")
}

func (b Bucket) YValue() int {
	return len(b.hits)
}

func (tb *TimeBuckets) updateMinMax() {
	if len(tb.buckets) < 1 {
		return
	}
	minC := len(tb.buckets[0].hits)
	maxC := len(tb.buckets[0].hits)
	minT := tb.buckets[0].t
	maxT := tb.buckets[0].t
	for _, b := range tb.buckets {
		c := len(b.hits)
		if c < minC {
			minC = c
		}
		if c > maxC {
			maxC = c
		}
		if b.t.Before(minT) {
			minT = b.t
		}
		if b.t.After(maxT) {
			maxT = b.t
		}
	}
	tb.TimeMin = minT
	tb.TimeMax = maxT
	tb.CountMin = minC
	tb.CountMax = maxC
}

func diffDurations(t1, t2 time.Time, d time.Duration) int {
	t1n := t1.Unix()
	t2n := t2.Unix()
	var diff int64
	if t1n > t2n {
		diff = t1n - t2n
	} else {
		diff = t2n - t1n
	}
	return int(float64(diff) / d.Seconds())
}
