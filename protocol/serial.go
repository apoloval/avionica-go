package protocol

import "github.com/jacobsa/go-serial/serial"

func NewSerialDevice(portName string, baudRate uint) (*Device, error) {
	opts := serial.OpenOptions{
		PortName: portName,
		BaudRate: baudRate,
	}
	port, err := serial.Open(opts)
	if err != nil {
		return nil, err
	}
	return NewDevice(port), nil
}
