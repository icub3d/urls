package urls

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestValidID(t *testing.T) {
	tests := []struct {
		id       string
		expected bool
	}{
		// Test an empty string.
		{
			id:       "",
			expected: false,
		},

		// Test one with an invalid character.
		{
			id:       "98Nb-a",
			expected: false,
		},

		// Test a valid id.
		{
			id:       "9A8b1hdkalKAHgkalhg",
			expected: true,
		},
	}

	for k, test := range tests {
		result := ValidID(test.id)
		if result != test.expected {
			t.Errorf("Test %v: expected %v from ValidID(%v), but got %v",
				k, test.expected, test.id, result)
		}
	}
}

func TestIntToShort(t *testing.T) {

	tests := []struct {
		id       int64
		expected string
	}{
		// Test some known cases.
		{
			id:       -1,
			expected: "0",
		},
		{
			id:       0,
			expected: "0",
		},

		// Test some random cases.
		{
			id:       25883599,
			expected: "1kbVP",
		},
		{
			id:       13992492,
			expected: "wi5M",
		},
		{
			id:       24025617,
			expected: "1co9x",
		},
		{
			id:       10702,
			expected: "2mc",
		},
		{
			id:       56351835,
			expected: "3oRgZ",
		},
		{
			id:       60799965,
			expected: "476qb",
		},
		{
			id:       30313280,
			expected: "23Brs",
		},
		{
			id:       11612025,
			expected: "mioj",
		},
		{
			id:       19568944,
			expected: "1K6m8",
		},
		{
			id:       23213121,
			expected: "1ZOnB",
		},
		{
			id:       36808014,
			expected: "2URRe",
		},
		{
			id:       73173873,
			expected: "4x1rl",
		},
		{
			id:       33140829,
			expected: "2F3RV",
		},
		{
			id:       49448467,
			expected: "3LTnv",
		},
		{
			id:       23215655,
			expected: "1ZPS3",
		},
		{
			id:       31857603,
			expected: "29fcJ",
		},
		{
			id:       83196636,
			expected: "5d5FE",
		},
		{
			id:       64932637,
			expected: "4ORwb",
		},
		{
			id:       62056666,
			expected: "4CNly",
		},
		{
			id:       53483481,
			expected: "3cPUn",
		},
		{
			id:       6655522,
			expected: "RvP8",
		},
		{
			id:       20680798,
			expected: "1Om1G",
		},
		{
			id:       41750143,
			expected: "2pB7P",
		},
		{
			id:       3023948,
			expected: "CgfM",
		},
		{
			id:       10194536,
			expected: "gm40",
		},
		{
			id:       60294934,
			expected: "44zSw",
		},
		{
			id:       1985607,
			expected: "8KXv",
		},
		{
			id:       89591901,
			expected: "63uwf",
		},
		{
			id:       21746119,
			expected: "1TF9r",
		},
		{
			id:       66617413,
			expected: "4VWEP",
		},
		{
			id:       14147405,
			expected: "xMNx",
		},
		{
			id:       50948628,
			expected: "3Rm44",
		},
		{
			id:       67471159,
			expected: "4Z6KV",
		},
		{
			id:       57769279,
			expected: "3uOQZ",
		},
		{
			id:       80925232,
			expected: "5TYLg",
		},
		{
			id:       51234358,
			expected: "3SyOc",
		},
		{
			id:       21762868,
			expected: "1TJW0",
		},
		{
			id:       85775710,
			expected: "5nuBC",
		},
		{
			id:       2068155,
			expected: "8g1L",
		},
		{
			id:       38720804,
			expected: "2cT36",
		},
		{
			id:       40975171,
			expected: "2lvVr",
		},
		{
			id:       75084343,
			expected: "552rn",
		},
		{
			id:       69963207,
			expected: "4jYcp",
		},
		{
			id:       68132074,
			expected: "4bsGQ",
		},
		{
			id:       77810961,
			expected: "5GUBV",
		},
		{
			id:       49361381,
			expected: "3L79J",
		},
		{
			id:       68236415,
			expected: "4cJPL",
		},
		{
			id:       67175877,
			expected: "4XrVt",
		},
		{
			id:       65592921,
			expected: "4RDiL",
		},
		{
			id:       41238522,
			expected: "2n21S",
		},
		{
			id:       94283385,
			expected: "6NbPl",
		},
		{
			id:       41956665,
			expected: "2q2qP",
		},
		{
			id:       37707752,
			expected: "2YDVY",
		},
		{
			id:       68610149,
			expected: "4dsdJ",
		},
		{
			id:       28209392,
			expected: "1uMYC",
		},
		{
			id:       38691855,
			expected: "2cLWB",
		},
		{
			id:       99522769,
			expected: "6jaPx",
		},
		{
			id:       69682314,
			expected: "4iNYI",
		},
		{
			id:       16162743,
			expected: "15ofP",
		},
		{
			id:       38386419,
			expected: "2b43n",
		},
		{
			id:       42666958,
			expected: "2t1ck",
		},
		{
			id:       98849355,
			expected: "6glER",
		},
		{
			id:       17320289,
			expected: "1AfnV",
		},
		{
			id:       25808476,
			expected: "1kHxk",
		},
		{
			id:       91304382,
			expected: "6B6RK",
		},
		{
			id:       33144228,
			expected: "2F4KK",
		},
		{
			id:       61813865,
			expected: "4BMbp",
		},
		{
			id:       90135250,
			expected: "66CIM",
		},
		{
			id:       84605551,
			expected: "5izlf",
		},
		{
			id:       75446667,
			expected: "56Z7j",
		},
		{
			id:       35347002,
			expected: "2OJMw",
		},
		{
			id:       84988268,
			expected: "5kbKW",
		},
		{
			id:       344523,
			expected: "1Rcp",
		},
		{
			id:       34542600,
			expected: "2Kw6i",
		},
		{
			id:       11579636,
			expected: "maOK",
		},
		{
			id:       50145891,
			expected: "3OPEh",
		},
		{
			id:       97880514,
			expected: "6chBy",
		},
		{
			id:       44866383,
			expected: "32FnL",
		},
		{
			id:       29717072,
			expected: "20glc",
		},
		{
			id:       39625641,
			expected: "2gGRF",
		},
		{
			id:       9977484,
			expected: "frbA",
		},
		{
			id:       38415131,
			expected: "2bBWt",
		},
		{
			id:       25011862,
			expected: "1gwj8",
		},
		{
			id:       4742753,
			expected: "Jto1",
		},
		{
			id:       30922277,
			expected: "25kIP",
		},
		{
			id:       18561924,
			expected: "1Fsns",
		},
		{
			id:       86782082,
			expected: "5s7z0",
		},
		{
			id:       42781204,
			expected: "2tVLQ",
		},
		{
			id:       21860274,
			expected: "1Tir4",
		},
		{
			id:       82510071,
			expected: "5aCdb",
		},
		{
			id:       37078850,
			expected: "2VZty",
		},
		{
			id:       55022721,
			expected: "3irvF",
		},
		{
			id:       51293216,
			expected: "3TDhw",
		},
		{
			id:       90557044,
			expected: "67y1U",
		},
		{
			id:       94867994,
			expected: "6Q3Uw",
		},
		{
			id:       26714036,
			expected: "1o5XY",
		},
		{
			id:       11013256,
			expected: "kD3A",
		},
		{
			id:       37169131,
			expected: "2VxO7",
		},
		{
			id:       7988373,
			expected: "XW8j",
		},
		{
			id:       52386136,
			expected: "3Xo1g",
		},
	}

	for k, test := range tests {
		result := IntToShort(test.id)
		if result != test.expected {
			t.Errorf("Test %v: expected %v from IntToShort(%v), but got %v",
				k, test.expected, test.id, result)
		}
	}
}

