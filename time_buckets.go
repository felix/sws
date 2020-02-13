package sws

import (
	"fmt"
	"sort"
	"time"
)

type TimeBuckets struct {
	Duration         time.Duration
	TimeMin, TimeMax time.Time
	//CountMin, CountMax int
	Buckets []Bucket
}

type Bucket struct {
	Time  time.Time
	Count int
}

func (tb TimeBuckets) FilledXYValues() ([]time.Time, []float64) {
	x := make([]time.Time, len(tb.Buckets))
	y := make([]float64, len(tb.Buckets))
	for i, b := range tb.Buckets {
		x[i] = b.Time
		y[i] = float64(b.Count)
	}
	return durationsFilled(x, y, tb.Duration)
}

func (b Bucket) String() string {
	return fmt.Sprintf("%s => %d", b.Time, b.Count)
}

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
	sort.Slice(out.Buckets, func(i, j int) bool {
		return out.Buckets[i].Time.Before(out.Buckets[j].Time)
	})
	return out
}

// TODO pull from chart upstream
func durations(start time.Time, total int, d time.Duration) []time.Time {
	times := make([]time.Time, total)

	last := start
	for i := 0; i < total; i++ {
		times[i] = last
		last = last.Add(d)
	}
	return times
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

func timeMinMax(times ...time.Time) (min, max time.Time) {
	if len(times) == 0 {
		return
	}
	min = times[0]
	max = times[0]

	for index := 1; index < len(times); index++ {
		if times[index].Before(min) {
			min = times[index]
		}
		if times[index].After(max) {
			max = times[index]
		}
	}
	return
}

func durationsFilled(xdata []time.Time, ydata []float64, d time.Duration) ([]time.Time, []float64) {
	start, end := timeMinMax(xdata...)
	totalHours := diffDurations(start, end, d)

	finalTimes := durations(start, totalHours+1, d)
	finalValues := make([]float64, totalHours+1)

	var hoursFromStart int
	for i, xd := range xdata {
		hoursFromStart = diffDurations(start, xd, d)
		finalValues[hoursFromStart] = ydata[i]
	}

	return finalTimes, finalValues
}
