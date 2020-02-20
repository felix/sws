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
	Data               []Bucket
}

type Bucket struct {
	Time  time.Time
	Count int
}

/*
func (tb TimeBuckets) YMax() int {
	return tb.CountMax
}

func (tb TimeBuckets) Next() TimeCountable {
	return tb.Data
}

func (b Bucket) XValue() time.Time {
	return b.Time
}

func (b Bucket) YValue() int {
	return b.Count
}
*/

// XYValues splits the buckets into two data series, one with the times
// and the other with the values.
func (tb TimeBuckets) XYValues() ([]time.Time, []float64) {
	x := make([]time.Time, len(tb.Data))
	y := make([]float64, len(tb.Data))
	for i, b := range tb.Data {
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
		Data:     make([]Bucket, 0),
	}
	for j, h := range hits {
		k := h.CreatedAt.Truncate(d)
		if j == 0 || k.Before(out.TimeMin) {
			out.TimeMin = k
		}
		if j == 0 || k.After(out.TimeMax) {
			out.TimeMax = k
		}
		var found bool
		for i, tb := range out.Data {
			if tb.Time.Equal(k) {
				out.Data[i].Count++
				found = true
			}
		}
		if !found {
			out.Data = append(out.Data, Bucket{Time: k, Count: 1})
		}
	}
	out.Sort()
	out.updateMinMax()
	return out
}

// Sort order the buckets in ascending order by time.
func (tb *TimeBuckets) Sort() {
	sort.Slice(tb.Data, func(i, j int) bool {
		return tb.Data[i].Time.Before(tb.Data[j].Time)
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
		case existing >= len(tb.Data):
			newBuckets[idx] = Bucket{Time: n, Count: 0}

		case n.Before(tb.Data[existing].Time):
			newBuckets[idx] = Bucket{Time: n, Count: 0}

		default:
			newBuckets[idx] = tb.Data[existing]
			existing++
		}
		idx++
	}
	tb.updateMinMax()
	tb.Data = newBuckets
}

func (tb *TimeBuckets) updateMinMax() {
	if len(tb.Data) < 1 {
		return
	}
	minC := tb.Data[0].Count
	maxC := tb.Data[0].Count
	minT := tb.Data[0].Time
	maxT := tb.Data[0].Time
	for _, b := range tb.Data {
		if b.Count < minC {
			minC = b.Count
		}
		if b.Count > maxC {
			maxC = b.Count
		}
		if b.Time.Before(minT) {
			minT = b.Time
		}
		if b.Time.After(maxT) {
			maxT = b.Time
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
