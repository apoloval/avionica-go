package protocol

import (
	"fmt"
)

const (
	portDataMax      = 128
	portDataLast     = portDataMax - 1
	portWordSize     = 16
	portControlFirst = portDataMax
	portControlMax   = portDataMax / portWordSize
	portControlLast  = portControlFirst + portControlMax - 1
)

type BasicPort struct {
	raw byte
}

func (p BasicPort) Byte() byte {
	return p.raw
}

type DataPort struct {
	BasicPort
}

func NewDataPort(raw byte) (DataPort, error) {
	if raw <= portDataLast {
		return DataPort{BasicPort{raw}}, nil
	}
	return DataPort{}, fmt.Errorf("invalid dataport %v", raw)
}

func MustDataPort(raw byte) DataPort {
	port, err := NewDataPort(raw)
	if err != nil {
		panic(err)
	}
	return port
}

func (p DataPort) Byte() byte {
	return p.raw
}

func (p DataPort) Control() (ControlPort, ControlBitmask) {
	port := ControlPort{BasicPort{portControlFirst + p.raw/portWordSize}}
	bitmask := ControlBitmask(1 << (p.raw % portWordSize))
	return port, bitmask
}

type ControlPort struct {
	BasicPort
}

func NewControlPort(raw byte) (ControlPort, error) {
	if raw >= portControlFirst && raw <= portControlLast {
		return ControlPort{BasicPort{raw}}, nil
	}
	return ControlPort{}, fmt.Errorf("invalid controlport %v", raw)
}

func MustControlPort(raw byte) ControlPort {
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
	index := cp.Byte() - portControlFirst
	return pc.raw[index]&cbm > 0
}

func (pc PortConfig) Enable(port DataPort) {
	cp, cbm := port.Control()
	index := cp.Byte() - portControlFirst
	pc.raw[index] |= cbm
}

func (pc PortConfig) Disable(port DataPort) {
	cp, cbm := port.Control()
	index := cp.Byte() - portControlFirst
	pc.raw[index] &= ^cbm
}

func (pc PortConfig) Bitmask(port ControlPort) ControlBitmask {
	index := port.Byte() - portControlFirst
	return pc.raw[index]
}

func (pc PortConfig) Copy() PortConfig {
	result := NewPortConfig()
	copy(result.raw, pc.raw)
	return result
}
