package main

import (
	"net/http"

	"src.userspace.com.au/sws"
)

func handleDomains(db sws.DomainStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		return
	}
}

func handleDomain(db sws.DomainStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		return
	}
}
