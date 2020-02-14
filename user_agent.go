package sws

import (
	"crypto/sha1"
	"fmt"
	"net/http"
	"time"
)

type UserAgent struct {
	Hash       string    `json:"hash"`
	Name       string    `json:"name"`
	LastSeenAt time.Time `json:"last_seen_at"`
}

func UserAgentHash(s string) string {
	return fmt.Sprintf("%x", sha1.Sum([]byte(s)))
}

func UserAgentFromRequest(r *http.Request) (*UserAgent, error) {
	q := r.URL.Query()
	agent := q.Get("u")
	if agent == "" {
		return nil, nil
	}
	ua := r.UserAgent()

	return &UserAgent{
		Name:       ua,
		LastSeenAt: time.Now(),
		Hash:       UserAgentHash(ua),
	}, nil
}
