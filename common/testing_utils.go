package common

import (
	"fmt"
	"io"
)

func LastID() ID {
	return lastID
}

var buf []byte

func ReaderContent(body io.Reader) (string, error) {
	if buf == nil {
		buf = make([]byte, 16384)	// should be enough for everyone
	}
	if closer, ok := body.(io.Closer); ok {
		defer closer.Close()
	}
	n, err := body.Read(buf)
	return string(buf[:n]), err
}

func TestMessage(i int) string {
	return fmt.Sprintf("test case no. %d", i)
}
