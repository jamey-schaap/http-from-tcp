package response

import (
	"fmt"
	"http-from-tcp/internal/headers"
	"io"
	"strings"
)

type writeState int

const (
	writeStateStatusLine writeState = iota
	writeStateHeaders
	writeStateBody
)

type Writer struct {
	io.Writer
	writeState
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{
		w,
		writeStateStatusLine,
	}
}

func getStatusLine(statusCode StatusCode) []byte {
	reasonPhrase := ""
	switch statusCode {
	case StatusCodeSuccess:
		reasonPhrase = "OK"
	case StatusCodeBadRequest:
		reasonPhrase = "Bad Request"
	case StatusCodeInternalServerError:
		reasonPhrase = "Internal Server Error"
	}
	return []byte(fmt.Sprintf("HTTP/1.1 %d %s\r\n", statusCode, reasonPhrase))
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	if w.writeState != writeStateStatusLine {
		return fmt.Errorf("invalid write state, expected %d actual %d", writeStateStatusLine, w.writeState)
	}

	statusLine := getStatusLine(statusCode)
	_, err := w.Write(statusLine)
	w.writeState = writeStateHeaders
	return err
}

func (w *Writer) WriteHeaders(headers headers.Headers) error {
	if w.writeState != writeStateHeaders {
		return fmt.Errorf("invalid write state, expected %d actual %d", writeStateHeaders, w.writeState)
	}

	var sb strings.Builder
	for k, v := range headers {
		sb.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}
	sb.WriteString("\r\n")
	_, err := w.Write([]byte(sb.String()))
	w.writeState = writeStateBody
	return err
}

func (w *Writer) WriteBody(p []byte) (int, error) {
	if w.writeState != writeStateBody {
		return 0, fmt.Errorf("invalid write state, expected %d actual %d", writeStateBody, w.writeState)
	}

	n, err := w.Write(p)
	// TODO: How to keep the content-length header up-to-date

	return n, err
}
