package protocol

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDevice_Conn(t *testing.T) {
	raw := newFakeRawDevice()
	dev, err := NewDevice(raw)
	x, y := dev.Version()

	assert.NoError(t, err)
	assert.Equal(t, byte(1), x)
	assert.Equal(t, byte(21), y)
}

func TestDevice_EnablePort(t *testing.T) {
	raw := newFakeRawDevice()
	dev, _ := NewDevice(raw)
	config := NewPortConfig()
	config.Enable(MustDataPort(0))
	config.Enable(MustDataPort(1))
	config.Enable(MustDataPort(16))
	err := dev.ConfigurePorts(config)

	assert.NoError(t, err)
	assert.Equal(t, []byte{0x80, 0x00, 0x03, 0x81, 0x00, 0x01}, raw.output.Bytes())
}

func TestDevice_Write(t *testing.T) {
	raw := newFakeRawDevice()
	dev, _ := NewDevice(raw)
	err := dev.Write(MustDataPort(0x42), 0x1234)

	assert.NoError(t, err)
	assert.Equal(t, []byte{0x42, 0x12, 0x34}, raw.output.Bytes())
}

func TestDevice_AddHandler(t *testing.T) {
	raw := newFakeRawDevice()
	dev, _ := NewDevice(raw)
	handler := dev.AddHandler(0x42)
	raw.receive([]byte{0x42, 0x12, 0x34})

	assert.Equal(t, MessageData(0x1234), <-handler)
}

type fakeRawDevice struct {
	output *bytes.Buffer
	input  chan byte
}

func newFakeRawDevice() *fakeRawDevice {
	dev := &fakeRawDevice{
		output: bytes.NewBuffer([]byte{}),
		input:  make(chan byte, 256),
	}
	dev.receive([]byte("AVIONICA"))
	dev.receive([]byte{1, 21})
	return dev
}

func (raw *fakeRawDevice) receive(data []byte) {
	for _, b := range data {
		raw.input <- b
	}
}

func (raw *fakeRawDevice) Read(p []byte) (n int, err error) {
	p[0] = <-raw.input
	return 1, nil
}

func (raw *fakeRawDevice) Write(p []byte) (n int, err error) {
	return raw.output.Write(p)
}

func (fakeRawDevice) Close() error { return nil }
