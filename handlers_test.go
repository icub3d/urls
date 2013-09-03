package urls

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"sort"
	"strconv"
	"testing"
	"time"
)

func TestGetUrls(t *testing.T) {
	ds := prep()

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

	a := ds.URLsArray()
	for k, test := range tests {
		if test.err != nil {
			ds.SetError(test.err, test.when)
		}

		w := httptest.NewRecorder()
		r, _ := http.NewRequest("GET",
			fmt.Sprintf("http://localhost/admin/urls?limit=%v&offset=%v",
				test.limit, test.offset), nil)

		GetURLs(ds, w, r)

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
	ds := prep()

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
			ds.SetError(test.err, test.when)
		}

		w := httptest.NewRecorder()
		r, _ := http.NewRequest("GET", "http://localhost/admin/count/urls", nil)

		CountURLs(ds, w, r)

		body := w.Body.String()

		if test.expected != body {
			t.Errorf("Test %v: bodies not equal: expecting %v, got %v",
				k, test.expected, body)
		}
	}
}

func TestDeleteURL(t *testing.T) {
	ds := prep()

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
			ds.SetError(test.err, test.when)
		}

		w := httptest.NewRecorder()
		r, _ := http.NewRequest("GET", "http://localhost/"+test.id, nil)

		DeleteURL(ds, w, r)

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
	ds := prep()

	tests := []struct {
		long     string
		short    int64
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
			ds.SetError(test.err, test.when)
		}

		var b bytes.Buffer
		b.Write([]byte(`{"Long":"` + test.long + `"}`))
		w := httptest.NewRecorder()
		r, _ := http.NewRequest("POST", "http://localhost/urls", &b)

		NewURL(ds, w, r)

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

func TestGetStatistics(t *testing.T) {
	ds := prep()

	tests := []struct {
		id       string
		err      error
		when     int
		expected string
	}{
		// Test in the middle
		{
			id:       "1c",
			expected: `{"Short":"1c","Clicks":100,"LastUpdated":"0001-01-01T00:00:00Z","Referrers":null,"Browsers":null,"Countries":null,"Platforms":null,"Hours":null}`,
		},

		// Test a failure.
		{
			id:       "1c",
			err:      fmt.Errorf("failure"),
			when:     1,
			expected: "oops",
		},
	}

	for k, test := range tests {
		if test.err != nil {
			ds.SetError(test.err, test.when)
		}

		w := httptest.NewRecorder()
		r, _ := http.NewRequest("GET",
			fmt.Sprintf("http://localhost/admin/stats/%v",
				test.id), nil)

		GetStatistics(ds, w, r)

		if w.Body.String() != test.expected {
			t.Errorf("Test %v: bodies not equal: expecting %v, got %v",
				k, test.expected, w.Body.String())
		}
	}
}

