package server

import (
	"bytes"
	"fmt"
	"io"
	"net"

	"github.com/Numeez/go-http/internal/request"
	"github.com/Numeez/go-http/internal/response"
)

type HandlerError struct {
	StatusCode response.StatusCode
	Message    string
}
type Handler func(w *response.Writer, req *request.Request) *HandlerError

type Server struct {
	port     uint16
	listener net.Listener
	handler  Handler
	close    bool
}

func Serve(port uint16, handler Handler) (*Server, error) {
	server := &Server{}
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}
	server.listener = listener
	server.port = port
	server.close = false
	server.handler = handler
	go server.listen()
	return server, nil
}
func (s *Server) listen() {
	for {
		if s.close {
			return
		}
		conn, err := s.listener.Accept()
		if err != nil {
			return
		}
		go handleConnection(s, conn)
	}

}

func (s *Server) Close() error {
	s.close = true
	if err := s.listener.Close(); err != nil {
		return err
	}
	return nil
}

func handleConnection(server *Server, conn io.ReadWriteCloser) {
	headers := response.GetDefaultHeaders(0)
	r, err := request.RequestFromReader(conn)
	responseWriter:=response.NewWriter(conn)
	if err != nil {
		_ = responseWriter.WriteStatusLine(response.StatusCode(response.HttpStatusBadRequest))
		_ = responseWriter.WriteHeaders(headers)
		return
	}
	writer := bytes.NewBuffer([]byte{})
	statusCode := response.HttpStatusOk
	handlerError := server.handler(responseWriter, r)
	var body []byte = nil
	if handlerError != nil {
		statusCode = handlerError.StatusCode
		body = []byte(handlerError.Message)
	} else {
		body = writer.Bytes()
	}
	headers.Replace("Content-length", fmt.Sprintf("%d", len(body)))
	_ = responseWriter.WriteStatusLine(statusCode)
	_ = responseWriter.WriteHeaders(headers)
	_,_ = conn.Write(body)
	_ = conn.Close()
}
