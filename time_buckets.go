package sws

import (
	"fmt"
	"sort"
	"time"
)

type TimeBuckets struct {
	Duration           time.Duration
	TimeMin, TimeMax   time.Time
	CountMin, CountMax int
	Buckets            []Bucket
}

type Bucket struct {
	Time  time.Time
	Count int
}

// XYValues splits the buckets into two data series, one with the times
// and the other with the values.
func (tb TimeBuckets) XYValues() ([]time.Time, []float64) {
	x := make([]time.Time, len(tb.Buckets))
	y := make([]float64, len(tb.Buckets))
	for i, b := range tb.Buckets {
		x[i] = b.Time
		y[i] = float64(b.Count)
	}
	return x, y
}

func (b Bucket) String() string {
	return fmt.Sprintf("%s => %d", b.Time, b.Count)
}

// HitsToTimeBuckets converts a slice of hits to time buckets, group by durtation.
func HitsToTimeBuckets(hits []*Hit, d time.Duration) TimeBuckets {
	out := TimeBuckets{
		Duration: d,
		Buckets:  make([]Bucket, 0),
	}
	for j, h := range hits {
		k := h.CreatedAt.Round(d)
		if j == 0 || k.Before(out.TimeMin) {
			out.TimeMin = k
		}
		if j == 0 || k.After(out.TimeMax) {
			out.TimeMax = k
		}
		var found bool
		for i, tb := range out.Buckets {
			if tb.Time.Equal(k) {
				out.Buckets[i].Count++
				found = true
			}
		}
		if !found {
			out.Buckets = append(out.Buckets, Bucket{Time: k, Count: 1})
		}
	}
	out.Sort()
	return out
}

// Sort order the buckets in ascending order by time.
func (tb *TimeBuckets) Sort() {
	sort.Slice(tb.Buckets, func(i, j int) bool {
		return tb.Buckets[i].Time.Before(tb.Buckets[j].Time)
	})
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

	total := diffDurations(begin, end, tb.Duration)
	tb.Sort()

	newBuckets := make([]Bucket, total)

	var existing int
	var idx int
	for n := begin; idx < total && !n.After(end); n = n.Add(tb.Duration) {
		switch {
		case existing >= len(tb.Buckets):
			newBuckets[idx] = Bucket{Time: n, Count: 0}

		case n.Before(tb.Buckets[existing].Time):
			newBuckets[idx] = Bucket{Time: n, Count: 0}

		default:
			newBuckets[idx] = tb.Buckets[existing]
			existing++
		}
		idx++
	}
	tb.Buckets = newBuckets
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