func TestDigit(t *testing.T) {
	tests := []struct {
		id       int64
		expected string
	}{
		// Test values in int range.
		{
			id:       0,
			expected: "0",
		},
		{
			id:       5,
			expected: "5",
		},
		{
			id:       9,
			expected: "9",
		},

		// Test values in A-Z range.
		{
			id:       10,
			expected: "A",
		},
		{
			id:       25,
			expected: "P",
		},
		{
			id:       35,
			expected: "Z",
		},

		// Test values in a-z range.
		{
			id:       36,
			expected: "a",
		},
		{
			id:       51,
			expected: "p",
		},
		{
			id:       61,
			expected: "z",
		},

		// Test an invalid id.
		{
			id:       62,
			expected: "",
		},
	}

	for k, test := range tests {
		result := digit(test.id)
		if result != test.expected {
			t.Errorf("Test %v: expected %v from digit(%v), but got %v",
				k, test.expected, test.id, result)
		}
	}

}

func TestChar(t *testing.T) {
	tests := []struct {
		result int64
		id     string
	}{
		// Test values in int range.
		{
			result: 0,
			id:     "0",
		},
		{
			result: 5,
			id:     "5",
		},
		{
			result: 9,
			id:     "9",
		},

		// Test values in A-Z range.
		{
			result: 10,
			id:     "A",
		},
		{
			result: 25,
			id:     "P",
		},
		{
			result: 35,
			id:     "Z",
		},

		// Test values in a-z range.
		{
			result: 36,
			id:     "a",
		},
		{
			result: 51,
			id:     "p",
		},
		{
			result: 61,
			id:     "z",
		},

		// Test an invalid id.
		{
			result: 0,
			id:     "-",
		},

		{
			result: 0,
			id:     "",
		},
	}

	for k, test := range tests {
		result := char(test.id)
		if result != test.result {
			t.Errorf("Test %v: expected %v from char(%v), but got %v",
				k, test.result, test.id, result)
		}
	}

}

