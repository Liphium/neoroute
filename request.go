package neoroute

import (
	"errors"
	"fmt"
	"io"

	"github.com/tinylib/msgp/msgp"
)

// maxRequestHeaderSize is the max number of bytes read while parsing the request envelope (id, route, data size).
//
// This number might change in the future, but it should be fine to fit basically anything.
const maxRequestHeaderSize = 1024

//go:generate msgp -unexported

// request is encoded as a msgp tuple: [id, route, data]. Order is guaranteed on the wire.
//
//msgp:tuple request
type request struct {
	Id    int    `msg:"id"`
	Route string `msg:"route"`
	Data  []byte `msg:"data"`
}

// readRequestHeader parses the tuple header, the id and route, and the size of the data from the given bytes
// without reading the data itself. It returns the bytes that were not consumed by the header.
//
// If the bytes are incomplete, the returned error wraps msgp.ErrShortBytes.
func readRequestHeader(b []byte) (id int, route string, size int64, rest []byte, err error) {
	n, b, err := msgp.ReadArrayHeaderBytes(b)
	if err != nil {
		return 0, "", 0, nil, fmt.Errorf("%w: %w", ErrInvalidRequestFormat, err)
	}
	if n != 3 {
		return 0, "", 0, nil, fmt.Errorf("%w: expected 3 fields, got %d", ErrInvalidRequestFormat, n)
	}
	id, b, err = msgp.ReadIntBytes(b)
	if err != nil {
		return 0, "", 0, nil, fmt.Errorf("%w: %w", ErrInvalidRequestFormat, err)
	}
	route, b, err = msgp.ReadStringBytes(b)
	if err != nil {
		return 0, "", 0, nil, fmt.Errorf("%w: %w", ErrInvalidRequestFormat, err)
	}
	dataSize, b, err := msgp.ReadBytesHeader(b)
	if err != nil {
		return 0, "", 0, nil, fmt.Errorf("%w: %w", ErrInvalidRequestFormat, err)
	}
	return id, route, int64(dataSize), b, nil
}

// readRequestHeaderFrom reads from r in small chunks until the request header can be parsed.
// It returns the header values and the bytes read from r that belong to the data after the header.
// At most maxRequestHeaderSize bytes are read from r.
func readRequestHeaderFrom(r io.Reader) (id int, route string, size int64, leftover []byte, err error) {
	buf := make([]byte, 0, 64)
	chunk := make([]byte, 64)
	for {
		if len(buf) >= maxRequestHeaderSize {
			return 0, "", 0, nil, fmt.Errorf("%w: header exceeds %d bytes", ErrInvalidRequestFormat, maxRequestHeaderSize)
		}
		want := maxRequestHeaderSize - len(buf)
		if want > len(chunk) {
			want = len(chunk)
		}
		n, readErr := r.Read(chunk[:want])
		buf = append(buf, chunk[:n]...)
		if n > 0 {
			id, route, size, rest, parseErr := readRequestHeader(buf)
			if parseErr == nil {
				return id, route, size, rest, nil
			}
			if !errors.Is(parseErr, msgp.ErrShortBytes) {
				return 0, "", 0, nil, parseErr
			}
		}
		if readErr != nil {
			if errors.Is(readErr, io.EOF) {
				readErr = io.ErrUnexpectedEOF
			}
			return 0, "", 0, nil, fmt.Errorf("%w: %v", ErrInvalidRequestFormat, readErr)
		}
	}
}

// readRequestData reads the data of a request with the given size.
func readRequestData(r io.Reader, size int64) ([]byte, error) {
	data := make([]byte, size)
	if _, err := io.ReadFull(r, data); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrReadingBody, err)
	}
	return data, nil
}
