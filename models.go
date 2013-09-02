package urls

import (
	"net/http"
	"path"
	"time"
)

// Url is a representation of a shortened URL.
type Url struct {
	// The shortened ID.
	Short string

	// The full version of the URL.
	Long string

	// The date and time this Url was created.
	Created time.Time

	// The number of clicks this URL has received.
	Clicks int
}

// NewUrl creates a new url with the currend time and 0 clicks. The
// shortened url is not set.
func NewUrl(long string) *Url {
	return &Url{
		Long:    long,
		Created: time.Now(),
	}
}

// Log is a log of a click.
type Log struct {
	// The id of the URL clicked.
	Short string

	// The time the item was clicked.
	When time.Time

	// The ip address of the request.
	Addr string

	// The referrer of the request.
	Referrer string

	// The user agent of the request.
	UserAgent string
}

// NewLog creates a new log entry from the given request.
func NewLog(r *http.Request) *Log {
	return &Log{
		Short:     path.Base(r.URL.Path),
		When:      time.Now(),
		Addr:      r.RemoteAddr,
		Referrer:  r.Header.Get("Referer"),
		UserAgent: r.Header.Get("User-Agent"),
	}
}

// Statistics contain the information about the clicks a url has
// received.
type Statistics struct {
	// The name of the URL.
	Short string

	// The number of clicks this URL has received.
	Clicks int

	// A breakdown of the count by referrers. 'Unknown' is used for
	// clicks without a referrer.
	Referrers map[string]int

	// A breakdown of the count by browser. 'Unknown' is used for clicks
	// without a recognizable browser.
	Browsers map[string]int

	// A breakdown of the count by country. 'Unknown' is used for clicks
	// without a recognizable country.
	Countries map[string]int

	// A breakdown of the count by platform. 'Unknown' is used for clicks
	// without a recognizable platform.
	Platforms map[string]int

	// A breakdown of the count by Hours.
	Hours map[time.Time]int
}

// NewStatistics creates an empty set of statistics.
func NewStatistics(short string) *Statistics {
	return &Statistics{
		Short:     short,
		Referrers: make(map[string]int),
		Browsers:  make(map[string]int),
		Countries: make(map[string]int),
		Platforms: make(map[string]int),
		Hours:     make(map[time.Time]int),
	}
}
