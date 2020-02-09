package sws

import (
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type Hit struct {
	DomainID *int    `json:"domain_id,omitempty"`
	Addr     *string `json:"addr,omitempty"`
	// URL components
	Scheme   *string `json:"scheme,omitempty"`
	Host     *string `json:"host,omitempty"`
	Path     *string `json:"page,omitempty"`
	Query    *string `json:"query,omitempty"`
	Fragment *string `json:"fragment,omitempty"`

	Title     *string `json:"title,omitempty"`
	Referrer  *string `json:"referrer,omitempty"`
	UserAgent *string `json:"user_agent,omitempty"`
	ViewPort  *string `json:"view_port,omitempty"`
	//Features  map[string]string `json:"features,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitempty"`

	// TODO
	Domain *Domain `json:"-"`
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
	if agent != "" {
		out.UserAgent = &agent
	} else {
		s := r.UserAgent()
		out.UserAgent = &s
	}

	if view := q.Get("v"); view != "" {
		out.ViewPort = &view
	}
	return &out, nil
}
