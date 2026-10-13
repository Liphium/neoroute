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

// request is encoded as a msgp tuple: [id, route, data]. Order is guaranteed on the wire.
//
//msgp:tuple request
type request struct {
	Id    int    `msg:"id"`
	Route string `msg:"route"`
	Data  []byte `msg:"data"`
}

//go:generate msgp -unexported

// readRequestHeader reads the tuple header, the id and route, and the size of the data without reading the data itself.
//
// It also returns the bytes that were already buffered after the header, these are the start of the body.
func readRequestHeader(dc *msgp.Reader) (id int, route string, size int64, leftover []byte, err error) {
	n, err := dc.ReadArrayHeader()
	if err != nil {
		return 0, "", 0, nil, fmt.Errorf("%w: %v", ErrInvalidRequestFormat, err)
	}
	if n != 3 {
		return 0, "", 0, nil, fmt.Errorf("%w: expected 3 fields, got %d", ErrInvalidRequestFormat, n)
	}
	id, err = dc.ReadInt()
	if err != nil {
		return 0, "", 0, nil, fmt.Errorf("%w: %v", ErrInvalidRequestFormat, err)
	}
	route, err = dc.ReadString()
	if err != nil {
		return 0, "", 0, nil, fmt.Errorf("%w: %v", ErrInvalidRequestFormat, err)
	}
	dataSize, err := dc.ReadBytesHeader()
	if err != nil {
		return 0, "", 0, nil, fmt.Errorf("%w: %v", ErrInvalidRequestFormat, err)
	}
	// Take the bytes the reader already buffered ahead of the header, they belong to the body
	buffered, err := dc.R.Next(dc.Buffered())
	if err != nil {
		return 0, "", 0, nil, fmt.Errorf("%w: %v", ErrInvalidRequestFormat, err)
	}
	// Copy, the slice points into the reader's internal buffer
	leftover = append([]byte(nil), buffered...)
	return id, route, int64(dataSize), leftover, nil
}

// ensureBodyEnded checks that r has no bytes left after the body, reading at most one byte.
func ensureBodyEnded(r io.Reader) error {
	var one [1]byte
	n, err := r.Read(one[:])
	if n > 0 {
		return fmt.Errorf("%w: trailing bytes after request body", ErrInvalidRequestFormat)
	}
	if err != nil && !errors.Is(err, io.EOF) {
		return fmt.Errorf("%w: %v", ErrReadingBody, err)
	}
	return nil
}

// readRequestData reads the data of a request with the given size.
func readRequestData(r io.Reader, size int64) ([]byte, error) {
	data := make([]byte, size)
	if _, err := io.ReadFull(r, data); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrReadingBody, err)
	}
	return data, nil
}
