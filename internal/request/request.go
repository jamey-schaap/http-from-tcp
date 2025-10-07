package request

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"
)

type Request struct {
	RequestLine RequestLine
	parserState ParserState
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type ParserState int

const (
	initialized ParserState = iota
	done
)

const crlf = "\r\n"

func RequestFromReader(reader io.Reader) (*Request, error) {
	request := &Request{}

	buffer := make([]byte, 0)
	totalReadBytes := 0
	totalParsedBytes := 0
	for request.parserState != done {
		readBuffer := make([]byte, 8)

		readBytes, err := reader.Read(readBuffer)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			return nil, err
		}
		buffer = append(buffer, readBuffer[:readBytes]...)
		totalReadBytes += readBytes

		parsedBytes, err := request.parse(buffer[totalParsedBytes:totalReadBytes])
		if err != nil {
			return nil, err
		}
		totalParsedBytes += parsedBytes
	}
	return request, nil
}

func (r *Request) parse(data []byte) (int, error) {
	if r.parserState != initialized {
		return 0, nil
	}

	requestLine, n, err := parseRequestLine(data)
	if err != nil {
		return 0, err
	}

	if n == 0 {
		return n, nil
	}

	r.parserState = done
	r.RequestLine = *requestLine
	fmt.Println(r.RequestLine)
	return n, nil
}

func parseRequestLine(data []byte) (*RequestLine, int, error) {
	idx := bytes.Index(data, []byte(crlf))
	if idx == -1 {
		return nil, 0, nil
	}

	requestLineText := string(data[:idx])
	requestLine, err := parseRequestLineFromString(requestLineText)
	if err != nil {
		return nil, idx, err
	}
	return requestLine, idx, nil
}

func parseRequestLineFromString(str string) (*RequestLine, error) {
	parts := strings.Split(str, " ")
	if len(parts) != 3 {
		return nil, fmt.Errorf("poorly formatted request-line: %s", str)
	}

	method := parts[0]
	for _, c := range method {
		if c < 'A' || c > 'Z' {
			return nil, fmt.Errorf("invalid method: %s", method)
		}
	}

	requestTarget := parts[1]
	if strings.Contains(requestTarget, " ") {
		return nil, fmt.Errorf("invalid request-target: %s", requestTarget)
	}

	versionParts := strings.Split(parts[2], "/")
	if len(versionParts) != 2 {
		return nil, fmt.Errorf("malformed start-linef; %s", str)
	}

	httpPart := versionParts[0]
	if httpPart != "HTTP" {
		return nil, fmt.Errorf("unrecognized HTTP-version: %s", httpPart)
	}

	version := versionParts[1]
	if version != "1.1" {
		return nil, fmt.Errorf("unrecognized HTTP-version: %s", version)
	}

	return &RequestLine{
		Method:        method,
		RequestTarget: requestTarget,
		HttpVersion:   versionParts[1],
	}, nil
}
