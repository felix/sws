// +build ignore

package sws

import (
	"net/http"
	"testing"
)

func TestHitFromRequest(t *testing.T) {
	tests := []struct {
		in       string
		expected Hit
	}{
		{
			`http://example.com/?h=example.com&s=http&i=1&p=/&t=title`,
			Hit{
				SiteID:   ptrInt(1),
				Scheme:   ptrString("http"),
				Host:     ptrString("example.com"),
				Path:     ptrString("/"),
				Title:    ptrString("title"),
				Query:    nil,
				Fragment: nil,
			},
		},
	}

	for i, tt := range tests {
		r, err := http.NewRequest("GET", tt.in, nil)
		if err != nil {
			t.Fatalf("%d => failed: %s", i, err)
		}
		actual, err := HitFromRequest(r)
		if err != nil {
			t.Fatalf("%d => failed: %s", i, err)
		}

		if tt.expected.SiteID != nil && *actual.SiteID != *tt.expected.SiteID {
			t.Errorf("%d => expected %d, got %d", i, *tt.expected.SiteID, *actual.SiteID)
		}
		if tt.expected.Scheme != nil && *actual.Scheme != *tt.expected.Scheme {
			t.Errorf("%d => expected %q, got %q", i, *tt.expected.Scheme, *actual.Scheme)
		}
		if tt.expected.Host != nil && *actual.Host != *tt.expected.Host {
			t.Errorf("%d => expected %q, got %q", i, *tt.expected.Host, *actual.Host)
		}
		if tt.expected.Path != nil && *actual.Path != *tt.expected.Path {
			t.Errorf("%d => expected %q, got %q", i, *tt.expected.Path, *actual.Path)
		}
		if tt.expected.Query != nil && *actual.Query != *tt.expected.Query {
			t.Errorf("%d => expected %q, got %q", i, *tt.expected.Query, *actual.Query)
		}
		if tt.expected.Fragment != nil && *actual.Fragment != *tt.expected.Fragment {
			t.Errorf("%d => expected %q, got %q", i, *tt.expected.Fragment, *actual.Fragment)
		}
	}
}
