package urls

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sort"
	"strconv"
	"testing"
	"time"
)

var datastore *mds

func TestGetUrls(t *testing.T) {
	prep()

	tests := []struct {
		limit  int
		offset int
		err    error
		when   int
		start  int
		end    int
	}{
		// Test beginning.
		{
			limit:  20,
			offset: 0,
			start:  0,
			end:    20,
		},

		// Test in the middle
		{
			limit:  25,
			offset: 150,
			start:  150,
			end:    175,
		},

		// Test end.
		{
			limit:  25,
			offset: 175,
			start:  175,
			end:    200,
		},

		// Test a failure.
		{
			limit:  25,
			offset: 175,
			err:    fmt.Errorf("failure"),
			when:   1,
		},
	}

	a := datastore.URLsArray()
	for k, test := range tests {
		if test.err != nil {
			datastore.SetError(test.err, test.when)
		}

		w := httptest.NewRecorder()
		r, _ := http.NewRequest("GET",
			fmt.Sprintf("http://localhost/admin/urls?limit=%v&offset=%v",
				test.limit, test.offset), nil)

		GetURLs(w, r)

		enc, _ := json.Marshal(a[test.start:test.end])

		if test.err != nil {
			enc = []byte("oops")
		}

		body := w.Body.Bytes()

		if !bytes.Equal(enc, body) {
			t.Errorf("Test %v: bodies not equal: expecting %v, got %v",
				k, string(enc), string(body))
		}
	}
}

func TestCountUrls(t *testing.T) {
	prep()

	tests := []struct {
		err      error
		when     int
		expected string
	}{
		// Test normal get.
		{
			expected: `{"count":200}`,
		},

		// Test an error.
		{
			err:      fmt.Errorf("failure"),
			when:     1,
			expected: "oops",
		},
	}

	for k, test := range tests {
		if test.err != nil {
			datastore.SetError(test.err, test.when)
		}

		w := httptest.NewRecorder()
		r, _ := http.NewRequest("GET", "http://localhost/admin/count/urls", nil)

		CountURLs(w, r)

		body := w.Body.String()

		if test.expected != body {
			t.Errorf("Test %v: bodies not equal: expecting %v, got %v",
				k, test.expected, body)
		}
	}
}

func TestDeleteURL(t *testing.T) {
	prep()

	tests := []struct {
		id       string
		code     int
		expected string
		err      error
		when     int
	}{
		// Test a normal delete.
		{
			id:       "1c",
			code:     http.StatusOK,
			expected: ``,
		},

		// Test an error
		{
			id:       "1d",
			code:     http.StatusInternalServerError,
			expected: `oops`,
			err:      fmt.Errorf("failure"),
			when:     1,
		},
	}

	for k, test := range tests {
		if test.err != nil {
			datastore.SetError(test.err, test.when)
		}

		w := httptest.NewRecorder()
		r, _ := http.NewRequest("GET", "http://localhost/"+test.id, nil)

		DeleteURL(w, r)

		body := w.Body.String()
		if test.expected != body {
			t.Errorf("Test %v: bodies not equal: expecting %v, got %v",
				k, test.expected, body)
		}

		if test.code != w.Code {
			t.Errorf("Test %v: codes not equal: expecting %v, got %v",
				k, test.code, w.Code)
		}
	}
}

func TestNewURL(t *testing.T) {
	prep()

	tests := []struct {
		long     string
		short    int
		expected string
		err      error
		when     int
	}{
		// Test a new.
		{
			long:  "http://test.new/1000.html",
			short: 1000,
		},

		// Test an error
		{
			long:     "http://test.new/blah.html",
			short:    1001,
			expected: `oops`,
			err:      fmt.Errorf("failure"),
			when:     1,
		},
	}

	for k, test := range tests {
		if test.err != nil {
			datastore.SetError(test.err, test.when)
		}

		var b bytes.Buffer
		b.Write([]byte(`{"Long":"` + test.long + `"}`))
		w := httptest.NewRecorder()
		r, _ := http.NewRequest("POST", "http://localhost/urls", &b)

		NewURL(w, r)

		if test.err != nil {
			body := w.Body.String()
			if test.expected != body {
				t.Errorf("Test %v: bodies not equal: expecting %v, got %v",
					k, test.expected, body)
			}
		} else {
			body := w.Body.Bytes()
			if !bytes.Contains(body, []byte(IntToShort(test.short))) {
				t.Errorf("Test %v: short id not found: looking for %v, got %v",
					k, IntToShort(test.short), string(body))

			}
		}
	}
}

