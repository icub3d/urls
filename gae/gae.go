// Package gae is an implementation of urls that works on app engine.
package gae

import (
	"appengine"
	"appengine/user"
	"github.com/icub3d/urls"
	"net/http"
)

func init() {
	http.HandleFunc("/api/user", userHandler)
	http.HandleFunc("/api/urls", urlsHandler)
	http.HandleFunc("/api/urls/", urlHandler)
	http.HandleFunc("/api/count/urls", getOrNotFound(urls.CountURLs))

	http.HandleFunc("/api/stats/", getOrNotFound(urls.GetStatistics))

	http.HandleFunc("/", getOrNotFound(urls.Redirect))
}

// userHandler get the currently logged in user and returns their
// e-mail address.
func userHandler(w http.ResponseWriter, r *http.Request) {
	cxt := appengine.NewContext(r)
	if r.Method == "GET" {
		u := user.Current(cxt)
		lo, _ := user.LogoutURL(cxt, "/admin/")
		w.Write([]byte(`{"Email":"` + u.Email + `","LogoutURL":"` + lo + `"}`))
	} else {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("not found"))
	}
}

// getOrNotFound is a helper function that returns a handle function
// that accepts GET request with the given handler or a not found.
func getOrNotFound(f urls.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ds := NewDataStore(appengine.NewContext(r))
		if r.Method == "GET" {
			f(ds, w, r)
		} else {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("not found"))
		}
	}
}

// urlsHandler handles the GET/POST for /admin/urls
func urlsHandler(w http.ResponseWriter, r *http.Request) {
	ds := NewDataStore(appengine.NewContext(r))
	if r.Method == "GET" {
		urls.GetURLs(ds, w, r)
	} else if r.Method == "POST" {
		urls.NewURL(ds, w, r)
	} else {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("not found"))
	}
}

// urlsHandler handles the DELETE for /admin/urls/{id}
func urlHandler(w http.ResponseWriter, r *http.Request) {
	ds := NewDataStore(appengine.NewContext(r))

	if r.Method == "DELETE" {
		urls.DeleteURL(ds, w, r)
	} else {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("not found"))
	}
}
