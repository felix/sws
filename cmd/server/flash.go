package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
)

type flashLvl string

type flashCtxKey string

const (
	flashInfo    flashLvl = "info"
	flashError   flashLvl = "error"
	flashWarn    flashLvl = "warn"
	flashSuccess flashLvl = "success"

	flashQueryKey = "_flash"
)

type flashMsg struct {
	Level   flashLvl
	Message string
}

func flashSet(r *http.Request, l flashLvl, m string) *http.Request {
	flashes := flashGet(r)
	if flashes == nil {
		flashes = make([]flashMsg, 0)
	}
	flashes = append(flashes, flashMsg{l, m})
	return r.WithContext(context.WithValue(r.Context(), flashCtxKey("flash"), flashes))
}

func flashGet(r *http.Request) []flashMsg {
	if msg := r.URL.Query().Get(flashQueryKey); msg != "" {
		if b, err := base64.RawURLEncoding.DecodeString(msg); err == nil {
			debug("found flash from query", string(b))
			var f []flashMsg
			if err = json.Unmarshal(b, &f); err == nil {
				return f
			}
		}
	}
	if f, ok := r.Context().Value(flashCtxKey("flash")).([]flashMsg); ok {
		debug("found flash from context", f)
		return f
	}

	return nil
}

func flashURL(r *http.Request, url string) string {
	f := flashGet(r)
	b, err := json.Marshal(f)
	if err != nil {
		return url
	}
	qs := base64.RawURLEncoding.EncodeToString(b)
	vals := r.URL.Query()
	vals.Set(flashQueryKey, qs)
	return fmt.Sprintf("%s?%s", url, vals.Encode())
}

func flashQuery(r *http.Request) *http.Request {
	f := flashGet(r)
	b, err := json.Marshal(f)
	if err != nil {
		return nil
	}
	qs := base64.RawURLEncoding.EncodeToString(b)
	vals := r.URL.Query()
	vals.Set("_flash", qs)
	r.URL.RawQuery = vals.Encode()
	return r
}
