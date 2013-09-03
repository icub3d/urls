package urls

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

var (
	// re is the regular expression used to check for valid ids
	re = regexp.MustCompile("^[0-9a-zA-Z]+$")
)

const (
	// The base we are working in ([0-9a-zA-Z])
	base = 62
)

// ValidID returns true if the given string is a valid ID.
func ValidID(id string) bool {
	return re.MatchString(id)
}

// IntToShort returns the string representation of the given
// integer. Values less than 0 return 0. Otherwise, it will be some
// string that includes the characters 0-9, a-z, and A-Z.
func IntToShort(i int) string {
	if i <= 0 {
		return "0"
	}

	s := ""
	for i > 0 {
		r := i % base
		s = digit(r) + s
		i = i / base
	}

	return s
}

// Digit convers the given integer into its representative single
// digit in the language.
func digit(i int) string {
	if i < 10 {
		i = i + 48 // 48 is where 0-9 starts
	} else if i >= 10 && i <= 35 {
		i = i + 55 // 55+10 is where A-Z starts.
	} else if i >= 36 && i < 62 {
		i = i + 61 // 61+36 is where a-z starts.
	} else {
		// Anything else doesn't fit into [0-9a-zA-Z].
		return ""
	}

	return string([]byte{byte(i)})
}

// paramGetInt is a helper function that returns the integer value of the
// query paramGetInt with the given key.
func paramGetInt(q url.Values, key string) int {
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

// getLimitOffset is a helper function that gets the limit and offset
// values from the query parameters and sets them to sane values if
// they are not sane. Limit defaults to 20 and offset 0. If limit >
// 100, limit will be set to 100. If values are negative, they are set
// to their default.
func getLimitOffset(q url.Values) (int, int) {
	// Get the query parameters.
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

	return limit, offset
}

// marshalAndWrite is a helper function that marshals the given data
// and writes it to the ResponseWrite. If marshalling fails, "oops" is
// returns as well as the http.StatusInternalServerError.
func marshalAndWrite(w http.ResponseWriter, i interface{}) {
	// Marshal it to JSON.
	enc, err := json.Marshal(i)
	if err != nil {
		log.Printf("Marshal(%v) failed with: %v", i, err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("oops"))
		return
	}

	// Write the response.
	w.Write(enc)

}

// parseUserAgent looks for keywords in the given string and returns
// the best guess for the browser and platform.
func parseUserAgent(ua string) (string, string) {
	browser := "Unknown"
	platform := "Unknown"

	// These are in this order because some contain the others.
	browsers := []string{
		"Chrome",
		"Safari",
		"Firefox",
		"MSIE",
		"Opera",
	}
	for _, search := range browsers {
		if strings.Contains(ua, search) {
			browser = search
			break
		}
	}

	platforms := map[string]string{
		"Linux":          "Linux",
		"Windows NT 5.1": "Windows XP",
		"Windows NT 6.1": "Windows 7",
		"Windows NT 6.2": "Windows 8",
		"Mac OS X":       "Mac OS X",
		"iPhone":         "iOS",
		"iPad":           "iOS",
	}

	for search, value := range platforms {
		if strings.Contains(ua, search) {
			platform = value
			break
		}
	}

	return browser, platform
}

// Determine country attempts to determine the country of origin by
// the IP Address.
func determineCountry(addr string) string {
	country := "Unknown"

	// Get the parts.
	addr = strings.Split(addr, ":")[0]
	parts := strings.Split(addr, ".")
	if len(parts) != 4 {
		return country
	}

	// get the integer equivalent.
	total := 0
	i, err := strconv.Atoi(parts[0])
	if err != nil {
		return country
	}
	total = total + 16777216*i

	i, err = strconv.Atoi(parts[1])
	if err != nil {
		return country
	}
	total = total + 65536*i

	i, err = strconv.Atoi(parts[2])
	if err != nil {
		return country
	}
	total = total + 256*i

	i, err = strconv.Atoi(parts[3])
	if err != nil {
		return country
	}
	total = total + i

	for _, ipr := range ipLookup {
		if total > ipr.start && total < ipr.end {
			return ipr.country
		}
	}

	return country
}
