package sws

import (
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Hit struct {
	ID     *int   `json:"id"`
	SiteID *int   `json:"site_id" db:"site_id"`
	Addr   string `json:"addr"`
	// URL components
	Scheme string  `json:"scheme"`
	Host   string  `json:"host"`
	Path   string  `json:"path"`
	Query  *string `json:"query,omitempty"`

	Title         *string   `json:"title,omitempty"`
	Referrer      *string   `json:"referrer,omitempty"`
	UserAgentHash *string   `json:"user_agent_hash,omitempty" db:"user_agent_hash"`
	ViewPort      *string   `json:"view_port,omitempty" db:"view_port"`
	NoScript      bool      `json:"no_script" db:"no_script"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	//Features  map[string]string `json:"features,omitempty"`

	// TODO
	//Site *Site `db:"s"`
	UserAgent *UserAgent `db:"ua"`
}

type Hitter interface {
	//Filter(FilterFunc) *HitSet
	Hits() []*Hit
	Begin() time.Time
	End() time.Time
	Duration() time.Duration
	Location() *time.Location
}

func (h Hit) String() string {
	var out strings.Builder
	for _, sp := range []*string{&h.Scheme, &h.Host, &h.Path, h.Query} {
		if sp != nil {
			out.WriteString(*sp)
		}
	}
	return out.String()
}

// SortHits in ascending order by time.
func SortHits(hits []*Hit) {
	sort.Slice(hits, func(i, j int) bool {
		return hits[i].CreatedAt.Before(hits[j].CreatedAt)
	})
}

func HitFromRequest(r *http.Request) (*Hit, error) {
	out := &Hit{
		CreatedAt: time.Now(),
		Addr:      r.RemoteAddr,
	}

	q := r.URL.Query()
	siteIDs := q.Get("id")
	if siteIDs == "" {
		if siteIDs = q.Get("site"); siteIDs == "" {
			return nil, fmt.Errorf("missing site")
		}
	}
	siteID, err := strconv.Atoi(siteIDs)
	if err != nil {
		return nil, fmt.Errorf("invalid site")
	}
	out.SiteID = &siteID

	// Host and referrer
	ref, err := url.ParseRequestURI(r.Referer())
	if err != nil || ref == nil {
		ref = new(url.URL)
	}
	host := q.Get("h")
	if host == "" {
		host = ref.Host
	}
	out.Host = host

	if h := r.Header.Get("HTTP_X_REQUESTED_WITH"); h == "" {
		out.NoScript = true
	}

	scheme := q.Get("s")
	if scheme != "" {
		out.Scheme = strings.TrimSuffix(scheme, ":")
	} else {
		out.Scheme = ref.Scheme
	}

	path := q.Get("p")
	if path != "" {
		out.Path = path
	} else {
		out.Path = ref.RawPath
	}

	query := q.Get("q")
	if query != "" {
		out.Query = &query
	} else {
		if ref.RawQuery != "" {
			out.Query = ptrString("?" + ref.RawQuery)
		}
	}

	if title := q.Get("t"); title != "" {
		out.Title = &title
	}

	referrer := q.Get("r")
	if referrer != "" {
		out.Referrer = &referrer
	} else {
		s := ref.String()
		out.Referrer = &s
	}

	agent := q.Get("u")
	if agent == "" {
		agent = r.UserAgent()
	}
	uaHash := UserAgentHash(agent)
	out.UserAgentHash = &uaHash
	out.UserAgent = &UserAgent{
		Hash:       uaHash,
		Name:       agent,
		LastSeenAt: time.Now(),
	}

	if view := q.Get("v"); view != "" {
		out.ViewPort = &view
	}
	return out, nil
}
