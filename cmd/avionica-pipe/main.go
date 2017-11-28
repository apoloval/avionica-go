package main

import (
	"errors"
	"fmt"
	"os"

	"strconv"

	"time"

	"github.com/apoloval/avionica-go/protocol"
	"github.com/op/go-logging"
)

func main() {
	configLogging()

	if len(os.Args) < 3 {
		exit(errors.New("invalid argument count"))
	}
	serialPortName := os.Args[1]
	dataPort := parseDataPort(os.Args[2])

	device, err := protocol.NewSerialDevice(serialPortName, 9600)
	if err != nil {
		exit(err)
	}

	// Let the Arduino to boot completely
	time.Sleep(time.Second)

	portConfig := protocol.NewPortConfig()
	portConfig.Enable(dataPort)
	if err := device.ConfigurePorts(portConfig); err != nil {
		exit(err)
	}

	output := device.AddHandler(dataPort.Basic())
	for {
		o, ok := <-output
		if !ok {
			break
		}
		fmt.Printf("0x%06x\n", o)
	}
	device.Close()
}

func parseDataPort(str string) protocol.DataPort {
	n, err := strconv.ParseInt(str, 10, 8)
	if err != nil {
		exit(err)
	}

	port, err := protocol.NewDataPort(protocol.Port(n))
	if err != nil {
		exit(err)
	}

	return port
}

func configLogging() {
	backend := logging.NewLogBackend(os.Stderr, "", 0)
	leveled := logging.AddModuleLevel(backend)
	leveled.SetLevel(logging.WARNING, "")
	logging.SetBackend(leveled)
}

func exit(err error) {
	fmt.Fprintf(os.Stderr, "Error: %s\n", err)
	fmt.Fprintf(os.Stderr, "Usage: %s <serial port> <port number>\n\n", os.Args[0])
	os.Exit(-1)
}
