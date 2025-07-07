package response

import (
	"fmt"
	"voylento/httpfromtcp/internal/headers"
)

func GetDefaultHeaders(contentLen int) headers.Headers {
	headers := headers.NewHeaders()
	headers.Set("content-type", "text/plain")
	headers.Set("connection", "close")
	headers.Set("content-length", fmt.Sprintf("%d", contentLen))

	return headers
}

func GetDefaultHeadersForChunkEncoding() headers.Headers {
	h:= headers.NewHeaders()
	h.Set("transfer-encoding", "chunked")
	h.Set("connection", "close")
	h.Set("content-type", "application/json")
	return h
}
