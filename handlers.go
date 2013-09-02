package urls

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"path"
)

var (
	// DS is the datastore the handlers will use. You should set this to
	// your implementation or expect errors.
	DS DataStore
)

// GetUrls is a handler func for getting a list of urls sorted by
// create date. If limit and offset are query parameters, they are
// used to limit the return set and offset from the beginning. Offset
// defaults to 0 and limit defaults to 20. The max offset is 100.
//
// This would normally map to something like GET /urls. It does not
// check any session or admin cookies or anything like that. If you
// are checking those (and you probably should), you can wrap this
// handler in another handler.
func GetUrls(w http.ResponseWriter, r *http.Request) {
	// Get the query parameters.
	q := r.URL.Query()
	limit := paramGetInt(q, "limit")
	offset := paramGetInt(q, "offset")

	// Set sane values if we don't find any.
	if limit <= 0 {
		limit = 20
	} else if limit > 100 {
		limit = 100
	}

	if offset < 0 {
		offset = 0
	}

	// Get the data.
	u, err := DS.GetURLs(limit, offset)
	if err != nil {
		log.Printf("GetUrls(%v, %v) failed with: %v", limit, offset, err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("oops"))
		return
	}

	// Marshal it to JSON.
	enc, err := json.Marshal(u)
	if err != nil {
		log.Printf("Marshal(%v) failed with: %v", u, err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("oops"))
		return
	}

	// Write the response.
	w.Write(enc)
}

// CountUrls is a handler func that returns the number of urls in the
// system. It returns json in the form: {"count":%v}.
//
// This would normally map to something like GET /count/urls. It does
// not check any session or admin cookies or anything like that. If
// you are checking those (and you probably should), you can wrap this
// handler in another handler.
func CountUrls(w http.ResponseWriter, r *http.Request) {
	c, err := DS.CountURLs()
	if err != nil {
		log.Printf("CountUrls() failed with: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("oops"))
		return
	}

	w.Write([]byte(fmt.Sprintf(`{"count":%v}`, c)))

}

// TODO implement NewUrl POST /urls
// Expects POST data to be JSON of *Url without a short ID.
// The short ID will be blanked and a time will be created.
// returns *Url as JSON with fixed short ID

// TODO implement DeleteUrl DELETE /urls/{id}
// returns 200 OK or 404 Not Found

// TODO implement GetLogs GET /logs/{id}
// limit defaults to 20 max is 100
// offset defaults to 0
// returns []*Log as JSON

// TODO implement CountLogs GET /count/logs/{id}
// returns {"count": count} as JSON

// TODO implement GetStatistics GET /stats/{id}
// returns *Statistics as JSON.

// TODO implement CreateStatistics PUT /stats/
// Re-analyze the statistics.

// Redirect is a handler func that handles the redirect. Given a short
// id, it sets the HTTP code to 302 and the Location header. If the
// short id isn't found, a 404 not found is returned.
//
// This would normally map to something like GET /{id}.
func Redirect(w http.ResponseWriter, r *http.Request) {
	id := path.Base(r.URL.Path)

	u, err := DS.GetURL(id)
	if err != nil {
		if err != nil {
			log.Printf("GetUrl(%v) failed with: %v", id, err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("oops"))
			return
		}
	}

	if u != nil {
		w.Header().Add("Location", u.Long)
		w.WriteHeader(http.StatusFound)
		return
	}

	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("not found"))

}
