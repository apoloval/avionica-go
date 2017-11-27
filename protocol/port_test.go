package protocol

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewDataPort(t *testing.T) {
	tests := []struct {
		raw   Port
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
			assert.Equal(t, test.raw, p.raw)
		} else {
			assert.Error(t, err)
		}
	}
}

func TestDataPort_Control(t *testing.T) {
	tests := []struct {
		raw     Port
		control Port
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

		assert.Equal(t, test.control, ctrl.raw)
		assert.Equal(t, test.bitmask, bitmask)
	}
}

func TestPortConfig_Enable(t *testing.T) {
	cfg := NewPortConfig()

	cfg.Enable(MustDataPort(0))
	assert.Equal(t, ControlBitmask(0x0001), cfg.Bitmask(MustControlPort(128)))

	cfg.Enable(MustDataPort(1))
	assert.Equal(t, ControlBitmask(0x0003), cfg.Bitmask(MustControlPort(128)))

	cfg.Enable(MustDataPort(15))
	assert.Equal(t, ControlBitmask(0x8003), cfg.Bitmask(MustControlPort(128)))

	cfg.Enable(MustDataPort(16))
	assert.Equal(t, ControlBitmask(0x0001), cfg.Bitmask(MustControlPort(129)))

	cfg.Enable(MustDataPort(31))
	assert.Equal(t, ControlBitmask(0x8001), cfg.Bitmask(MustControlPort(129)))
}

func TestPortConfig_Disable(t *testing.T) {
	cfg := NewPortConfig()

	cfg.Enable(MustDataPort(0))
	cfg.Enable(MustDataPort(1))
	cfg.Enable(MustDataPort(2))
	cfg.Disable(MustDataPort(1))
	assert.Equal(t, ControlBitmask(0x0005), cfg.Bitmask(MustControlPort(128)))
}
