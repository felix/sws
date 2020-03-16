package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"net"
	"net/http"
	"strings"
	"text/template"

	"github.com/hashicorp/golang-lru"
	"src.userspace.com.au/sws"
)

const (
	gif = "R0lGODlhAQABAIAAAAAAAP///yH5BAEAAAAALAAAAAABAAEAAAIBRAA7"
)

func handleHitCounter(db sws.CounterStore, mmdbPath string) http.HandlerFunc {
	gifBytes, err := base64.StdEncoding.DecodeString(gif)
	if err != nil {
		panic(err)
	}

	cache, err := lru.New(100)
	if err != nil {
		panic(err)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Cache-Control", "no-store")

		hit, err := sws.HitFromRequest(r)
		if err != nil {
			log("failed to extract hit", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		site, err := verifyHit(db, hit)
		if err != nil {
			log("failed to verify site", err)
			http.Error(w, "invalid site", http.StatusBadRequest)
			return
		}

		hit.Addr = r.RemoteAddr
		if strings.Contains(r.RemoteAddr, ":") {
			hit.Addr, _, err = net.SplitHostPort(r.RemoteAddr)
		}

		if r.Header.Get("X-Moz") == "prefetch" || r.Header.Get("X-Purpose") == "preview" {
			w.Header().Set("Content-Type", "image/gif")
			w.Write(gifBytes)
			return
		}

		// Ignore IPs
		if site.IgnoreIPs != "" && strings.Contains(site.IgnoreIPs, hit.Addr) {
			w.Header().Set("Content-Type", "image/gif")
			w.Write(gifBytes)
			return
		}

		if err == nil && hit.Addr != "" {
			var cc *string
			if v, ok := cache.Get(hit.Addr); ok {
				cc = v.(*string)
			} else if mmdbPath != "" {
				if cc, err = sws.FetchCountryCode(mmdbPath, hit.Addr); err != nil {
					log("geoip lookup failed:", err)
				}
				cache.Add(hit.Addr, cc)
			}
			hit.CountryCode = cc
			debug("geolocated:", hit.Addr, "to", *hit.CountryCode)
		}

		if err := db.SaveHit(hit); err != nil {
			log("failed to save hit", err)
			//http.Error(w, err.Error(), http.StatusInternalServerError)
			//return
		}
		w.Header().Set("Content-Type", "image/gif")
		w.Write(gifBytes)
		return
	}
}

func verifyHit(db sws.SiteGetter, h *sws.Hit) (*sws.Site, error) {
	if h.SiteID == nil {
		return nil, fmt.Errorf("invalid site ID")
	}
	site, err := db.GetSiteByID(*h.SiteID)
	if err != nil {
		return nil, err
	}
	if site.Name == h.Host {
		return site, nil
	}
	if strings.Contains(site.Aliases, h.Host) {
		return site, nil
	}
	if site.AcceptSubdomains && strings.HasSuffix(h.Host, site.Name) {
		return site, nil
	}
	return nil, fmt.Errorf("invalid host")
}

func handleCounter(addr string) http.HandlerFunc {
	counter := getCounter()
	tmpl, err := template.New("counter").Parse(counter)
	if err != nil || tmpl == nil {
		panic(err)
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, newTemplateData(nil)); err != nil {
		panic(err)
	}
	b := buf.Bytes()
	etag := fmt.Sprintf(`"%x"`, sha1.Sum(b))

	return func(w http.ResponseWriter, r *http.Request) {
		if match := r.Header.Get("If-None-Match"); match != "" {
			if strings.Contains(match, etag) {
				w.WriteHeader(http.StatusNotModified)
				return
			}
		}
		// TODO restrict to site sites
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Etag", etag)
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Content-Type", "application/javascript")

		if _, err := w.Write(b); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
