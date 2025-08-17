package request

import (
	"bytes"
	"fmt"
	"io"
)

type parserState string

const (
	StateInit  parserState = "init"
	StateDone  parserState = "done"
	StateError parserState = "error"
)

type Request struct {
	RequestLine RequestLine
	State       parserState
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

var ERR_MALFORMED_REQUEST_LINE = fmt.Errorf("malformed request line")
var ERR_MALFORMED_HTTP_VERSION = fmt.Errorf("malformed http version")
var ERR_UNSUPPORTED_HTTP_VERSION = fmt.Errorf("unsupported http version")
var ErrReqestErrorState = fmt.Errorf("request is in error state")
var SEPARATOR = []byte("\r\n")

func newRequest() *Request {
	return &Request{
		State: StateInit,
	}
}

func parseRequestLine(line []byte) (*RequestLine, int, error) {
	idx := bytes.Index(line, SEPARATOR)
	if idx == -1 {
		return nil, 0, nil
	}
	requestLine := line[:idx]
	read := idx + len(SEPARATOR)
	parts := bytes.Split(requestLine, []byte(" "))
	if len(parts) != 3 {
		return nil, read, ERR_MALFORMED_REQUEST_LINE
	}
	method := parts[0]
	target := parts[1]
	rawHttpVersion := parts[2]
	if string(rawHttpVersion) != "HTTP/1.1" {
		return nil, read, ERR_MALFORMED_HTTP_VERSION
	}
	httpVersion := bytes.Split(rawHttpVersion, []byte("/"))
	if len(httpVersion) != 2 {
		return nil, read, ERR_MALFORMED_HTTP_VERSION
	}
	version := httpVersion[1]
	rl := &RequestLine{
		Method:        string(method),
		RequestTarget: string(target),
		HttpVersion:   string(version),
	}
	return rl, read, nil
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	request := newRequest()
	buf := make([]byte, 1024)
	bufLen := 0
	for !request.done() {
		n, err := reader.Read(buf[bufLen:])
		if err != nil {
			return nil, err
		}
		bufLen += n
		readN, err := request.parse(buf[:bufLen])
		if err != nil {
			return nil, err
		}
		copy(buf, buf[readN:bufLen])
		bufLen -= readN
	}
	return request, nil
}

func (r *Request) parse(data []byte) (int, error) {
	read := 0
outer:
	for {
		switch r.State {
		case StateError:
			return 0, ErrReqestErrorState
		case StateInit:
			rl, n, err := parseRequestLine(data)
			if err != nil {
				r.State = StateError
				return 0, err
			}
			if n == 0 {
				break outer
			}
			r.RequestLine = *rl
			read += n
			r.State = StateDone

		case StateDone:
			break outer

		}
	}
	return read, nil

}

func (r *Request) done() bool {
	return r.State == StateDone || r.State == StateError
}
