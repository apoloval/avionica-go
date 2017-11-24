package protocol

import (
	"encoding/binary"
	"io"
)

type Message struct {
	Port byte
	Data uint16
}

func (m Message) Encode(writer io.Writer) error {
	return binary.Write(writer, binary.BigEndian, m)
}

func DecodeMessage(reader io.Reader, msg *Message) error {
	return binary.Read(reader, binary.BigEndian, msg)
}
