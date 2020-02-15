package sws

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Hit struct {
	ID     *int    `json:"id"`
	SiteID *int    `json:"site_id,omitempty"`
	Addr   *string `json:"addr,omitempty"`
	// URL components
	Scheme   *string `json:"scheme,omitempty"`
	Host     *string `json:"host,omitempty"`
	Path     *string `json:"path,omitempty"`
	Query    *string `json:"query,omitempty"`
	Fragment *string `json:"fragment,omitempty"`

	Title         *string `json:"title,omitempty"`
	Referrer      *string `json:"referrer,omitempty"`
	UserAgentHash *string `json:"user_agent_hash,omitempty"`
	ViewPort      *string `json:"view_port,omitempty"`
	NoScript      bool    `json:"no_script"`
	//Features  map[string]string `json:"features,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitempty"`

	// TODO
	Site      *Site      `json:"-"`
	UserAgent *UserAgent `json:"-"`
}

func (h Hit) String() string {
	var out strings.Builder
	for _, sp := range []*string{h.Scheme, h.Host, h.Path, h.Query, h.Fragment} {
		if sp != nil {
			out.WriteString(*sp)
		}
	}
	return out.String()
}

func HitFromRequest(r *http.Request) (*Hit, error) {
	out := &Hit{
		CreatedAt: ptrTime(time.Now()),
	}

	q := r.URL.Query()
	host := q.Get("h")
	ref, err := url.ParseRequestURI(r.Referer())
	if err != nil || ref == nil {
		ref = new(url.URL)
	}
	if host == "" {
		if ref.Host == "" {
			return nil, fmt.Errorf("missing host")
		}
		out.NoScript = true
		host = ref.Host
	}
	out.Host = &host

	scheme := q.Get("s")
	if scheme != "" {
		out.Scheme = ptrString(strings.TrimSuffix(scheme, ":"))
	} else {
		out.Scheme = &ref.Scheme
	}

	path := q.Get("p")
	if path != "" {
		out.Path = &path
	} else {
		out.Path = &ref.RawPath
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
