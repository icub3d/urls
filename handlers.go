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

// CreateStatistics is a handler function that updates the statistics
// of each URL. It checks each URL, but only updates the statistics
// with new log entries.
//
// This would normally map to something like PUT /stats. It does not
// check any session or admin cookies or anything like that. If you
// are checking those (and you probably should), you can wrap this
// handler in another handler.
func CreateStatistics(ds DataStore, w http.ResponseWriter, r *http.Request) {
	// These are the number of records we'll get per loop.
	const LIMIT = 1000

	offset := 0
	for {
		// Fetch the next set of URLs.
		urls, err := ds.GetURLs(LIMIT, offset)
		if err != nil {
			log.Printf(
				"GetURLs(100, %v) during CreateStatistics failed. Stopping with: %v",
				offset, err)
			return
		}

		if len(urls) == 0 {
			return
		}
		offset += LIMIT

		// Loop through each URL.
		for _, url := range urls {
			// Get the current stats for this url.
			stats, err := ds.GetStatistics(url.Short)
			if err != nil {
				log.Printf(
					"GetStatistics(%v) failed. Skipping this URL: %v",
					url.Short, err)
				return
			}

			// Create the maps if they weren't created.
			if stats.Referrers == nil {
				stats.Referrers = make(map[string]int)
			}
			if stats.Browsers == nil {
				stats.Browsers = make(map[string]int)
			}
			if stats.Countries == nil {
				stats.Countries = make(map[string]int)
			}
			if stats.Platforms == nil {
				stats.Platforms = make(map[string]int)
			}
			if stats.Hours == nil {
				stats.Hours = make(map[string]int)
			}

			// Now loop through all the Logs until we get to the last time
			// the log was updated.
			logOffset := 0
			var latest time.Time
			for {
				logs, err := ds.GetLogs(url.Short, LIMIT, logOffset)
				if err != nil {
					log.Printf(
						"GetLogs(%v, 100, %v) failed. Stopping: %v",
						url.Short, logOffset, err)
					return
				}
				logOffset += LIMIT

				if len(logs) == 0 {
					break
				}

				for _, log := range logs {
					if stats.LastUpdated.After(log.When) {
						// We have already processed this log.
						continue
					}

					if log.When.After(latest) {
						latest = log.When
					}

					browser, platform := parseUserAgent(log.UserAgent)
					country := determineCountry(log.Addr)
					hour := fmt.Sprintf("%04d%02d%02d%02d",
						log.When.Year(), log.When.Month(), log.When.Day(),
						log.When.Hour())

					// Update the values.
					stats.Referrers[log.Referrer] = stats.Referrers[log.Referrer] + 1
					stats.Browsers[browser] = stats.Browsers[browser] + 1
					stats.Countries[country] = stats.Countries[country] + 1
					stats.Platforms[platform] = stats.Platforms[platform] + 1
					stats.Hours[hour] = stats.Hours[hour] + 1

					// Update the clicks.
					stats.Clicks += 1
					url.Clicks += 1
				}
			}

			// Set the update time to the newest time.
			stats.LastUpdated = latest

			// Put the Url for the Clicks count.
			ds.PutURL(url)

			// Put the Statistics.
			ds.PutStatistics(stats)
		}
	}
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

	// Write the redirect.
	w.Header().Add("Location", u.Long)
	w.WriteHeader(http.StatusFound)
}
