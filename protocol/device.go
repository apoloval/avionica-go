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
	for rawPort := Port(portControlFirst); rawPort <= portControlLast; rawPort++ {
		port := MustControlPort(rawPort)
		prevBitmask := dev.portConfig.Bitmask(port)
		newBitmask := config.Bitmask(port)
		if prevBitmask != newBitmask {
			msg := Message{Port: port.raw, Data: MessageData(newBitmask)}
			err := msg.encode(dev.raw)
			if err != nil {
				return err
			}
		}
	}
	dev.portConfig = config.Copy()
	return nil
}

func (dev *Device) Write(port DataPort, data MessageData) error {
	msg := Message{Port: port.raw, Data: data}
	return msg.encode(dev.raw)
}
