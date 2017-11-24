package protocol

import (
	"encoding/binary"
	"io"
)

type message struct {
	Port byte
	Data uint16
}

func (m message) encode(writer io.Writer) error {
	return binary.Write(writer, binary.BigEndian, m)
}

func decodeMessage(reader io.Reader, msg *message) error {
	return binary.Read(reader, binary.BigEndian, msg)
}