func TestGetLogs(t *testing.T) {
	prep()

	tests := []struct {
		id     string
		limit  int
		offset int
		err    error
		when   int
		start  int
		end    int
	}{
		// Test beginning.
		{
			id:     "1c",
			limit:  20,
			offset: 0,
			start:  0,
			end:    20,
		},

		// Test in the middle
		{
			id:     "1c",
			limit:  25,
			offset: 25,
			start:  25,
			end:    50,
		},

		// Test a failure.
		{
			id:     "1c",
			limit:  25,
			offset: 50,
			err:    fmt.Errorf("failure"),
			when:   1,
		},
	}

	for k, test := range tests {
		a := datastore.LogsArray(test.id)

		if test.err != nil {
			datastore.SetError(test.err, test.when)
		}

		w := httptest.NewRecorder()
		r, _ := http.NewRequest("GET",
			fmt.Sprintf("http://localhost/admin/logs/%v?limit=%v&offset=%v",
				test.id, test.limit, test.offset), nil)

		GetLogs(w, r)

		enc, _ := json.Marshal(a[test.start:test.end])

		if test.err != nil {
			enc = []byte("oops")
		}

		body := w.Body.Bytes()

		if !bytes.Equal(enc, body) {
			t.Errorf("Test %v: bodies not equal: expecting %v, got %v",
				k, string(enc), string(body))
		}
	}
}

func TestCountLogs(t *testing.T) {
	prep()

	tests := []struct {
		id       string
		err      error
		when     int
		expected string
	}{
		// Test normal get.
		{
			id:       "1c",
			expected: `{"count":100}`,
		},

		// Test an error.
		{
			id:       "1c",
			err:      fmt.Errorf("failure"),
			when:     1,
			expected: "oops",
		},
	}

	for k, test := range tests {
		if test.err != nil {
			datastore.SetError(test.err, test.when)
		}

		w := httptest.NewRecorder()
		r, _ := http.NewRequest("GET",
			fmt.Sprintf("http://localhost/admin/count/logs/%v", test.id), nil)

		CountLogs(w, r)

		body := w.Body.String()

		if test.expected != body {
			t.Errorf("Test %v: bodies not equal: expecting %v, got %v",
				k, test.expected, body)
		}
	}
}

func TestRedirect(t *testing.T) {
	prep()

	tests := []struct {
		id       string
		location string
		code     int
		expected string
		err      error
		when     int
	}{
		// Test a normal get.
		{
			id:       "1c",
			location: "http://longurl.com/100.html",
			code:     http.StatusFound,
			expected: ``,
		},

		// Test a not found.
		{
			id:       "198djd81jd",
			location: "",
			code:     http.StatusNotFound,
			expected: `not found`,
		},

		// Test an error
		{
			id:       "123",
			location: "",
			code:     http.StatusInternalServerError,
			expected: `oops`,
			err:      fmt.Errorf("failure"),
			when:     1,
		},
	}

	for k, test := range tests {
		if test.err != nil {
			datastore.SetError(test.err, test.when)
		}

		w := httptest.NewRecorder()
		r, _ := http.NewRequest("GET", "http://localhost/"+test.id, nil)

		Redirect(w, r)

		body := w.Body.String()
		if test.expected != body {
			t.Errorf("Test %v: bodies not equal: expecting %v, got %v",
				k, test.expected, body)
		}

		if test.code == http.StatusFound {
			loc := w.HeaderMap.Get("Location")
			if test.location != loc {
				t.Errorf("Test %v: location not equal: expecting %v, got %v",
					k, test.location, loc)
			}
		}

		if test.code != w.Code {
			t.Errorf("Test %v: codes not equal: expecting %v, got %v",
				k, test.code, w.Code)
		}
	}
}

func prep() {
	datastore = &mds{
		urls:  make(map[string]*URL),
		stats: make(map[string]*Statistics),
		logs:  make(map[string][]*Log),
	}

	datastore.count = 1000

	t, _ := time.Parse("Jan 2 2006", "Jan 2 2013")

	// Add a set of urls.
	for x := 0; x < 200; x++ {
		u := &URL{
			Short:   IntToShort(x),
			Long:    "http://longurl.com/" + strconv.Itoa(x) + ".html",
			Created: t.AddDate(0, 0, -x),
			Clicks:  x,
		}

		datastore.PutURL(u)

		l := make([]*Log, 0, x)
		for y := 0; y < x; y++ {
			l = append(l, &Log{
				Short:     u.Short,
				When:      t.AddDate(0, 0, -y),
				Addr:      "127.0.0." + strconv.Itoa(y),
				Referrer:  "www.google.com",
				UserAgent: " Mozilla/5.0 (iPad; U; CPU OS 3_2_1 like Mac OS X; en-us) AppleWebKit/531.21.10 (KHTML, like Gecko) Mobile/7B405",
			})
		}

		datastore.logs[u.Short] = l

	}

	DS = datastore
}

