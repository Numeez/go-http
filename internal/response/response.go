package response

import (
	"fmt"
	"io"
	"strconv"

	"github.com/Numeez/go-http/internal/headers"
)

type StatusCode int

const (
	HttpStatusOk                  StatusCode = 200
	HttpStatusBadRequest          StatusCode = 400
	HttpStatusInternalServerError StatusCode = 500
)

func GetDefaultHeaders(contentLength int) headers.Headers {
	header := headers.NewHeaders()
	header.Set("Content-Length", strconv.Itoa(contentLength))
	header.Set("Connection", "close")
	header.Set("Content-Type", "text/plain")
	return *header

}

type Writer struct {
	writer io.Writer
}

func NewWriter(writer io.Writer) *Writer {
	return &Writer{
		writer: writer,
	}
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	statusLine := []byte{}
	switch statusCode {
	case HttpStatusOk:
		statusLine = []byte("HTTP/1.1 200 OK\r\n")
	case HttpStatusBadRequest:
		statusLine = []byte("HTTP/1.1 400 Bad Request\r\n")
	case HttpStatusInternalServerError:
		statusLine = []byte("HTTP/1.1 500 Internal Server Error\r\n")
	default:
		return fmt.Errorf("unrecognized status code")
	}
	_, err := w.writer.Write(statusLine)
	if err != nil {
		return err
	}
	return nil

}
func (w *Writer) WriteHeaders(h headers.Headers) error {
	var resultHeaders []byte
	h.ForEach(func(k, v string) {
		resultHeaders = append(resultHeaders, []byte(fmt.Sprintf("%s: %s\r\n", k, v))...)

	})
	resultHeaders = append(resultHeaders, []byte("\r\n")...)
	_, err := w.writer.Write(resultHeaders)
	if err != nil {
		return err
	}
	return nil

}
func (w *Writer) WriteBody(p []byte) (int, error) {
	n, err := w.writer.Write(p)
	if err != nil {
		return n, err
	}
	return n, nil
}
