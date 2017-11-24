package protocol

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMessageEncode(t *testing.T) {
	buf := bytes.NewBuffer(make([]byte, 0, 3))
	msg := Message{Port: 0x42, Data: 0x3040}
	err := msg.Encode(buf)

	assert.NoError(t, err)
	assert.Equal(t, []byte{0x42, 0x30, 0x40}, buf.Bytes())
}

func TestMessageDecode(t *testing.T) {
	buf := []byte{0x42, 0x30, 0x40}
	var msg Message
	err := DecodeMessage(bytes.NewReader(buf), &msg)

	assert.NoError(t, err)
	assert.Equal(t, byte(0x42), msg.Port)
	assert.Equal(t, uint16(0x3040), msg.Data)
}
