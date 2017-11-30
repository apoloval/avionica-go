package protocol

import (
	"github.com/tarm/serial"
)

func NewSerialDevice(portName string, baudRate int) (*Device, error) {
	config := serial.Config{
		Name: portName,
		Baud: baudRate,
	}
	port, err := serial.OpenPort(&config)
	if err != nil {
		return nil, err
	}

	if err := port.Flush(); err != nil {
		return nil, err
	}

	return NewDevice(port)
}
