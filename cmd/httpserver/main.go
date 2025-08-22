package main

import (
	"crypto/sha256"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/Numeez/go-http/internal/headers"
	"github.com/Numeez/go-http/internal/request"
	"github.com/Numeez/go-http/internal/response"
	"github.com/Numeez/go-http/internal/server"
)

const port = 42069

func respond400() []byte {
	return []byte(`<html>
  <head>
    <title>400 Bad Request</title>
  </head>
  <body>
    <h1>Bad Request</h1>
    <p>Your request honestly kinda sucked.</p>
  </body>
</html>`)
}

func respond500() []byte {
	return []byte(`<html>
  <head>
    <title>500 Internal Server Error</title>
  </head>
  <body>
    <h1>Internal Server Error</h1>
    <p>Okay, you know what? This one is on me.</p>
  </body>
</html>`)
}

func respond200() []byte {
	return []byte(`<html>
  <head>
    <title>200 OK</title>
  </head>
  <body>
    <h1>Success!</h1>
    <p>Your request was an absolute banger.</p>
  </body>
</html>`)
}
func toStr(data []byte) string {
	out := ""
	for _, d := range data {
		out += fmt.Sprintf("%02x", d)
	}
	return out
}
func main() {
	server, err := server.Serve(port, func(w *response.Writer, req *request.Request) *server.HandlerError {
		h := response.GetDefaultHeaders(0)
		body := respond200()
		status := response.HttpStatusOk
		if req.RequestLine.RequestTarget == "/yourproblem" {
			body = respond400()
			status = response.HttpStatusBadRequest
		} else if req.RequestLine.RequestTarget == "/myproblem" {
			body = respond500()
			status = response.HttpStatusBadRequest
		} else if req.RequestLine.RequestTarget == "/video" {
			f, err := os.ReadFile("assets/vim.mp4")
			if err != nil {
				body = respond500()
				status = response.HttpStatusBadRequest
			}
			h.Replace("content-type", "video.mp4")
			h.Replace("content-length", fmt.Sprintf("%d", len(f)))
			_ = w.WriteStatusLine(response.HttpStatusOk)
			_ = w.WriteHeaders(h)
			_, _ = w.WriteBody(f)
			return nil

		} else if strings.HasPrefix(req.RequestLine.RequestTarget, "/httpbin/") {
			traget := req.RequestLine.RequestTarget
			resp, err := http.Get("https://httpbin.org" + traget[len("/httpbin/"):])
			if err != nil {
				body = respond500()
				status = response.HttpStatusBadRequest
			} else {
				_ = w.WriteStatusLine(response.HttpStatusOk)
				h.Delete("Content-Length")
				h.Set("transfer-encoding", "chunked")
				h.Replace("Content-Type", "text/plain")
				h.Set("Trailer", "X-Content-SHA256")
				h.Set("Trailer", "X-Content-Length")
				_ = w.WriteHeaders(h)
				fullBody := []byte{}
				for {
					data := make([]byte, 32)
					n, err := resp.Body.Read(data)
					if err != nil {
						break
					}
					fullBody = append(fullBody, data[:n]...)
					_, _ = w.WriteBody([]byte(fmt.Sprintf("%x\r\n", n)))
					_, _ = w.WriteBody(data[:n])
					_, _ = w.WriteBody([]byte("\r\n"))
				}
				_, _ = w.WriteBody([]byte("0\r\n"))
				tailer := headers.NewHeaders()
				out := sha256.Sum256(fullBody)
				tailer.Set("X-Content-SHA256", toStr(out[:]))
				tailer.Set("X-Content-Length", fmt.Sprintf("%d", len(fullBody)))
				_ = w.WriteHeaders(*tailer)
			}
			return nil

		}
		_ = w.WriteStatusLine(status)
		h.Replace("Content-length", fmt.Sprintf("%d", len(body)))
		h.Replace("Content-Type", "text/html; charset=utf-8")
		_ = w.WriteHeaders(h)
		_, _ = w.WriteBody(body)

		return nil
	})
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