func TestParamGetInt(t *testing.T) {
	tests := []struct {
		q        url.Values
		key      string
		expected int
	}{
		// Test a missing key.
		{
			q: url.Values{
				"test": []string{},
			},
			key:      "notfound",
			expected: 0,
		},

		// Test a non-integer value.
		{
			q: url.Values{
				"test": []string{"nonint"},
			},
			key:      "test",
			expected: 0,
		},

		// Test a normal value.
		{
			q: url.Values{
				"test": []string{"54818"},
			},
			key:      "test",
			expected: 54818,
		},
	}
	for k, test := range tests {
		result := paramGetInt(test.q, test.key)
		if result != test.expected {
			t.Errorf("Test %v: expected %v from paramGetInt(%v, %v), but got %v",
				k, test.expected, test.q, test.key, result)
		}
	}

}

func TestGetLimitOffset(t *testing.T) {
	tests := []struct {
		q      url.Values
		limit  int
		offset int
	}{
		// Test no limit or offset
		{
			q:      url.Values{},
			limit:  20,
			offset: 0,
		},

		// Test values too small
		{
			q: url.Values{
				"limit":  []string{"0"},
				"offset": []string{"-1"},
			},
			limit:  20,
			offset: 0,
		},

		// Test a limit over 100
		{
			q: url.Values{
				"limit":  []string{"101"},
				"offset": []string{"2"},
			},
			limit:  100,
			offset: 2,
		},

		// Test a normal get.
		{
			q: url.Values{
				"limit":  []string{"55"},
				"offset": []string{"23"},
			},
			limit:  55,
			offset: 23,
		},
	}
	for k, test := range tests {
		limit, offset := getLimitOffset(test.q)
		if limit != test.limit || offset != test.offset {
			t.Errorf(
				"Test %v: expected (%v,%v) from getLimitOffset(%v), but got (%v,%v)",
				k, test.limit, test.offset, test.q, limit, offset)
		}
	}

}

func TestMarshalAndWrite(t *testing.T) {
	tests := []struct {
		i        interface{}
		code     int
		expected string
	}{
		// Test a normal write.
		{
			i:        123,
			code:     http.StatusOK,
			expected: "123",
		},

		// Test a failed marshal.
		{
			i:        complex(1, 1),
			code:     http.StatusInternalServerError,
			expected: "oops",
		},
	}

	for k, test := range tests {
		w := httptest.NewRecorder()

		marshalAndWrite(w, test.i)

		if w.Code != test.code {
			t.Errorf("Test %v: expected code %v but got %v", k, test.code, w.Code)
		}

		if w.Body.String() != test.expected {
			t.Errorf("Test %v: expected body '%v' but got '%v'",
				k, test.expected, w.Body.String())
		}
	}
}

func TestParseUserAgent(t *testing.T) {
	tests := []struct {
		ua       string
		browser  string
		platform string
	}{
		{
			ua:       "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/28.0.1500.95 Safari/537.36",
			browser:  "Chrome",
			platform: "Windows 7",
		},
		{
			ua:       "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_8_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/28.0.1500.95 Safari/537.36",
			browser:  "Chrome",
			platform: "Mac OS X",
		},
		{
			ua:       "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_8_4) AppleWebKit/536.30.1 (KHTML, like Gecko) Version/6.0.5 Safari/536.30.1",
			browser:  "Safari",
			platform: "Mac OS X",
		},
		{
			ua:       "Mozilla/5.0 (Windows NT 6.1; WOW64; rv:23.0) Gecko/20100101 Firefox/23.0",
			browser:  "Firefox",
			platform: "Windows 7",
		},
		{
			ua:       "Mozilla/5.0 (Windows NT 6.2; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/28.0.1500.95 Safari/537.36",
			browser:  "Chrome",
			platform: "Windows 8",
		},
		{
			ua:       "Opera/9.80 (Windows NT 6.1; WOW64) Presto/2.12.388 Version/12.16",
			browser:  "Opera",
			platform: "Windows 7",
		},
		{
			ua:       "",
			browser:  "Unknown",
			platform: "Unknown",
		},
	}

	for k, test := range tests {
		browser, platform := parseUserAgent(test.ua)

		if test.browser != browser {
			t.Errorf("Test %v: expected browser '%v' but got '%v': %v",
				k, test.browser, browser, test.ua)
		}

		if test.platform != platform {
			t.Errorf("Test %v: expected platform '%v' but got '%v': %v",
				k, test.platform, platform, test.ua)
		}
	}
}

