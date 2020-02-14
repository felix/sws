package sws

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Hit struct {
	ID       *int    `json:"id"`
	DomainID *int    `json:"domain_id,omitempty"`
	Addr     *string `json:"addr,omitempty"`
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
	//Features  map[string]string `json:"features,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitempty"`

	// TODO
	Domain    *Domain    `json:"-"`
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
		host = ref.Host
	}
	out := Hit{
		Host:      &host,
		CreatedAt: ptrTime(time.Now()),
	}

	scheme := q.Get("s")
	if scheme != "" {
		out.Scheme = &scheme
	} else {
		out.Scheme = &ref.Scheme
	}

	path := q.Get("p")
	if path != "" {
		out.Path = &path
	} else {
		out.Path = &ref.Path
	}

	query := q.Get("q")
	if query != "" {
		out.Path = &query
	} else {
		out.Path = &ref.Path
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
	return &out, nil
}
