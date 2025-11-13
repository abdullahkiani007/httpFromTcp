package headers

import (
	"bytes"
	"fmt"
	"strings"
)

type Headers struct {
	Header map[string]string
}

func (h *Headers) Set(key string, val string) {
	v := h.Get(key)
	if len(v) >= 1 {
		h.Header[strings.ToLower(key)] = v + "," + val
	} else {
		h.Header[strings.ToLower(key)] = val
	}
}

func (h *Headers) Get(key string) string {
	return h.Header[strings.ToLower(key)]
}

func isValidToken(key string) bool {
	var found bool
	for _, c := range key {
		found = false
		switch c {
		case '!', '#', '$', '%', '&', '\'', '*', '+', '-', '.', '^', '_', '`', '|', '~':
			found = true
		default:
			if c >= 'A' && c <= 'Z' || c >= 'a' && c <= 'z' {
				found = true
			}

			// fmt.Printf(" chr %v code %v found %v\n", string(c), c, found)

			if !found {
				return false
			}
		}
	}

	return true

}

func parseHeader(b []byte) (string, string, error) {

	parts := strings.SplitN(string(b), ":", 2)
	fmt.Printf("parts %v", parts)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("malformedsd header line parts %v", parts)
	}
	if len(parts[0]) < 1 {
		return "", "", fmt.Errorf("malformedd header line %v", len(parts))
	}
	// fmt.Printf("\nlen of name before %v , %v.\n", len(parts[0]), parts[0])
	name := strings.TrimLeft(parts[0], " ")

	if bytes.IndexByte([]byte(name), 32) != -1 {
		return "", "", fmt.Errorf("invalid spacing in header")
	}

	if !isValidToken(name) {
		return "", "", fmt.Errorf("invalid characters")
	}

	val := strings.TrimSpace(parts[1])
	return name, val, nil
}

var rn = []byte("\r\n")

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	n = 0
	done = false

	fmt.Printf("Parsing header '%v'\n", string(data))

	for {
		idx := bytes.Index(data[n:], rn)
		fmt.Printf("idx %v bytes read %v\n", idx, n)
		fmt.Printf("Header %v\n", string(data[n:]))
		if idx == -1 {
			break
			// return n, done, nil
		}

		if idx == 0 {
			done = true
			n += len(rn)
			break
		}

		name, value, err := parseHeader(data[n : n+idx])
		if err != nil {
			return 0, false, err
		}
		n = n + idx + len(rn)

		h.Set(name, value)
		fmt.Printf("bytes read %v\n", n)
		fmt.Printf("header %v\n", h.Header)
	}
	return n, done, nil
}

func NewHeaders() *Headers {
	header := Headers{
		Header: map[string]string{},
	}
	return &header
	// return header
}
