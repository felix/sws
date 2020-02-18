package sws

import (
	"testing"
	"time"
)

func TestHitsToTimeBuckets(t *testing.T) {
	now := time.Now()

	tests := []struct {
		in       []*Hit
		expected TimeBuckets
	}{
		{
			[]*Hit{
				{ID: ptrInt(1), CreatedAt: ptrTime(now.Add(-5 * time.Second))},
				{ID: ptrInt(2), CreatedAt: ptrTime(now.Add(-4 * time.Second))},
				{ID: ptrInt(3), CreatedAt: ptrTime(now.Add(-4 * time.Second))},
				{ID: ptrInt(4), CreatedAt: ptrTime(now.Add(-4 * time.Second))},
				{ID: ptrInt(5), CreatedAt: ptrTime(now.Add(-3 * time.Second))},
				{ID: ptrInt(6), CreatedAt: ptrTime(now.Add(-2 * time.Second))},
			},
			TimeBuckets{
				Duration: time.Second,
				TimeMin:  now.Add(-5 * time.Second).Round(time.Second),
				TimeMax:  now.Add(-2 * time.Second).Round(time.Second),
				Buckets: []Bucket{
					{Time: now.Add(-5 * time.Second).Round(time.Second), Count: 1},
					{Time: now.Add(-4 * time.Second).Round(time.Second), Count: 3},
					{Time: now.Add(-3 * time.Second).Round(time.Second), Count: 1},
					{Time: now.Add(-2 * time.Second).Round(time.Second), Count: 1},
				},
			},
		},
	}

	for i, tt := range tests {
		actual := HitsToTimeBuckets(tt.in, time.Second)
		if !actual.TimeMin.Equal(tt.expected.TimeMin) {
			t.Errorf("%d => expected %s, got %s", i, tt.expected.TimeMin, actual.TimeMin)
		}
		if !actual.TimeMax.Equal(tt.expected.TimeMax) {
			t.Errorf("%d => expected %s, got %s", i, tt.expected.TimeMax, actual.TimeMax)
		}
	}
}

func TestXYValues(t *testing.T) {
	now := time.Now()

	tests := []struct {
		in   TimeBuckets
		outX []time.Time
		outY []float64
	}{
		{
			TimeBuckets{
				Duration: time.Minute,
				Buckets: []Bucket{
					{Time: now.Add(-3 * time.Minute), Count: 1},
					{Time: now.Add(-2 * time.Minute), Count: 2},
					{Time: now.Add(-1 * time.Minute), Count: 1},
				},
			},
			[]time.Time{
				now.Add(-3 * time.Minute),
				now.Add(-2 * time.Minute),
				now.Add(-1 * time.Minute),
			},
			[]float64{1, 2, 1},
		},
	}

	for i, tt := range tests {
		aX, aY := tt.in.XYValues()
		for j, x := range aX {
			if x != tt.outX[j] {
				t.Errorf("%d => expected [%d] = %s, got %s", i, j, tt.outX[j], x)
			}
		}
		for j, y := range aY {
			if y != tt.outY[j] {
				t.Errorf("%d => expected [%d] = %f, got %f", i, j, tt.outY[j], y)
			}
		}
	}
}

func TestTimeBucketsFill(t *testing.T) {
	now := time.Now().Round(time.Second)

	tests := []struct {
		in         TimeBuckets
		begin, end *time.Time
		expected   []Bucket
	}{
		{
			// End provided
			TimeBuckets{
				Duration: time.Second,
				TimeMin:  now.Add(-5 * time.Second),
				TimeMax:  now.Add(-5 * time.Second),
				Buckets:  []Bucket{{Time: now.Add(-5 * time.Second), Count: 1}},
			},
			nil,
			ptrTime(now),
			[]Bucket{
				{Time: now.Add(-5 * time.Second), Count: 1},
				{Time: now.Add(-4 * time.Second), Count: 0},
				{Time: now.Add(-3 * time.Second), Count: 0},
				{Time: now.Add(-2 * time.Second), Count: 0},
				{Time: now.Add(-1 * time.Second), Count: 0},
				{Time: now, Count: 0},
			},
		},
		{
			// Begin provided
			TimeBuckets{
				Duration: time.Second,
				TimeMin:  now.Add(-5 * time.Second),
				TimeMax:  now.Add(-5 * time.Second),
				Buckets:  []Bucket{{Time: now.Add(-5 * time.Second), Count: 1}},
			},
			ptrTime(now.Add(-6 * time.Second)),
			nil,
			[]Bucket{
				{Time: now.Add(-6 * time.Second), Count: 0},
				{Time: now.Add(-5 * time.Second), Count: 1},
			},
		},
		{
			// No begin or end
			TimeBuckets{
				Duration: time.Second,
				TimeMin:  now.Add(-5 * time.Second),
				TimeMax:  now.Add(-5 * time.Second),
				Buckets:  []Bucket{{Time: now.Add(-5 * time.Second), Count: 1}},
			},
			nil,
			nil,
			[]Bucket{
				{Time: now.Add(-5 * time.Second), Count: 1},
			},
		},
	}

	for i, tt := range tests {
		tt.in.Fill(tt.begin, tt.end)
		for j, b := range tt.in.Buckets {
			if b.Time != tt.expected[j].Time {
				t.Errorf("%d => [%d] expected %s, got %s", i, j, tt.expected[j].Time, b.Time)
			}
			if b.Count != tt.expected[j].Count {
				t.Errorf("%d => [%d] expected %d, got %d", i, j, tt.expected[j].Count, b.Count)
			}
		}
	}
}
