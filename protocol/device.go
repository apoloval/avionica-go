package protocol

import (
	"io"
)

type Device struct {
	raw        io.ReadWriteCloser
	portConfig PortConfig
}

func NewDevice(raw io.ReadWriteCloser) Device {
	return Device{raw, NewPortConfig()}
}

func (dev *Device) ConfigurePorts(config PortConfig) error {
	for rawPort := byte(portControlFirst); rawPort <= portControlLast; rawPort++ {
		port := MustControlPort(rawPort)
		prevBitmask := dev.portConfig.Bitmask(port)
		newBitmask := config.Bitmask(port)
		if prevBitmask != newBitmask {
			msg := message{Port: port.Byte(), Data: uint16(newBitmask)}
			err := msg.encode(dev.raw)
			if err != nil {
				return err
			}
		}
	}
	dev.portConfig = config.Copy()
	return nil
}

func (dev *Device) Write(port DataPort, data uint16) error {
	msg := message{Port: port.Byte(), Data: data}
	return msg.encode(dev.raw)
}
