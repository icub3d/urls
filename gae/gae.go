// Package gae is an implementation of urls that works on app engine.
package gae

import (
	"github.com/icub3d/urls"
	"net/http"
)

func init() {
	http.HandleFunc("/admin/urls", urlsHandler)
	http.HandleFunc("/admin/urls/", urlHandler)
	http.HandleFunc("/admin/count/urls", getOrNotFound(urls.CountURLs))

	http.HandleFunc("/admin/stats", statsHandler)
	http.HandleFunc("/admin/stats/", getOrNotFound(urls.GetStatistics))

	http.HandleFunc("/", getOrNotFound(urls.Redirect))
}

// getOrNotFound is a helper function that returns a handle function
// that accepts GET request with the given handler or a not found.
func getOrNotFound(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			f(w, r)
		} else {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("not found"))
		}
	}
}

// urlsHandler handles the GET/POST for /admin/urls
func urlsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		urls.GetURLs(w, r)
	} else if r.Method == "POST" {
		urls.NewURL(w, r)
	} else {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("not found"))
	}
}

// urlsHandler handles the DELETE for /admin/urls/{id}
func urlHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "DELETE" {
		urls.DeleteURL(w, r)
	} else {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("not found"))
	}
}

// statsHandler handles the PUT for /admin/stats
func statsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "PUT" {
		urls.CreateStatistics(w, r)
	} else {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("not found"))
	}
}
