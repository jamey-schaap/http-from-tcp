package headers

import (
	"bytes"
	"fmt"
	"strings"
)

type Headers map[string]string

const crlf = "\r\n"

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	idx := bytes.Index(data, []byte(crlf))
	if idx == -1 {
		return 0, false, nil
	}

	if idx == 0 {
		return /* CRLF */ 2, true, nil
	}

	parts := bytes.SplitAfterN(data[:idx], []byte(":"), 2)
	key := string(parts[0])
	if key != strings.TrimRight(key, " ") {
		return 0, false, fmt.Errorf("invalid header name: %s", key)
	}

	value := bytes.TrimSpace(parts[1])
	key = strings.TrimSpace(key)

	h[key] = string(value)
	return idx + /* CRLF */ 2, done, nil
}

func NewHeaders() Headers {
	return map[string]string{}
}