// mds implements a DataStore in memory suitable for testing.
type mds struct {
	urls  map[string]*URL
	stats map[string]*Statistics
	logs  map[string][]*Log
	err   error
	when  int
	count int
}

// SetError marks an error to occur after when calls.
func (ds *mds) SetError(err error, when int) {
	ds.err = err
	ds.when = when
}

// Error is a helper function that keeps track of the error and when
// and returns it when it's been called a sufficient number of times.
func (ds *mds) error() error {
	if ds.err == nil {
		return nil
	}

	if ds.when == 1 {
		err := ds.err
		ds.err = nil
		return err
	}

	ds.when--
	return nil
}

// Implement the URL parts.
func (ds *mds) CountURLs() (int, error) {
	if err := ds.error(); err != nil {
		return 0, err
	}

	return len(ds.urls), nil
}

func (ds *mds) URLsArray() []*URL {
	// Get an array of the urls
	u := surls{}
	for _, v := range ds.urls {
		u = append(u, v)
	}

	sort.Sort(u)

	return u
}

func (ds *mds) GetURLs(limit, offset int) ([]*URL, error) {
	if err := ds.error(); err != nil {
		return nil, err
	}

	u := ds.URLsArray()

	if offset > len(u) {
		return []*URL{}, nil
	}

	if offset+limit > len(u) {
		limit = len(u)
	} else {
		limit = offset + limit
	}

	return u[offset:limit], nil
}

func (ds *mds) GetURL(short string) (*URL, error) {
	if err := ds.error(); err != nil {
		return nil, err
	}

	got := ds.urls[short]
	if got == nil {
		return nil, nil
	}

	return &URL{
		Short:   got.Short,
		Long:    got.Long,
		Created: got.Created,
		Clicks:  got.Clicks,
	}, nil
}

func (ds *mds) PutURL(url *URL) (string, error) {
	if err := ds.error(); err != nil {
		return "", err
	}

	if url.Short == "" {
		id := IntToShort(ds.count)
		ds.count++

		url.Short = id
	}

	ds.urls[url.Short] = &URL{
		Short:   url.Short,
		Long:    url.Long,
		Created: url.Created,
		Clicks:  url.Clicks,
	}

	return url.Short, nil
}

func (ds *mds) DeleteURL(short string) error {
	if err := ds.error(); err != nil {
		return err
	}

	delete(ds.urls, short)

	return nil
}

func (ds *mds) GetStatistics(short string) (*Statistics, error) {
	if err := ds.error(); err != nil {
		return nil, err
	}

	return ds.stats[short], nil
}

func (ds *mds) PutStatistics(stats *Statistics) error {
	if err := ds.error(); err != nil {
		return err
	}

	ds.stats[stats.Short] = stats
	return nil
}

func (ds *mds) LogClick(l *Log) error {
	if err := ds.error(); err != nil {
		return err
	}

	ls, found := ds.logs[l.Short]
	if !found {
		ls = []*Log{}
	}

	ds.logs[l.Short] = append(ls, l)

	return nil
}

func (ds *mds) CountLogs(short string) (int, error) {
	if err := ds.error(); err != nil {
		return 0, err
	}

	return len(ds.logs[short]), nil
}

func (ds *mds) LogsArray(id string) []*Log {
	// Get an array of the urls
	u := slogs{}
	for _, v := range ds.logs[id] {
		u = append(u, v)
	}

	sort.Sort(u)

	return u
}

func (ds *mds) GetLogs(short string, limit, offset int) ([]*Log, error) {
	if err := ds.error(); err != nil {
		return nil, err
	}

	// Get an array of the urls
	u := ds.LogsArray(short)

	if offset > len(u) {
		return []*Log{}, nil
	}

	if offset+limit > len(u) {
		limit = len(u)
	} else {
		limit = offset + limit
	}

	return u[offset:limit], nil
}

// These are sort helpers for the url and logs.
type surls []*URL

func (s surls) Len() int {
	return len(s)
}

func (s surls) Less(i, j int) bool {
	return s[i].Created.After(s[j].Created)
}

func (s surls) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

type slogs []*Log

func (s slogs) Len() int {
	return len(s)
}

func (s slogs) Less(i, j int) bool {
	return s[i].When.After(s[j].When)
}

func (s slogs) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
