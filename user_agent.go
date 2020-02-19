package sws

import (
	"crypto/sha1"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	detector "github.com/mssola/user_agent"
)

// UserAgent of a hit.
type UserAgent struct {
	Hash       string    `json:"hash"`
	Name       string    `json:"name"`
	LastSeenAt time.Time `json:"last_seen_at" db:"last_seen_at"`
	Count      int

	ua *detector.UserAgent
}

var (
	reBotWord, reBotSite *regexp.Regexp
)

type browserMatcher func(string) (string, bool)

func init() {
	reBotWord = regexp.MustCompile("(?i)(bot|crawler|sp(i|y)der|search|worm|fetch|nutch)")
	reBotSite = regexp.MustCompile("http[s]?://.+\\.\\w+")
}

// UserAgentHash is the UA key.
func UserAgentHash(s string) string {
	return fmt.Sprintf("%x", sha1.Sum([]byte(s)))
}

// UserAgentFromRequest extracts a UA from a request.
func UserAgentFromRequest(r *http.Request) (*UserAgent, error) {
	q := r.URL.Query()
	ua := q.Get("u")
	if ua == "" {
		ua = r.UserAgent()
	}

	return &UserAgent{
		Name:       ua,
		LastSeenAt: time.Now(),
		Hash:       UserAgentHash(ua),
		ua:         detector.New(ua),
	}, nil
}

// UserAgentsFromHits collects the browsers from provided hits.
func UserAgentsFromHits(hits []*Hit) map[string]*UserAgent {
	out := make(map[string]*UserAgent)
	for _, h := range hits {
		if h.UserAgentHash != nil {
			b, ok := out[*h.UserAgentHash]
			if !ok {
				b = &UserAgent{
					Name:       h.UserAgent.Name,
					LastSeenAt: h.CreatedAt,
					ua:         detector.New(h.UserAgent.Name),
				}
			}
			if b.LastSeenAt.Before(h.CreatedAt) {
				b.LastSeenAt = h.CreatedAt
			}
			b.Count++
			out[*h.UserAgentHash] = b
		}
	}
	return out
}

func (ua UserAgent) IsBot() bool {
	return ua.ua.Bot()
}

func (ua UserAgent) IsMobile() bool {
	//return ua.ua.Mobile()
	return strings.Contains(ua.Name, "Mobi")
}

func (ua UserAgent) Platform() string {
	return ua.ua.Platform()
}

func (ua UserAgent) Browser() string {
	n, _ := ua.ua.Browser()
	return n
}

func (ua UserAgent) BrowserVersion() string {
	_, v := ua.ua.Browser()
	return v
}
