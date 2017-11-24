package protocol

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewDataPort(t *testing.T) {
	tests := []struct {
		raw   byte
		valid bool
	}{
		{raw: 0, valid: true},
		{raw: 100, valid: true},
		{raw: 127, valid: true},
		{raw: 128, valid: false},
		{raw: 255, valid: false},
	}
	for _, test := range tests {
		p, err := NewDataPort(test.raw)
		if test.valid {
			assert.NoError(t, err)
			assert.Equal(t, test.raw, p.Byte())
		} else {
			assert.Error(t, err)
		}
	}
}

func TestDataPort_Control(t *testing.T) {
	tests := []struct {
		raw     byte
		control byte
		bitmask ControlBitmask
	}{
		{raw: 0, control: 128, bitmask: 0x0001},
		{raw: 7, control: 128, bitmask: 0x0080},
		{raw: 15, control: 128, bitmask: 0x8000},
		{raw: 16, control: 129, bitmask: 0x0001},
	}
	for _, test := range tests {
		p, _ := NewDataPort(test.raw)
		ctrl, bitmask := p.Control()

		assert.Equal(t, test.control, ctrl.Byte())
		assert.Equal(t, test.bitmask, bitmask)
	}
}