func TestDetermineCountry(t *testing.T) {
	tests := []struct {
		addr    string
		country string
	}{
		// Test some invalid values
		{
			addr:    "abcdef",
			country: "Unknown",
		},
		{
			addr:    "10.84.8.abc",
			country: "Unknown",
		},
		{
			addr:    "10.84.abc.8",
			country: "Unknown",
		},
		{
			addr:    "10.aby.84.8",
			country: "Unknown",
		},
		{
			addr:    "abc.10.94.8",
			country: "Unknown",
		},
		{
			addr:    "127.0.0.1",
			country: "Unknown",
		},

		// Test a few good values
		{
			addr:    "1.0.0.23",
			country: "AU",
		},
		{
			addr:    "223.255.255.201",
			country: "AU",
		},
		{
			addr:    "206.251.44.94",
			country: "US",
		},
		{
			addr:    "190.109.96.35",
			country: "CO",
		},
	}

	for k, test := range tests {
		country := determineCountry(test.addr)
		if country != test.country {
			t.Errorf("Test %v: expected country '%v' but got '%v': %v",
				k, test.country, country, test.addr)
		}
	}
}

func TestShortToInt(t *testing.T) {

	tests := []struct {
		results int64
		id      string
	}{
		// Test some known cases.
		{
			results: 0,
			id:      "0",
		},
		{
			results: 0,
			id:      "0",
		},

		// Test some random cases.
		{
			results: 25883599,
			id:      "1kbVP",
		},
		{
			results: 13992492,
			id:      "wi5M",
		},
		{
			results: 24025617,
			id:      "1co9x",
		},
		{
			results: 10702,
			id:      "2mc",
		},
		{
			results: 56351835,
			id:      "3oRgZ",
		},
		{
			results: 60799965,
			id:      "476qb",
		},
		{
			results: 30313280,
			id:      "23Brs",
		},
		{
			results: 11612025,
			id:      "mioj",
		},
		{
			results: 19568944,
			id:      "1K6m8",
		},
		{
			results: 23213121,
			id:      "1ZOnB",
		},
		{
			results: 36808014,
			id:      "2URRe",
		},
		{
			results: 73173873,
			id:      "4x1rl",
		},
		{
			results: 33140829,
			id:      "2F3RV",
		},
		{
			results: 49448467,
			id:      "3LTnv",
		},
		{
			results: 23215655,
			id:      "1ZPS3",
		},
		{
			results: 31857603,
			id:      "29fcJ",
		},
		{
			results: 83196636,
			id:      "5d5FE",
		},
		{
			results: 64932637,
			id:      "4ORwb",
		},
		{
			results: 62056666,
			id:      "4CNly",
		},
		{
			results: 53483481,
			id:      "3cPUn",
		},
		{
			results: 6655522,
			id:      "RvP8",
		},
		{
			results: 20680798,
			id:      "1Om1G",
		},
		{
			results: 41750143,
			id:      "2pB7P",
		},
		{
			results: 3023948,
			id:      "CgfM",
		},
		{
			results: 10194536,
			id:      "gm40",
		},
		{
			results: 60294934,
			id:      "44zSw",
		},
		{
			results: 1985607,
			id:      "8KXv",
		},
		{
			results: 89591901,
			id:      "63uwf",
		},
		{
			results: 21746119,
			id:      "1TF9r",
		},
		{
			results: 66617413,
			id:      "4VWEP",
		},
		{
			results: 14147405,
			id:      "xMNx",
		},
		{
			results: 50948628,
			id:      "3Rm44",
		},
		{
			results: 67471159,
			id:      "4Z6KV",
		},
		{
			results: 57769279,
			id:      "3uOQZ",
		},
		{
			results: 80925232,
			id:      "5TYLg",
		},
		{
			results: 51234358,
			id:      "3SyOc",
		},
		{
			results: 21762868,
			id:      "1TJW0",
		},
		{
			results: 85775710,
			id:      "5nuBC",
		},
		{
			results: 2068155,
			id:      "8g1L",
		},
		{
			results: 38720804,
			id:      "2cT36",
		},
		{
			results: 40975171,
			id:      "2lvVr",
		},
		{
			results: 75084343,
			id:      "552rn",
		},
		{
			results: 69963207,
			id:      "4jYcp",
		},
		{
			results: 68132074,
			id:      "4bsGQ",
		},
		{
			results: 77810961,
			id:      "5GUBV",
		},
		{
			results: 49361381,
			id:      "3L79J",
		},
		{
			results: 68236415,
			id:      "4cJPL",
		},
		{
			results: 67175877,
			id:      "4XrVt",
		},
		{
			results: 65592921,
			id:      "4RDiL",
		},
		{
			results: 41238522,
			id:      "2n21S",
		},
		{
			results: 94283385,
			id:      "6NbPl",
		},
		{
			results: 41956665,
			id:      "2q2qP",
		},
		{
			results: 37707752,
			id:      "2YDVY",
		},
		{
			results: 68610149,
			id:      "4dsdJ",
		},
		{
			results: 28209392,
			id:      "1uMYC",
		},
		{
			results: 38691855,
			id:      "2cLWB",
		},
		{
			results: 99522769,
			id:      "6jaPx",
		},
		{
			results: 69682314,
			id:      "4iNYI",
		},
		{
			results: 16162743,
			id:      "15ofP",
		},
		{
			results: 38386419,
			id:      "2b43n",
		},
		{
			results: 42666958,
			id:      "2t1ck",
		},
		{
			results: 98849355,
			id:      "6glER",
		},
		{
			results: 17320289,
			id:      "1AfnV",
		},
		{
			results: 25808476,
			id:      "1kHxk",
		},
		{
			results: 91304382,
			id:      "6B6RK",
		},
		{
			results: 33144228,
			id:      "2F4KK",
		},
		{
			results: 61813865,
			id:      "4BMbp",
		},
		{
			results: 90135250,
			id:      "66CIM",
		},
		{
			results: 84605551,
			id:      "5izlf",
		},
		{
			results: 75446667,
			id:      "56Z7j",
		},
		{
			results: 35347002,
			id:      "2OJMw",
		},
		{
			results: 84988268,
			id:      "5kbKW",
		},
		{
			results: 344523,
			id:      "1Rcp",
		},
		{
			results: 34542600,
			id:      "2Kw6i",
		},
		{
			results: 11579636,
			id:      "maOK",
		},
		{
			results: 50145891,
			id:      "3OPEh",
		},
		{
			results: 97880514,
			id:      "6chBy",
		},
		{
			results: 44866383,
			id:      "32FnL",
		},
		{
			results: 29717072,
			id:      "20glc",
		},
		{
			results: 39625641,
			id:      "2gGRF",
		},
		{
			results: 9977484,
			id:      "frbA",
		},
		{
			results: 38415131,
			id:      "2bBWt",
		},
		{
			results: 25011862,
			id:      "1gwj8",
		},
		{
			results: 4742753,
			id:      "Jto1",
		},
		{
			results: 30922277,
			id:      "25kIP",
		},
		{
			results: 18561924,
			id:      "1Fsns",
		},
		{
			results: 86782082,
			id:      "5s7z0",
		},
		{
			results: 42781204,
			id:      "2tVLQ",
		},
		{
			results: 21860274,
			id:      "1Tir4",
		},
		{
			results: 82510071,
			id:      "5aCdb",
		},
		{
			results: 37078850,
			id:      "2VZty",
		},
		{
			results: 55022721,
			id:      "3irvF",
		},
		{
			results: 51293216,
			id:      "3TDhw",
		},
		{
			results: 90557044,
			id:      "67y1U",
		},
		{
			results: 94867994,
			id:      "6Q3Uw",
		},
		{
			results: 26714036,
			id:      "1o5XY",
		},
		{
			results: 11013256,
			id:      "kD3A",
		},
		{
			results: 37169131,
			id:      "2VxO7",
		},
		{
			results: 7988373,
			id:      "XW8j",
		},
		{
			results: 52386136,
			id:      "3Xo1g",
		},
	}

	for k, test := range tests {
		result := ShortToInt(test.id)
		if result != test.results {
			t.Errorf("Test %v: int %v from ShortToInt(%v), but got %v",
				k, test.id, test.results, result)
		}
	}
}
