package protocol

import (
	"fmt"
)

const (
	portMax      = 128
	portWordSize = 16
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
	if raw < portMax {
		return DataPort{BasicPort{raw}}, nil
	}
	return DataPort{}, fmt.Errorf("invalid dataport %v", raw)
}

func (p DataPort) Byte() byte {
	return p.raw
}

func (p DataPort) Control() (ControlPort, ControlBitmask) {
	port := ControlPort{BasicPort{portMax + p.raw/portWordSize}}
	bitmask := ControlBitmask(1 << (p.raw % deviceWordSize))
	return port, bitmask
}

type ControlPort struct {
	BasicPort
}

type ControlBitmask uint16