func TestCreateStatistics(t *testing.T) {
	td, _ := time.Parse("Jan 2 15:04:05 2006", "Aug 11 13:41:45 2013")

	// It's going to be easier to do this ourselves.
	ds := &mds{
		urls: map[string]*URL{
			"1c": &URL{
				Short:   "1c",
				Long:    "http://localhost/",
				Created: td,
				Clicks:  10,
			},
			"1d": &URL{
				Short:   "1d",
				Long:    "http://localhost/",
				Created: td,
				Clicks:  5,
			},
		},
		stats: map[string]*Statistics{
			"1c": &Statistics{
				Short:       "1c",
				Clicks:      10,
				LastUpdated: td.AddDate(0, 0, -5),
				Referrers: map[string]int{
					"localhost": 7,
					"127.0.0.1": 3,
				},
				Browsers: map[string]int{
					"Chrome": 7,
					"Safari": 3,
				},
				Countries: map[string]int{
					"US": 7,
					"AU": 3,
				},
				Platforms: map[string]int{
					"Linux":    7,
					"Mac OS X": 3,
				},
				Hours: map[string]int{
					"2013080111": 7,
					"2013080112": 3,
				},
			},
			"1d": &Statistics{
				Short:       "1d",
				Clicks:      5,
				LastUpdated: td.AddDate(0, 0, -5),
				Referrers: map[string]int{
					"localhost": 4,
					"127.0.0.1": 1,
				},
				Browsers: map[string]int{
					"Chrome": 3,
					"Safari": 2,
				},
				Countries: map[string]int{
					"US": 4,
					"AU": 1,
				},
				Platforms: map[string]int{
					"Linux":    1,
					"Mac OS X": 4,
				},
				Hours: map[string]int{
					"2013080111": 3,
					"2013080112": 2,
				},
			},
		},
		logs: map[string][]*Log{
			"1c": []*Log{
				&Log{
					Short:     "1c",
					When:      td.AddDate(0, 0, -7),
					Addr:      "1.0.0.45",
					Referrer:  "127.0.0.1",
					UserAgent: "Chrome Mac OS X",
				},
				&Log{
					Short:     "1c",
					When:      td.AddDate(0, 0, -2),
					Addr:      "1.0.0.45",
					Referrer:  "127.0.0.1",
					UserAgent: "Chrome Mac OS X",
				},
				&Log{
					Short:     "1c",
					When:      td.AddDate(0, 0, -1),
					Addr:      "1.0.0.43",
					Referrer:  "localhost",
					UserAgent: "Chrome Linux",
				},
			},
			"1d": []*Log{
				&Log{
					Short:     "1d",
					When:      td.AddDate(0, 0, -7),
					Addr:      "1.0.0.45",
					Referrer:  "127.0.0.1",
					UserAgent: "Chrome Mac OS X",
				},
				&Log{
					Short:     "1d",
					When:      td.AddDate(0, 0, -2),
					Addr:      "1.0.0.45",
					Referrer:  "127.0.0.1",
					UserAgent: "Chrome Mac OS X",
				},
				&Log{
					Short:     "1d",
					When:      td.AddDate(0, 0, -1),
					Addr:      "1.0.0.43",
					Referrer:  "localhost",
					UserAgent: "Chrome Linux",
				},
			},
		},
	}

	r, _ := http.NewRequest("GET", "/", nil)
	CreateStatistics(ds, httptest.NewRecorder(), r)

	onec := &Statistics{
		Short:       "1c",
		Clicks:      12,
		LastUpdated: td.AddDate(0, 0, -1),
		Referrers: map[string]int{
			"localhost": 8,
			"127.0.0.1": 4,
		},
		Browsers: map[string]int{
			"Chrome": 9,
			"Safari": 3,
		},
		Countries: map[string]int{
			"US": 7,
			"AU": 5,
		},
		Platforms: map[string]int{
			"Linux":    8,
			"Mac OS X": 4,
		},
		Hours: map[string]int{
			"2013080913": 1,
			"2013081013": 1,
			"2013080111": 7,
			"2013080112": 3,
		},
	}

	oned := &Statistics{
		Short:       "1d",
		Clicks:      7,
		LastUpdated: td.AddDate(0, 0, -1),
		Referrers: map[string]int{
			"localhost": 5,
			"127.0.0.1": 2,
		},
		Browsers: map[string]int{
			"Chrome": 5,
			"Safari": 2,
		},
		Countries: map[string]int{
			"US": 4,
			"AU": 3,
		},
		Platforms: map[string]int{
			"Linux":    2,
			"Mac OS X": 5,
		},
		Hours: map[string]int{
			"2013080913": 1,
			"2013081013": 1,
			"2013080111": 3,
			"2013080112": 2,
		},
	}

	if !reflect.DeepEqual(onec, ds.stats["1c"]) {
		t.Errorf("Onec not equal: %v %v", onec, ds.stats["1c"])
	}

	if !reflect.DeepEqual(oned, ds.stats["1d"]) {
		t.Errorf("Oned not equal: %v %v", oned, ds.stats["1d"])
	}
}

func TestRedirect(t *testing.T) {
	ds := prep()

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

		// Test a not valid.
		{
			id:       "this is invalid",
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
			ds.SetError(test.err, test.when)
		}

		w := httptest.NewRecorder()
		r, _ := http.NewRequest("GET", "http://localhost/"+test.id, nil)

		Redirect(ds, w, r)

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

func prep() *mds {
	ds := &mds{
		urls:  make(map[string]*URL),
		stats: make(map[string]*Statistics),
		logs:  make(map[string][]*Log),
	}

	ds.count = 1000

	t, _ := time.Parse("Jan 2 2006", "Jan 2 2013")

	// Add a set of urls.
	var x int64
	for x = 0; x < 200; x++ {
		u := &URL{
			Short:   IntToShort(x),
			Long:    "http://longurl.com/" + strconv.Itoa(int(x)) + ".html",
			Created: t.AddDate(0, 0, -int(x)),
			Clicks:  int(x),
		}

		ds.PutURL(u)

		ds.PutStatistics(&Statistics{
			Short:  u.Short,
			Clicks: int(x),
		})

		l := make([]*Log, 0, x)
		var y int64
		for y = 0; y < x; y++ {
			l = append(l, &Log{
				Short:     u.Short,
				When:      t.AddDate(0, 0, int(-x+y)),
				Addr:      "1.0.0.23",
				Referrer:  "www.google.com",
				UserAgent: "Mozilla/5.0 (Windows NT 6.1) Chrome/28.0.1500.95",
			})
		}

		ds.logs[u.Short] = l

	}

	return ds
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
		id := IntToShort(int64(ds.count))
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
	return s[j].When.After(s[i].When)
}

func (s slogs) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
