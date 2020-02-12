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
