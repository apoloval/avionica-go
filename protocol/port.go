package protocol

import (
	"fmt"
)

const (
	portDataMax    = 128
	portWordSize   = 16
	portControlMax = portDataMax / portWordSize

	portDataLast     = Port(portDataMax - 1)
	portControlFirst = Port(portDataMax)
	portControlLast  = Port(portControlFirst + portControlMax - 1)
)

type Port byte

type DataPort struct {
	raw Port
}

func NewDataPort(raw Port) (DataPort, error) {
	if raw <= portDataLast {
		return DataPort{raw}, nil
	}
	return DataPort{}, fmt.Errorf("invalid dataport %v", raw)
}

func MustDataPort(raw Port) DataPort {
	port, err := NewDataPort(raw)
	if err != nil {
		panic(err)
	}
	return port
}

func (p DataPort) Basic() Port {
	return p.raw
}

func (p DataPort) Control() (ControlPort, ControlBitmask) {
	port := ControlPort{Port(portControlFirst + p.raw/portWordSize)}
	bitmask := ControlBitmask(1 << (p.raw % portWordSize))
	return port, bitmask
}

type ControlPort struct {
	raw Port
}

func NewControlPort(raw Port) (ControlPort, error) {
	if raw >= portControlFirst && raw <= portControlLast {
		return ControlPort{Port(raw)}, nil
	}
	return ControlPort{}, fmt.Errorf("invalid controlport %v", raw)
}

func MustControlPort(raw Port) ControlPort {
	port, err := NewControlPort(raw)
	if err != nil {
		panic(err)
	}
	return port
}

type ControlBitmask uint16

type PortConfig struct {
	raw []ControlBitmask
}

func NewPortConfig() PortConfig {
	return PortConfig{make([]ControlBitmask, portControlMax)}
}

func (pc PortConfig) IsEnabled(port DataPort) bool {
	cp, cbm := port.Control()
	index := cp.raw - portControlFirst
	return pc.raw[index]&cbm > 0
}

func (pc PortConfig) Enable(port DataPort) {
	cp, cbm := port.Control()
	index := cp.raw - portControlFirst
	pc.raw[index] |= cbm
}

func (pc PortConfig) Disable(port DataPort) {
	cp, cbm := port.Control()
	index := cp.raw - portControlFirst
	pc.raw[index] &= ^cbm
}

func (pc PortConfig) Bitmask(port ControlPort) ControlBitmask {
	index := port.raw - portControlFirst
	return pc.raw[index]
}

func (pc PortConfig) Copy() PortConfig {
	result := NewPortConfig()
	copy(result.raw, pc.raw)
	return result
}
