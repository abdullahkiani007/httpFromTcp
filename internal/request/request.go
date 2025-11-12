package request

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"slices"
	"strings"
)

type parserState string

const (
	StateInit parserState = "init"
	StateDone parserState = "done"
)

type Request struct {
	RequestLine RequestLine
	state       parserState
}

func (r *Request) parse(b []byte) (int, error) {
	reqLine, n, err := parseRequestLine(b)
	if r.state == StateDone {
		return 0, errors.New("finished parsing request")
	}
	if err != nil {
		return 0, err
	}
	if n == 0 {
		return n, nil
	}

	r.RequestLine = *reqLine
	r.state = StateDone
	return n, nil

}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	buff := make([]byte, 1024)
	buffLen := 0
	req := &Request{
		state: StateInit,
	}

	for {
		n, err := reader.Read(buff[buffLen:])

		if err != nil {
			if errors.Is(err, io.EOF) {
				req.state = StateDone
				break
			}
			break
		}

		buffLen += n

		readN, err := req.parse(buff[:buffLen])

		if req.state == StateDone {
			break
		}
		if err != nil {
			return req, err
		}

		//  todo Understand this one
		copy(buff, buff[readN:buffLen])
		buffLen -= readN

	}

	return req, nil
}

func parseRequestLine(b []byte) (*RequestLine, int, error) {

	requestLine := &RequestLine{}
	index := bytes.IndexByte(b, 13)
	if index == -1 {
		return requestLine, 0, nil
	}
	buff := b[:index]
	parts := strings.Split(string(buff), " ")

	if len(parts) != 3 {
		fmt.Printf("parts %v", parts)
		return nil, 0, fmt.Errorf("invalid request %v ", len(parts))
	}
	METHODS := []string{"GET", "POST", "PUT", "DELETE"}

	if !slices.Contains(METHODS, parts[0]) {
		return nil, 0, fmt.Errorf("invalid http method %v ", parts[0])
	}

	v := strings.Split(parts[2], "/")
	if len(v) != 2 {
		return nil, 0, fmt.Errorf("invalid http version %v ", parts[2])
	}
	version := v[1]

	line := RequestLine{
		Method:        string(parts[0]),
		RequestTarget: string(parts[1]),
		HttpVersion:   version,
	}

	return &line, len(b), nil
}
