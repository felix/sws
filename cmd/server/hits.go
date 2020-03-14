package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"text/template"

	"github.com/hashicorp/golang-lru"
	maxminddb "github.com/oschwald/maxminddb-golang"
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
		hit, err := sws.HitFromRequest(r)
		if err != nil {
			log("failed to extract hit", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		site, err := db.GetSiteByName(hit.Host)
		if err != nil {
			log("failed to get site", err)
			http.Error(w, "invalid site", http.StatusNotFound)
			return
		}
		hit.SiteID = site.ID
		hit.Addr = r.RemoteAddr

		host, _, err := net.SplitHostPort(addr)
		if err == nil {
			var cc *string
			if v, ok := cache.Get(host); ok {
				cc = v.(*string)
			} else if mmdbPath != "" {
				cc, _ = fetchCountryCode(mmdbPath, host)
				if cc != nil {
					debug("geoip", host, "=>", *cc)
				}
				cache.Add(host, cc)
			}
			hit.CountryCode = cc
		}

		if err := db.SaveHit(hit); err != nil {
			log("failed to save hit", err)
			//http.Error(w, err.Error(), http.StatusInternalServerError)
			//return
		}
		// TODO restrict to site sites
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Content-Type", "image/gif")
		w.Write(gifBytes)
		return
	}
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
	etag := fmt.Sprintf(`"%x"`, sha1.Sum(buf.Bytes()))

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
		w.Header().Set("Content-Type", "application/javascript")

		if _, err := io.Copy(w, &buf); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func fetchCountryCode(path, host string) (*string, error) {
	db, err := maxminddb.Open(path)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	ip := net.ParseIP(host)
	var r struct {
		Country struct {
			ISOCode string `maxminddb:"iso_code"`
		} `maxminddb:"country"`
	}
	if err := db.Lookup(ip, &r); err != nil {
		return nil, err
	}
	return &r.Country.ISOCode, nil
}
