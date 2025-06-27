package request

import (
	"bytes"
	"fmt"
	"io"
	"regexp"
	"strings"
)
type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion		string
	RequestTarget	string
	Method				string
}

var httpMethodRegex = regexp.MustCompile(`^[A-Z]+$`)
const crlf = "\r\n"

func RequestFromReader(reader io.Reader) (*Request, error) {
	requestBytes, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	requestLine, err := parseRequestLine(requestBytes)
	if err != nil {
		return nil, err
	}
	return &Request{
		RequestLine: *requestLine,
	}, nil
}

func parseRequestLine(data []byte) (*RequestLine, error) {
	idx := bytes.Index(data, []byte(crlf))
	if idx == -1 {
		return nil, fmt.Errorf("CRLF not found in request line")
	}

	requestLineText := string(data[:idx])
	requestLine, err := parseRequestLineFromString(requestLineText)
	if err != nil {
		return nil, err
	}

	return requestLine, nil
}

func parseRequestLineFromString(str string) (*RequestLine, error) {
	fields := strings.Fields(str)
	if len(fields) != 3 {
		return nil, fmt.Errorf("request line format error: %s", str)
	}

	method := fields[0]

	if !httpMethodRegex.MatchString(method) {
		return nil, fmt.Errorf("HTTP method must contain only uppercase letters")
	}

	requestTarget := fields[1]

	httpVersion := fields[2]
	if strings.HasPrefix(httpVersion, "HTTP/") {
		httpVersion = strings.TrimPrefix(httpVersion, "HTTP/")
		if httpVersion != "1.1" {
			return nil, fmt.Errorf("unrecognized HTTP-version: %s", httpVersion)
		}
	} else {
		return nil, fmt.Errorf("Http version invalid format")
	}

	return &RequestLine{
			Method:					method,
			RequestTarget:	requestTarget,	
			HttpVersion:		httpVersion,
		}, nil
}
