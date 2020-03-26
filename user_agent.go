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
	Browser    string    `json:"browser"`
	Platform   string    `json:"platform"`
	Version    string    `json:"version"`
	Bot        bool      `json:"bot"`
	Mobile     bool      `json:"mobile"`
	hitSet     *HitSet
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
func UserAgentFromRequest(r *http.Request) *UserAgent {
	q := r.URL.Query()
	ua := q.Get("u")
	if ua == "" {
		ua = r.UserAgent()
	}
	hash := UserAgentHash(ua)

	det := detector.New(ua)
	browser, version := det.Browser()

	return &UserAgent{
		Name:       ua,
		LastSeenAt: time.Now(),
		Hash:       hash,
		Browser:    browser,
		Platform:   det.Platform(),
		Version:    version,
		Bot:        det.Bot(),
		Mobile:     strings.Contains(ua, "Mobi") || det.Mobile(),
	}
}

func (ua UserAgent) Count() int {
	return ua.hitSet.Count()
}

func (ua UserAgent) Label() string {
	return ua.Browser
}
