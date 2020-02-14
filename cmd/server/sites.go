package main

import (
	"net/http"

	"src.userspace.com.au/sws"
)

func handleSites(db sws.SiteStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		return
	}
}

func handleSite(db sws.SiteStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		return
	}
}
