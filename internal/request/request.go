package request

import (
	"errors"
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

func parseRequestLine(requestString string) (*Request, error) {
	requestLines := strings.Split(requestString, "\r\n")

	if len(requestLines) == 0 || requestLines[0] == "" {
		return nil, errors.New("empty request")
	}

	requestLineParts := strings.Fields(requestLines[0])
	if len(requestLineParts) != 3 {
		return nil, errors.New("invalid request line format")
	}

	httpVersion := requestLineParts[2]
	if strings.HasPrefix(httpVersion, "HTTP/") {
		httpVersion = strings.TrimPrefix(httpVersion, "HTTP/")
	}

	if httpVersion != "1.1" {
		return nil, errors.New("Only HTTP/1.1 is supported")
	}

	if !httpMethodRegex.MatchString(requestLineParts[0]) {
		return nil, errors.New("HTTP method must contain only uppercase letters")
	}
		
	request := &Request{
		RequestLine: RequestLine{
			Method:					requestLineParts[0],
			RequestTarget:	requestLineParts[1],
			HttpVersion:		httpVersion,
		},
	}


	return request, nil
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	request, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read from reader: %w", err)
	}

	requestString := string(request)

	return parseRequestLine(requestString)

}
