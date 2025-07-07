package main

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"net/http"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"voylento/httpfromtcp/internal/headers"
	"voylento/httpfromtcp/internal/request"
	"voylento/httpfromtcp/internal/response"
	"voylento/httpfromtcp/internal/server"
)

const port = 42069

func main() {
	server, err := server.Serve(port, handler)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}

func handler(w *response.Writer, req *request.Request) {
	if strings.HasPrefix(req.RequestLine.RequestTarget, "/httpbin") {
		proxyHandler(w, req)
		return
	}
	if strings.HasPrefix(req.RequestLine.RequestTarget, "/video") {
		videoHandler(w, req)
		return
	}
	if req.RequestLine.RequestTarget == "/yourproblem" {
		handler400(w, req)
		return
	} 
	if req.RequestLine.RequestTarget == "/myproblem" {
		handler500(w, req)
		return
	}
	handler200(w, req)
}

func proxyHandler(w *response.Writer, req *request.Request) {
	target := strings.TrimPrefix(req.RequestLine.RequestTarget, "/httpbin/")
	url := fmt.Sprintf("https://httpbin.org/%s", target)
	fmt.Printf("Proxying to %s\n", url)
	resp, err := http.Get(url)
	if err != nil {
		handler500(w, req)
		return
	}
	defer resp.Body.Close()

	h := response.GetDefaultHeaders(0)
	h.Override("Transfer-Encoding", "chunked")
	h.Remove("Content-Length")
	h.Set("Trailer", "X-Content-SHA256, X-Content-Length")
	w.WriteStatusLine(response.StatusCodeSuccess)
	w.WriteHeaders(h)

	const maxChunkSize = 1024
	
	var fullResponse []byte
	buf := make([]byte, maxChunkSize)	
	
	for {
		n, err := resp.Body.Read(buf)
		fmt.Printf("Read %d bytes\n", n)
		if n > 0 {
			fullResponse = append(fullResponse, buf[:n]...)
			_, err = w.WriteChunkedBody(buf[:n])
			if err != nil {
				log.Printf("Error writing chunked body: %v\n", err)
				break
			}
		}
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			log.Printf("Error reading response body for httpbin.org: %v\n", err)
			break
		}
	}

	_, err = w.WriteChunkedBodyDone()
	if err != nil {
		log.Printf("Error writing chunked body done: %v\n", err)
	}

	hash := sha256.New()
	hash.Write(fullResponse)
	sum := hash.Sum(nil)
	fmt.Printf("sha256 = %x\n", sum)
	fmt.Printf("len fullResponse = %d\n", len(fullResponse))

	trailers := headers.NewHeaders()
	trailers.Set("X-Content-SHA256", fmt.Sprintf("%x", hash.Sum(nil)))
	trailers.Set("X-Content-Length", fmt.Sprintf("%d", len(fullResponse)))
	w.WriteTrailers(trailers)
	err = w.FinalizeChunkedResponse()
}

func videoHandler(w *response.Writer, req *request.Request) {
	fileName := "./assets/vim.mp4"
	fileBytes, err := os.ReadFile(fileName)
	if err != nil {
		fmt.Printf("Error reading video file: %v\n", err)
		return
	}


	h:= response.GetDefaultHeaders(len(fileBytes))
	h.Override("Content-Type", "video/mp4")	
	w.WriteStatusLine(response.StatusCodeSuccess)
	w.WriteHeaders(h)
	w.WriteBody(fileBytes)
}

func handler400(w *response.Writer, req *request.Request) {
	writeResponse(w, 
			req, 
			response.StatusCodeBadRequest, 
			[]byte(`
<html>
  <head>
    <title>400 Bad Request</title>
  </head>
  <body>
    <h1>Bad Request</h1>
    <p>Your request honestly kinda sucked.</p>
  </body>
</html>`))
}

func handler500(w *response.Writer, req *request.Request) {
	writeResponse(w, 
			req, 
			response.StatusCodeBadRequest, 
			[]byte(`
<html>
  <head>
    <title>500 Internal Server Error</title>
  </head>
  <body>
    <h1>Internal Server Error</h1>
    <p>Okay, you know what? This one is on me.</p>
  </body>
</html>`))
}

func handler200(w *response.Writer, req *request.Request) {
	writeResponse(w,
			req,
			response.StatusCodeSuccess,
			[]byte(`
<html>
  <head>
    <title>200 OK</title>
  </head>
  <body>
    <h1>Success!</h1>
    <p>Your request was an absolute banger.</p>
  </body>
</html>
`))
}


func writeResponse(w *response.Writer, _ *request.Request, code response.StatusCode, body []byte) {
	w.WriteStatusLine(code)
	h := response.GetDefaultHeaders(len(body))
	h.Override("Content-Type", "text/html")
	w.WriteHeaders(h)
	w.WriteBody(body)
}

