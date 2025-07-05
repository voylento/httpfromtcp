package headers

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
)

type Headers map[string]string

const crlf = "\r\n"
var headerNameRegex = regexp.MustCompile("^[A-Za-z0-9!#$%&'*+\\-.^_`|~]+$")

const colon = ":"

func ValidateHeaderName(name string) bool {
	return headerNameRegex.MatchString(name)
}

func NewHeaders() Headers {
	return make(Headers)
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	idx := bytes.Index(data, []byte(crlf))	
	if idx == -1 {
		return 0, false, nil
	}
	if idx == 0 {
		return 2, true, nil
	}

	before, after, found:= bytes.Cut(data[:idx], []byte(":"))
	if !found {
		return 0, false, fmt.Errorf("Error: invalid header format: %s", string(data[:idx]))
	}

	key := string(before)
	if key != strings.TrimRight(key, " ") {
		return 0, false, fmt.Errorf("Error: invalid header format: %s", key)
	}

	value := bytes.TrimSpace(after)
	key = strings.TrimSpace(key)

	if !ValidateHeaderName(key) {
		return 0, false, fmt.Errorf("Error: Invalid header name: %s", key)
	}

	h.Set(strings.ToLower(key), string(value))
	return idx+2, false, nil
}

func (h Headers) Set(key, value string) {
	existingValue, exists := h[key]
	if exists {
		h[key] = existingValue + ", " + value
	} else {
		h[key] = value
	}
}

func (h Headers) Get(key string) (string, bool) {
	value, exists := h[strings.ToLower(key)]
	return value, exists
}
