package headers

import (
	"bytes"
	"fmt"
	"strings"
)

type Headers struct {
	headers map[string]string
}

func (h *Headers) Set(key string, val string) {
	h.headers[strings.ToLower(key)] = val
}
func (h *Headers) Get(key string) string {
	return h.headers[strings.ToLower(key)]
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

			fmt.Printf(" chr %v code %v found %v\n", string(c), c, found)

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
		return "", "", fmt.Errorf("malformed header line %v", len(parts))
	}
	if len(parts[0]) < 1 {
		return "", "", fmt.Errorf("malformed header line %v", len(parts))
	}
	fmt.Printf("\nlen of name before %v , %v.\n", len(parts[0]), parts[0])
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

	fmt.Printf("Parsing header %v\n", string(data))
	idx := bytes.Index(data[n:], rn)
	fmt.Printf("idx %v\n", idx)

	if idx == -1 {
		return n, done, nil
	}

	if idx == 0 {
		done = true
		return 2, done, nil
	}

	name, value, err := parseHeader(data[:idx])
	if err != nil {
		return 0, false, err
	}
	n = idx + len(rn)
	h.Set(name, value)

	return n, done, nil
}

func NewHeaders() *Headers {
	header := Headers{
		headers: map[string]string{},
	}
	return &header
	// return header
}
