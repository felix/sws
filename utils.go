package sws

import (
	"github.com/speps/go-hashids"
	"time"
)

func ptrString(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func ptrInt(i int) *int {
	return &i
}

func ptrTime(t time.Time) *time.Time {
	return &t
}

func hashID(salt string, id int) string {
	hd := hashids.NewData()
	hd.Salt = salt
	h, _ := hashids.NewWithData(hd)
	out, _ := h.Encode([]int{id})
	return out
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
