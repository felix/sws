package sws

import (
	"testing"
	"time"
)

func TestNewReferrerSet(t *testing.T) {
	now := time.Now()
	site := Site{Name: "example.com"}

	tests := []struct {
		hits     []*Hit
		expected ReferrerSet
	}{
		{
			hits: []*Hit{
				{CreatedAt: now, Referrer: strPtr("http://example1.com")},
				{CreatedAt: now, Referrer: strPtr("http://example1.com")},
				{CreatedAt: now, Referrer: strPtr("http://example2.com")},
			},
			expected: ReferrerSet{
				&Referrer{Name: "example1.com", URL: "http://example1.com"},
				&Referrer{Name: "example2.com", URL: "http://example2.com"},
			},
		},
		{
			hits: []*Hit{
				{CreatedAt: now, Referrer: strPtr("http://example1.com")},
				{CreatedAt: now, Referrer: strPtr("http://example1.com")},
				{CreatedAt: now, Referrer: nil},
			},
			expected: ReferrerSet{
				&Referrer{Name: "example1.com", URL: "http://example1.com"},
				&Referrer{Name: "direct", URL: ""},
			},
		},
		{
			hits: []*Hit{
				{CreatedAt: now, Referrer: strPtr("http://example1.com")},
				{CreatedAt: now, Referrer: strPtr("http://example1.com")},
				{CreatedAt: now, Referrer: strPtr("http://example.com")},
				{CreatedAt: now, Referrer: strPtr("http://example2.com")},
			},
			expected: ReferrerSet{
				&Referrer{Name: "example1.com", URL: "http://example1.com"},
				&Referrer{Name: "example2.com", URL: "http://example2.com"},
			},
		},
	}

	for i, tt := range tests {
		hs, err := NewHitSet(FromHits(tt.hits))
		if err != nil {
			t.Fatalf("%d => failed %s", i, err)
		}
		rs := NewReferrerSet(hs, site)

		if len(*rs) != len(tt.expected) {
			t.Errorf("%d => expected %d, got %d", i, len(tt.expected), len(*rs))
		}
		for j := range *rs {
			if (*rs)[j].Name != tt.expected[j].Name {
				t.Errorf("%d => expected %s, got %s", i, tt.expected[j].Name, (*rs)[j].Name)
			}
			if (*rs)[j].URL != tt.expected[j].URL {
				t.Errorf("%d => expected %s, got %s", i, tt.expected[j].URL, (*rs)[j].URL)
			}
		}
	}
}

func strPtr(s string) *string { return &s }
