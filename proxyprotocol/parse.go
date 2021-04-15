package proxyprotocol

import (
	"bufio"
	"bytes"
	"errors"
)

var (
	sigV1 = []byte("PROXY %s %s %s %d %d\r\n")
	sigV2 = []byte{0x0D, 0x0A, 0x0D, 0x0A, 0x00, 0x0D, 0x0A, 0x51, 0x55, 0x49, 0x54, 0x0A}
)

var errNotProxyProto = errors.New("possibly not using proxy protocol")

// InvalidHeaderErr contains the parsing error as well as all data read from the reader.
type InvalidHeaderErr struct {
	error
	Read []byte
}

// Parse will parse detect and return a V1 or V2 header, otherwise InvalidHeaderErr is returned.
// Parse will
func Parse(r *bufio.Reader) (Header, error) {
	b, err := r.ReadByte()
	if err != nil {
		return nil, err
	}
	r.UnreadByte()

	switch b {
	case sigV1[0]:
		return parseV1(r)
	case sigV2[0]:
		return parseV2(r)
	}
	return nil, errNotProxyProto
}

// Detect detects if the proxy protocol used
func Detect(r *bufio.Reader) (bool, error) {
	b, err := r.ReadByte()
	if err != nil {
		return false, err
	}
	r.UnreadByte()

	if b == sigV1[0] {
		if pb, err := r.Peek(5); err == nil && string(pb) == "PROXY" {
			return true, nil
		}
	} else if b == sigV2[0] {
		if pb, err := r.Peek(len(sigV2)); err == nil && bytes.Equal(pb, sigV2) {
			return true, nil
		}
	}
	return false, nil
}
