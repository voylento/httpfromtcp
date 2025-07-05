package request

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"

	"voylento/httpfromtcp/internal/headers"
)

type Request struct {
	RequestLine 	RequestLine
	Headers 		headers.Headers
	Body			[]byte
	state			requestState
}

type RequestLine struct {
	HttpVersion		string
	RequestTarget	string
	Method			string
}

type requestState int

const (
	requestStateInitialized requestState = iota
	requestStateParsingHeaders
	requestStateParsingBody
	requestStateDone
)

var httpMethodRegex = regexp.MustCompile(`^[A-Z]+$`)
const crlf = "\r\n"
const bufferSize = 8

func RequestFromReader(reader io.Reader) (*Request, error) {
	buf := make([]byte, bufferSize)
	readToIndex := 0
	req := &Request{
		Headers: headers.NewHeaders(),
		state: requestStateInitialized,
	}

	for req.state != requestStateDone {
		if readToIndex >= len(buf) {
			newBuf := make([]byte, len(buf)*2)
			copy(newBuf, buf)
			buf = newBuf
		}
	
		numBytesRead, err := reader.Read(buf[readToIndex:])
		if err != nil {
			if errors.Is(err, io.EOF) {
				if req.state != requestStateDone {
					return nil, fmt.Errorf("Error: EOF before request fully processed")
				}
				break
			}
			return nil, err
		}

		readToIndex += numBytesRead // add count of bytes to total read from io

		numBytesParsed, err := req.parse(buf[:readToIndex])
		if err != nil {
			return nil, err
		}

		copy(buf, buf[numBytesParsed:])
		readToIndex -= numBytesParsed
	}

	return req, nil
}

func (r *Request) parse(data []byte) (int, error) {
	totalBytesParsed := 0

	for r.state != requestStateDone {
		n, err := r.parseSingle(data[totalBytesParsed:])
		if err != nil {
			return 0, err
		}
		if n == 0 {
			return totalBytesParsed, nil
		}
		totalBytesParsed += n
	}

	return totalBytesParsed, nil
}

func (r *Request) parseSingle(data []byte) (int, error) {
	switch r.state {
	case requestStateInitialized:
		requestLine, n, err := parseRequestLine(data)
		if err != nil {
			return 0, err
		}
		if n == 0 {
			return 0, nil
		}
		r.RequestLine = *requestLine
		r.state = requestStateParsingHeaders
		return n, nil
	case requestStateParsingHeaders:
		n, done, err := r.Headers.Parse(data)
		if err != nil {
			return 0, err
		}
		if n == 0 {
			return 0, nil
		}
		if done {
			r.state = requestStateParsingBody
		}
		return n, nil 
	case requestStateParsingBody:
		bodyLength, exists := r.Headers.Get("content-length")
		contentLength, err := strconv.Atoi(bodyLength)
		if !exists || contentLength < 1 {
			r.state = requestStateDone
			return 0, nil
		}
		body, n, err := parseBody(data, contentLength)
		if err != nil {
			return 0, err
		}
		if n == 0 {
			return 0, nil
		}
		if n != contentLength {
			return 0, fmt.Errorf("Error: body length does not match content-length: %d", contentLength)
		}
		r.Body = body
		r.state = requestStateDone
		return n, nil	
	default:
		return 0, fmt.Errorf("Error: unknown parse state: %d", r.state)
	}
}

func parseBody(data []byte, length int) ([]byte, int, error) {
	if len(data) > length {
		return nil, 0, fmt.Errorf("Error: body length (%d) greater than content-length (%d)", len(data), length)
	}
	if len(data) < length {
		// need more data
		return nil, 0, nil
	}
	buf := make([]byte, length)
	copy(buf, data[:length])
	return buf, length, nil
}

func parseRequestLine(data []byte) (*RequestLine, int, error) {
	idx := bytes.Index(data, []byte(crlf))
	if idx == -1 {
		return nil, 0, nil
	}

	requestLineText := string(data[:idx])
	requestLine, err := parseRequestLineFromString(requestLineText)
	if err != nil {
		return nil, 0, err
	}

	return requestLine, idx + 2, nil
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

	httpVersion, found := strings.CutPrefix(fields[2], "HTTP/") 
	if !found {
		return nil, fmt.Errorf("Http version invalid format")
	}
	if httpVersion != "1.1" {
		return nil, fmt.Errorf("unrecognized HTTP-version: %s", httpVersion)
	}

	return &RequestLine{
			Method:					method,
			RequestTarget:	requestTarget,	
			HttpVersion:		httpVersion,
		}, nil
}
