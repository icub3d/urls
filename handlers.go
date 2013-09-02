package urls

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
)

var (
	// DS is the datastore the handlers will use. You should set this to
	// your implementation or expect errors.
	DS DataStore
)

// getInt is a helper function that returns the integer value of the
// query parameter with the given key.
func getInt(q url.Values, key string) int {
	value := q.Get(key)
	if value == "" {
		return 0
	}

	i, err := strconv.Atoi(value)
	if err != nil {
		return 0
	}

	return i
}

// GetUrls gets a list of urls sorted by create date. If limit and
// offset are query parameters, they are used to limit the return set
// and offset from the beginning. offset defaults to 0 and limit
// defaults to 20. The max offset is 100.
//
// This would normally map to something like GET /urls. It does not
// check any session or admin cookies or anything like that. If you
// are checking those (and you probably should, you can wrap this
// handler in another handler.
func GetUrls(w http.ResponseWriter, r *http.Request) {
	// Get the query parameters.
	q := r.URL.Query()
	limit := getInt(q, "limit")
	offset := getInt(q, "offset")

	// Set sane values if we don't find any.
	if limit <= 0 {
		limit = 20
	} else if limit > 100 {
		limit = 100
	}

	if offset < 0 {
		offset = 0
	}
}

// TODO implement CountUrls GET /count/urls
// returns {"count": count} as JSON

// TODO implement NewUrl POST /urls
// Expects POST data to be JSON of *Url without a short ID.
// The short ID will be blanked and a time will be created.
// returns *Url as JSON with fixed short ID

// TODO implement GetUrl GET /urls/{id}
// returns *Url as JSON if one is found or 404 not found.

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

// TODO implement Redirect /{id}
// if the id is found, issues a redirect. Otherwise 404 not found.
