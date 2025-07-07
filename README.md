# HTTPFROMTCP

This is the code I developed for the [Boot.dev](https://www.boot.dev) project [Learn the HTTP Protocol in Go](https://www.boot.dev/courses/learn-http-protocol-golang). The project aimed to demystify the inner workings of how http libraries in languages like Go work. This is a project that builds a simplistic http library on top of tcp. The course was developed by [ThePrimeagen](https://github.com/ThePrimeagen) for [Boot.dev](https://www.boot.dev)

## Get the Code

In the usual way: `git clone git@github.com:voylento/httpfromtcp.git`

## Run the code

Really? You want to run this code? Just go sign up at [Boot.dev](https://www.boot.dev) and do the course yourself. This is for learning. Yes, I put my leanring projects on github (sometimes).

Requirements: 1) obviously, you have Go installed; 2) you have curl installed

From the project root:

```
go run ./cmd/httpserver
```

That will run a localhost http server at port 42069

From another terminal window on the same machine:

```
curl -v localhost:42069/yourproblem
```
Rhis request hits a simple handler built on top of our http on tcp implementation. Should send back an html response saying "400 Bad Request" with a special message.

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
This request will not work. If you create a directory off the root names `assets` and put an mp4 video in that directory names video.mp4 and then connect to localhost:42069/video, the video should play in the browser. Since this is a toy app I didn't want to upload a video to the github repository.
