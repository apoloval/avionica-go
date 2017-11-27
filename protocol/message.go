package protocol

import (
	"encoding/binary"
	"io"
)

type MessageData uint16

type Message struct {
	Port Port
	Data MessageData
}

func (m Message) encode(writer io.Writer) error {
	return binary.Write(writer, binary.BigEndian, m)
}

func decodeMessage(reader io.Reader, msg *Message) error {
	return binary.Read(reader, binary.BigEndian, msg)
}
