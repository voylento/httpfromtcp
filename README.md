# HTTPFROMTCP

This is the code I developed while going through the [Boot.dev](https://www.boot.dev) course [Learn the HTTP Protocol in Go](https://www.boot.dev/courses/learn-http-protocol-golang). The project aimed to demystify the inner workings of how http libraries work. In this course, the students developed a simplistic HTTP/1.1 library on top of Go's TCP libraries. The course was developed by [ThePrimeagen](https://github.com/ThePrimeagen) for [Boot.dev](https://www.boot.dev).

## Get the Code

In the usual way: `git clone git@github.com:voylento/httpfromtcp.git`

## Run the code

Really? You want to run this code? Just go sign up at [Boot.dev](https://www.boot.dev) and do the course yourself. This is for learning. You learn by thinking through the problem sets and doing the typity typity on the keyboardy (in neovim, for best results).

### Dependencies 
- Go installed (go1.22+)
- curl
- (optional) netcat. Useful for inspecting the proxied chunked responses for the part of the course that covers chunking.

From the project root:

```
go run ./cmd/httpserver
```

That will run a localhost http server at port 42069

From another terminal window on the same machine:

```
curl -v localhost:42069/yourproblem
```
This request hits a simple handler built on top of our http on tcp implementation. Should send back an html response saying "400 Bad Request" with a special message.

```
curl -v localhost:42069/myproblem
```
This request also hits a simple handler build on the of the http on tcp implementation. Should send back an html response saying "500 Internal Server Error" with a special message.

```
curl -v localhost:42069/success
```
This request hits a simple success handler that uses the http on tcp implemenation.

```
curl -v --raw localhost:42069/httpbin/stream/50
```
This request hits a handler that proxies the request to httpbin.org and sends back the chunked response from httpbin.org. Getting this to work involved learning how to send chunked responses including proper headers and trailing headers.

```
curl -v localhost:42069/video
```
This request will not work. If you create a directory off the root project named `assets` and put an mp4 video in that directory named video.mp4 and then direct your browser to localhost:42069/video, the video should play in the browser. Since this is a toy app I didn't want to upload a video to the github repository.
