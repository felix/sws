package sws

import (
	"testing"
	"time"
)

func TestHitSetSortByDate(t *testing.T) {
	now := time.Now()
	then := now.Add(-10 * time.Hour)
	tests := []struct {
		hits       []*Hit
		begin, end time.Time
		d          time.Duration
	}{
		{
			hits: []*Hit{
				{CreatedAt: now},
				{CreatedAt: then.Add(2 * time.Hour)},
				{CreatedAt: then.Add(time.Hour)},
				{CreatedAt: then.Add(3 * time.Hour)},
			},
			begin: then,
			end:   now,
			d:     time.Hour,
		},
	}

	for i, tt := range tests {
		hs, err := NewHitSet(FromHits(tt.hits))
		if err != nil {
			t.Fatalf("%d => failed %s", i, err)
		}
		hs.duration = tt.d
		hs.SortByDate()

		for j := 0; j < len(hs.buckets)-1; j++ {
			if hs.buckets[j].t.After(hs.buckets[j+1].t) {
				t.Errorf("%d => %d is after %d", i, j, j+1)
			}
		}
	}
}

func TestHitSetFill(t *testing.T) {
	dur := time.Hour
	now := time.Now()
	then := now.Add(-10 * dur)

	tests := []struct {
		hits  []*Hit
		begin *time.Time
		end   *time.Time
	}{
		{
			hits: []*Hit{
				{CreatedAt: now},
				{CreatedAt: then.Add(2 * time.Hour)},
				{CreatedAt: then.Add(time.Hour)},
				{CreatedAt: then.Add(3 * time.Hour)},
			},
			begin: &then,
			end:   &now,
		},
		{
			hits: []*Hit{
				{CreatedAt: then.Add(2 * time.Hour)},
			},
			begin: &then,
			end:   &now,
		},
		{
			hits: []*Hit{
				{CreatedAt: now},
			},
			begin: &then,
			end:   &now,
		},
		{
			hits: []*Hit{
				{CreatedAt: now},
			},
			begin: &then,
		},
		{
			hits: []*Hit{
				{CreatedAt: then},
			},
			end: &now,
		},
		{
			hits: []*Hit{},
			end:  &now,
		},
	}

	for i, tt := range tests {
		hs, err := NewHitSet(FromHits(tt.hits))
		if err != nil {
			t.Fatalf("%d => failed %s", i, err)
		}
		hs.duration = dur
		hs.Fill(tt.begin, tt.end)

		expectedTime := then.Truncate(dur)
		for j, b := range hs.buckets {
			if !b.t.Equal(expectedTime) {
				t.Errorf("%d => expected bucket %d to equal %s, got %s", i, j, expectedTime, b.t)
			}
			expectedTime = expectedTime.Add(dur)
		}
		if len(tt.hits) > 0 {
			total := diffDurations(then, now, dur)
			if len(hs.buckets) != total {
				t.Errorf("%d => expected %d buckets, got %d", i, total, len(hs.buckets))
			}
		}
	}
}
