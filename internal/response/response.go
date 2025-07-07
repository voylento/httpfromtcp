package response

import (
	"fmt"
	"io"
	"net/http"

	"voylento/httpfromtcp/internal/headers"
)

const crlf = "\r\n"

type WriteState int
const (
	WriteStateStatusLine WriteState = iota
	WriteStateHeaders
	WriteStateBody
	WriteStateDone
)

type Writer struct {
	Writer 	io.Writer
	State	WriteState
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{
		State: WriteStateStatusLine,
		Writer: w,
	}
}

func writeStateToString(state WriteState) string {
	switch state{
	case WriteStateStatusLine:
		return "writeStateStatusLine"
	case WriteStateHeaders:
		return "writeStateHeaders"
	case WriteStateBody:
		return "writeStateBody"
	case WriteStateDone:
		return "writeStateDone"
	default:
		return "Unknown state"
	}
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	if w.State != WriteStateStatusLine {
		return fmt.Errorf("Error: attempting to write status line when state is %x", writeStateToString(w.State))
	}
	defer func() {w.State = WriteStateHeaders}()
	_, err := w.Writer.Write(getStatusLine(statusCode))
	return err
}

func (w *Writer) WriteHeaders(h headers.Headers) error {
	if w.State != WriteStateHeaders {
		return fmt.Errorf("Error: attempting to write headers when state is %s", writeStateToString(w.State)) 
	}
	defer func() {w.State = WriteStateBody}()
	for k, v := range h {
		canonicalName := http.CanonicalHeaderKey(k)
		_, err := fmt.Fprintf(w.Writer, "%s: %s%s", canonicalName, v, crlf)
		if err != nil {
			return err
		}
	}
	_, err := fmt.Fprintf(w.Writer, crlf)
	return err
}

func (w *Writer) WriteBody(p []byte) (int, error) {
	if w.State != WriteStateBody {
		return 0, fmt.Errorf("Error: attempting to write body when state is %s", writeStateToString(w.State))
	}
	defer func() {w.State = WriteStateDone}()
	return w.Writer.Write(p)
}

func (w *Writer) WriteChunkedBody(p []byte) (int, error) {
	if w.State != WriteStateBody {
		return 0, fmt.Errorf("Error: attempting to write body when state is %s", writeStateToString(w.State))
	}

	chunkSize := len(p)
	nTotal := 0

	n, err := fmt.Fprintf(w.Writer, "%X\r\n", chunkSize)
	if err != nil {
		return nTotal, err
	}
	nTotal += n

	n, err = w.Writer.Write(p)
	if err != nil {
		return nTotal, err
	}
	nTotal += n

	n, err = fmt.Fprintf(w.Writer, "\r\n")
	if err != nil {
		return nTotal, err
	}
	nTotal += n
	return nTotal, nil
}

func (w *Writer) WriteChunkedBodyDone() (int, error) {
	return w.Writer.Write([]byte("0\r\n\r\n"))
}
