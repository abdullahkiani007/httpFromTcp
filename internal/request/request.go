package request

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"slices"
	"strconv"
	"strings"

	"github.com/abdullahkiani007/httpfromtcp/internal/headers"
)

type parserState string

// type requestStateParsingHeaders string

const (
	StateInit           parserState = "init"
	StateDone           parserState = "done"
	StateParsingHeaders parserState = "header"
	StateParseBody      parserState = "body"
)

type Request struct {
	RequestLine RequestLine
	Headers     headers.Headers
	state       parserState
	Body        []byte
}

var rn = []byte("\r\n")

func (r *Request) parse(b []byte) (int, error) {
	var n int
	var err error
	var reqLine *RequestLine
	var done bool
	var tb int
	if r.state == StateInit {
		reqLine, n, err = parseRequestLine(b)

		if r.state == StateDone {
			return 0, errors.New("finished parsing request")
		}
		if err != nil {
			return 0, err
		}
		if n == 0 {
			return n, nil
		}
		tb = n
		fmt.Printf("1.total bytes %v\n", tb)

		r.RequestLine = *reqLine
		r.state = StateParsingHeaders
		fmt.Println("changing state to stateparsing headers")
		// fmt.Printf("1.number of bytes read %v, bytes read:%v, bytes left:%v\n", n, string(b[:n]), string(b[:]))
		fmt.Println(r.state)
	}
	if r.state == StateParsingHeaders {
		fmt.Printf("parsgin headers %v \n", string(b[tb:]))
		fmt.Printf("2.total bytes %v\n", tb)
		n, done, err = r.Headers.Parse(b[tb:])
		tb += n
		fmt.Printf("3.total bytes %v\n", tb)

		fmt.Printf("23: number of bytes read %v, bytes read:%v, bytes left:%v\n", n, string(b[:tb]), string(b[tb:]))

		if err != nil {
			return 0, err
		} else if done {

			if l := r.Headers.Get("content-length"); len(l) != 0 {
				fmt.Printf("changing to state to parsing body of len %v\n", l)
				fmt.Printf("headers parsed %v\n", r.Headers.Header)
				r.state = StateParseBody
			} else {
				fmt.Println("Changing state to done")
				r.state = StateDone
			}
		} else {
			fmt.Printf("2.number of bytes read %v, bytes read:%v, bytes left:%v\n", n, string(b[:tb]), string(b[tb:]))
			fmt.Printf("2. returning number of bytes read %v", tb)
			return tb, nil
		}
	}
	if r.state == StateParseBody {
		offSet := tb
		fmt.Printf("Parsing body %v length of body %v\n", string(b[tb:]), len(b))
		fmt.Printf("body %v\n", string(r.Body))
		r.Body = append(r.Body, b[tb:]...)
		fmt.Printf("after body %v\n", string(r.Body))

		l := r.Headers.Get("content-length")
		if l != "" {
			contentLength, err := strconv.Atoi(l)
			fmt.Printf("Content lengths is %v\n", contentLength)
			if err != nil {
				return 0, fmt.Errorf("invalid content-length value: %v", err)
			}
			if contentLength > len(r.Body) {
				return len(b[tb:]) + offSet, nil
			} else {
				fmt.Println("Changing state to done in bocy")
				r.state = StateDone
				n = len(b[tb:]) + offSet
				fmt.Printf("bytes parsed %v\n", n)
			}

		}
	}
	fmt.Printf("3.number of bytes read %v, bytes read:%v, bytes left:%v\n", n, string(b[:n]), string(b[n:]))
	fmt.Printf("3. returning number of bytes read %v", n)

	return n, nil

}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	buff := make([]byte, 4096)
	buffLen := 0
	req := &Request{
		state:   StateInit,
		Headers: *headers.NewHeaders(),
		Body:    []byte{},
	}

	for {
		fmt.Printf("passing this buff to reader %v, len of buff %v\n", string((buff[buffLen:])), len((buff[buffLen:])))
		n, err := reader.Read(buff[buffLen:])

		fmt.Printf("Parser state %v\n", req.state)
		fmt.Printf("org buff %v, buffLen %v\n", string(buff), buffLen)

		if err != nil {
			if errors.Is(err, io.EOF) {
				req.state = StateDone
				l := req.Headers.Get("content-length")
				if l != "" {
					contentLength, err := strconv.Atoi(l)
					fmt.Printf("Content lengths is %v\n", contentLength)
					if err != nil {
						return req, fmt.Errorf("invalid content-length value: %v", err)
					}
					if contentLength > len(req.Body) {
						return req, fmt.Errorf("invalid body: %v", err)
					}
				}

				fmt.Println("breakoing out of loop 1")
				break
			}
			fmt.Println("breakoing out of loop 2")

			break
		}

		buffLen += n
		fmt.Println("passing this buff to parser")
		for _, v := range buff[:buffLen+10] {
			if v == 13 || v == 10 {
				fmt.Printf("-")
			} else {
				fmt.Printf("%v", string(v))
			}
		}
		fmt.Println()

		readN, err := req.parse(buff[:buffLen])
		fmt.Printf("Number of bytes read %v\n", readN)

		if req.state == StateDone {
			fmt.Println("breakoing out of loop 3")
			break
		}
		if err != nil {

			fmt.Println("breakoing out of loop 4")
			fmt.Printf("err : %e\n", err)
			return req, err
		}
		fmt.Printf("before  copy buff %v, buffLen %v\n", string(buff), buffLen)

		//  todo Understand this one
		copy(buff, buff[readN:buffLen])
		buffLen -= readN
		fmt.Printf("after copy buff %v, buffLen %v\n", string(buff), buffLen)

	}
	fmt.Printf("loop is going to end %v\n", req)
	fmt.Printf("parsing body %v\n", string(req.Body))
	return req, nil
}

func parseRequestLine(b []byte) (*RequestLine, int, error) {

	requestLine := &RequestLine{}
	index := bytes.Index(b, rn)

	// index := bytes.IndexByte(b, 13)
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
	fmt.Printf("number of bytes read is %v\n", len(b))
	return &line, len(buff) + 2, nil
}
