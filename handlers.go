// Copyright 2013 Joshua Marsh. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package urls

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path"
	"time"
)

// HandlerFunc is a handler for the URL system.
type HandlerFunc func(ds DataStore, w http.ResponseWriter, r *http.Request)

// GetURLs is a handler func for getting a list of urls sorted by
// create date. If limit and offset are query parameters, they are
// used to limit the return set and offset from the beginning. Offset
// defaults to 0 and limit defaults to 20. The max offset is 100.
//
// This would normally map to something like GET /urls. It does not
// check any session or admin cookies or anything like that. If you
// are checking those (and you probably should), you can wrap this
// handler in another handler.
func GetURLs(ds DataStore, w http.ResponseWriter, r *http.Request) {
	// Get the query parameters.
	limit, offset := getLimitOffset(r.URL.Query())

	// Get the data.
	u, err := ds.GetURLs(limit, offset)
	if err != nil {
		log.Printf("GetUrls(%v, %v) failed with: %v", limit, offset, err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("oops"))
		return
	}

	marshalAndWrite(w, u)
}

// CountURLs is a handler func that returns the number of urls in the
// system. It returns json in the form: {"count":%v}.
//
// This would normally map to something like GET /count/urls. It does
// not check any session or admin cookies or anything like that. If
// you are checking those (and you probably should), you can wrap this
// handler in another handler.
func CountURLs(ds DataStore, w http.ResponseWriter, r *http.Request) {
	c, err := ds.CountURLs()
	if err != nil {
		log.Printf("CountUrls() failed with: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("oops"))
		return
	}

	w.Write([]byte(fmt.Sprintf(`{"count":%v}`, c)))

}

// NewURL creates a new URL based on the URL given as JSON. The short
// ID is created, the count is zeroed and the time is set to the
// current time. The updated URL is returned.
//
// This would normally map to something like POST /urls. It
// does not check any session or admin cookies or anything like
// that. If you are checking those (and you probably should), you can
// wrap this handler in another handler.
func NewURL(ds DataStore, w http.ResponseWriter, r *http.Request) {
	// Get the posted data.
	u := &URL{}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("ReadAll() failed on body: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("oops"))
		return
	}

	err = json.Unmarshal(body, u)
	if err != nil {
		log.Printf("Unmarshal() failed on body: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("oops"))
		return
	}

	// Set the fields.
	u.Clicks = 0
	u.Short = ""
	u.Created = time.Now()

	// Put the URL.
	_, err = ds.PutURL(u)
	if err != nil {
		log.Printf("PutURL(%v) failed on body: %v", u, err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("oops"))
		return
	}

	marshalAndWrite(w, u)
}

// DeleteURL deletes the url with the short id in the URL.
//
// This would normally map to something like DELETE /urls/{id}. It
// does not check any session or admin cookies or anything like
// that. If you are checking those (and you probably should), you can
// wrap this handler in another handler.
func DeleteURL(ds DataStore, w http.ResponseWriter, r *http.Request) {
	id := path.Base(r.URL.Path)

	if !ValidID(id) {
		// An invalid ID should return a not found.
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("not found"))
		return
	}

	err := ds.DeleteURL(id)
	if err != nil {
		if err != nil {
			log.Printf("GetUrl(%v) failed with: %v", id, err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("oops"))
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}

// GetStatistics is a handler func for getting the statistics of a URL.
//
// This would normally map to something like GET /stats/{id}. It does not
// check any session or admin cookies or anything like that. If you
// are checking those (and you probably should), you can wrap this
// handler in another handler.
func GetStatistics(ds DataStore, w http.ResponseWriter, r *http.Request) {
	id := path.Base(r.URL.Path)

	if !ValidID(id) {
		// An invalid ID should return a not found.
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("not found"))
		return
	}

	// Get the data.
	u, err := ds.GetStatistics(id)
	if err != nil {
		log.Printf("GetStatistics(%v) failed with: %v", id, err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("oops"))
		return
	}

	marshalAndWrite(w, u)
}

// Redirect is a handler func that handles the redirect. Given a short
// id, it sets the HTTP code to 302 and the Location header. If the
// short id isn't found, a 404 not found is returned.
//
// This would normally map to something like GET /{id}.
func Redirect(ds DataStore, w http.ResponseWriter, r *http.Request) {
	id := path.Base(r.URL.Path)

	if !ValidID(id) {
		// An invalid ID should return a not found.
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("not found"))
		return
	}

	// Get the URL in question.
	u, err := ds.GetURL(id)
	if err != nil {
		if err != nil {
			log.Printf("GetUrl(%v) failed with: %v", id, err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("oops"))
			return
		}
	}

	// Check for nil.
	if u == nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("not found"))
		return
	}

	// Create a Log entry.
	l := NewLog(id, r)
	err = ds.LogClick(l)
	if err != nil {
		// We shouldn't error out here but we should log it.
		log.Printf("LogClick(%v) failed (not likely recorded with: %v",
			l, err)
	}

	updateStats(ds, u, r)

	// Write the redirect.
	w.Header().Add("Location", u.Long)
	w.WriteHeader(http.StatusFound)
}
