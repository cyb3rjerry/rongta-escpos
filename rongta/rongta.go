package rongta

import (
	"github.com/cyb3rjerry/rongta-escpos/commands"
)

var ()

type Printer struct {
	driver *commands.Driver
}

// Requires a config struct to initialize the printer
// Default values can be initiated by calling the config.Default() method
func New(config Config) (*Printer, error) {

	p := &Printer{}

	rwc, err := config.connect()
	if err != nil {
		return nil, err
	}

	p.driver = commands.NewDriver(rwc)

	return p, nil
}

func (p *Printer) Init() error {

	p.driver.Initialize()
	return nil
}
