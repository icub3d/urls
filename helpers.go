package urls

import (
	"net/url"
	"regexp"
	"strconv"
)

var (
	// re is the regular expression used to check for valid ids
	re = regexp.MustCompile("^[:alnum:]+$")
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
