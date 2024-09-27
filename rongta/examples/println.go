package main

import "github.com/cyb3rjerry/rongta-escpos/rongta"

func main() {

	config := &rongta.SerialConfig{}
	config.Default()

	printer, err := rongta.New(config)
	if err != nil {
		panic(err)
	}

	err = printer.Init()
	if err != nil {
		panic(err)
	}

	err = printer.Println("Hello, World!")
	if err != nil {
		panic(err)
	}

}
