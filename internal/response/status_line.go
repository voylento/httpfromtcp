package response

import (
	"fmt"
)

type StatusCode int

const (
	StatusCodeSuccess					StatusCode = 200
	StatusCodeBadRequest				StatusCode = 400
	StatusCodeInternalServerError		StatusCode = 500
)

func getStatusLine(statusCode StatusCode) []byte {
	var reasonPhrase string
	switch statusCode {
	case StatusCodeSuccess:
		reasonPhrase = "OK"
	case StatusCodeBadRequest:
		reasonPhrase = "Bad Request"
	case StatusCodeInternalServerError:
		reasonPhrase = "Internal Server Error"
	}

	var b []byte	
	b = fmt.Appendf(b, "HTTP/1.1 %d %s\r\n", statusCode, reasonPhrase)
	return b
}

