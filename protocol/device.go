package protocol

import (
	"io"
	"sync"

	"github.com/op/go-logging"
)

const deviceReadChannelBufferSize = 64

var logger = logging.MustGetLogger("avionica.device")

type Device struct {
	raw        io.ReadWriteCloser
	portConfig PortConfig
	handlers   messageHandlers
	mutex      sync.Mutex

	versionMajor byte
	versionMinor byte
}

func NewDevice(raw io.ReadWriteCloser) (*Device, error) {
	device := &Device{
		raw:        raw,
		portConfig: NewPortConfig(),
		handlers:   newMessageHandlers(),
	}

	if err := device.connect(); err != nil {
		return nil, err
	}

	go device.readLoop()

	return device, nil
}

func (dev *Device) Version() (major byte, minor byte) {
	return dev.versionMajor, dev.versionMinor
}

func (dev *Device) ConfigurePorts(config PortConfig) error {
	for rawPort := Port(portControlFirst); rawPort <= portControlLast; rawPort++ {
		port := ControlPort{rawPort}
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

func (dev *Device) AddHandler(port Port) chan MessageData {
	dev.mutex.Lock()
	defer dev.mutex.Unlock()

	return dev.handlers.addHandler(port)
}

func (dev *Device) Close() error {
	return dev.raw.Close()
}

func (dev *Device) connect() error {
	var conn connectionHeader
	if err := decodeConnectionHeader(dev.raw, &conn); err != nil {
		return err
	}

	dev.versionMajor = conn.VersionMajor
	dev.versionMinor = conn.VersionMinor

	return nil
}

func (dev *Device) readLoop() {
	var msg Message
	for {
		if err := decodeMessage(dev.raw, &msg); err != nil {
			logger.Errorf("Failed to decode Message from device: %s", err)
			dev.handlers.closeAll()
			return
		}
		dev.mutex.Lock()
		dev.handlers.handle(msg)
		dev.mutex.Unlock()
	}
}

type messageHandlers map[Port][]chan MessageData

func newMessageHandlers() messageHandlers {
	return make(messageHandlers)
}

func (mh messageHandlers) addHandler(port Port) chan MessageData {
	c := make(chan MessageData, deviceReadChannelBufferSize)
	handlers := mh.handlersFor(port)
	handlers = append(handlers, c)
	mh[port] = handlers
	return c
}

func (mh messageHandlers) handle(msg Message) {
	logger.Debugf("Handling read from port 0x%x, value 0x%x", msg.Port, msg.Data)
	handlers := mh.handlersFor(msg.Port)
	for _, handler := range handlers {
		handler <- msg.Data
	}
}

func (mh messageHandlers) closeAll() {
	for _, handlers := range mh {
		for _, handler := range handlers {
			close(handler)
		}
	}
}

func (mh messageHandlers) handlersFor(port Port) []chan MessageData {
	if handlers, found := mh[port]; found {
		return handlers
	}
	return []chan MessageData{}
}
