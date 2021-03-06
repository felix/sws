package sws

import (
	"fmt"
	"net"
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
	CountryCode   *string   `json:"country_code" db:"country_code"`
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
	//Location() *time.Location // TODO Time zone
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
	// Strip port from remote address
	addr := r.RemoteAddr
	if strings.Contains(r.RemoteAddr, ":") {
		addr, _, _ = net.SplitHostPort(r.RemoteAddr)
	}
	out := &Hit{
		CreatedAt: time.Now(),
		Addr:      addr,
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

	out.UserAgent = UserAgentFromRequest(r)
	out.UserAgentHash = &out.UserAgent.Hash

	if view := q.Get("v"); view != "" {
		out.ViewPort = &view
	}

	return out, nil
}
