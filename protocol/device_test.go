package protocol

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDevice_EnablePort(t *testing.T) {
	buf := closeableBuffer{bytes.NewBuffer([]byte{})}
	dev := NewDevice(buf)
	config := NewPortConfig()
	config.Enable(MustDataPort(0))
	config.Enable(MustDataPort(1))
	config.Enable(MustDataPort(16))
	err := dev.ConfigurePorts(config)

	assert.NoError(t, err)
	assert.Equal(t, []byte{0x80, 0x00, 0x03, 0x81, 0x00, 0x01}, buf.Bytes())
}

type closeableBuffer struct{ *bytes.Buffer }

func (closeableBuffer) Close() error { return nil }
