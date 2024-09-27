package rongta

import (
	"io"
	"runtime"

	"go.bug.st/serial"
)

type ConnType int

const (
	USB ConnType = iota
	Serial
	TCP
)

type Config interface {
	Default()
	connect() (io.ReadWriteCloser, error)
}

type SerialConfig struct {
	Port     string
	BaudRate int
	Parity   serial.Parity
	DataBits int
	StopBits serial.StopBits
}

func (c *SerialConfig) Default() {
	c.BaudRate = 19200
	c.Parity = serial.NoParity
	c.DataBits = 8
	c.StopBits = serial.OneStopBit

	switch runtime.GOOS {
	case "windows":
		c.Port = "COM1"
	default:
		c.Port = "/dev/ttyUSB0"
	}
}

// Connect to the serial port
// Returns a ReadWriteCloser interface
func (c *SerialConfig) connect() (io.ReadWriteCloser, error) {
	mode := &serial.Mode{
		BaudRate: c.BaudRate,
		Parity:   c.Parity,
		DataBits: c.DataBits,
		StopBits: c.StopBits,
	}

	s, err := serial.Open(c.Port, mode)
	if err != nil {
		return nil, err
	}

	return s, nil
}

// Unimplemented
type USBConfig struct {
	// Todo
}

func (c *USBConfig) Default() {
	panic("unimplemented")
}

func (c *USBConfig) connect() (io.ReadWriteCloser, error) {
	panic("unimplemented")
}

// Unimplemented
type TCPConfig struct {
	// Todo
}

func (c *TCPConfig) Default() {
	panic("unimplemented")
}

func (c *TCPConfig) connect() (io.ReadWriteCloser, error) {
	panic("unimplemented")
}
