package server

import (
	"fmt"
	"log"
	"net"
	"sync/atomic"

	"voylento/httpfromtcp/internal/response"
	"voylento/httpfromtcp/internal/request"
)

type Handler func(w *response.Writer, req *request.Request)

type Server struct {
	listener net.Listener
	handler	Handler
	closed	atomic.Bool
}

func Serve(port int, handler Handler) (*Server, error) {
	// Listen on the tcp port provided
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}

	s := &Server{
		listener: l,
		handler:  handler,
	}
	go s.listen()
	return s, nil	
}

func (s *Server) listen() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if s.closed.Load() {
				return
			}
			log.Printf("Error accepting connection: %v", err)
			continue
		}
		go s.handle(conn)
	}
}

func (s *Server) Close() error {
	s.closed.Store(true)
	if s.listener != nil {
		return s.listener.Close()
	}
	return nil
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()
	w := response.NewWriter(conn)
	req, err := request.RequestFromReader(conn)
	if err != nil {
		w.WriteStatusLine(response.StatusCodeBadRequest)
		var b []byte
		body := fmt.Appendf(b, "Error parsing request: %v", err)
		w.WriteHeaders(response.GetDefaultHeaders(len(body)))
		w.WriteBody(body)
		return
	}
	s.handler(w, req)
}

