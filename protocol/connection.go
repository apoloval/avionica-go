package protocol

import (
	"encoding/binary"
	"fmt"
	"io"
)

type connectionHeader struct {
	versionMajor byte
	versionMinor byte
}

func decodeConnectionHeader(reader io.Reader, header *connectionHeader) error {
	magic := make([]byte, 8, 8)
	_, err := reader.Read(magic)
	if err != nil {
		return err
	}
	if string(magic) != "AVIONICA" {
		return fmt.Errorf("invalid connection header: %s", magic)
	}

	return binary.Read(reader, binary.BigEndian, header)
}
