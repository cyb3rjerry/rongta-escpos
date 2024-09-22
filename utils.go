package rongta

import "errors"

type COMMAND []byte
type PrinterIDInfo uint8
type PaperStatus uint8

const (
	PrinterModelID  PrinterIDInfo = 0x01
	PrinterTypeID   PrinterIDInfo = 0x02
	FirmwareVersion PrinterIDInfo = 0x41
	ManufacturerID  PrinterIDInfo = 0x42
	PrinterName     PrinterIDInfo = 0x43
	SerialNumber    PrinterIDInfo = 0x44
	// TypeOfMountedAdditionalFonts PrinterIDInfo = 0x45 // TODO: Currently unsupported, requires experimentation

	PaperStatusOK  PaperStatus = 0x00
	PaperStatusLow PaperStatus = 0x0C
)

var (
	ErrInvalidNumberOfBeeps = errors.New("invalid number of beeps")
	ErrInvalidBeepTime      = errors.New("invalid beep time")
	ErrInvalidTypePrinterID = errors.New("invalid type of printer ID requested")
)

// Recovers from a recoverable error and restarts printing from the line where the
// error occurred
// This command is effective only when an auto cutter error, a
// BM detecting error or a platen-open error occurs.
// This command is executed even when the printer is offline,
// the receive buffer is full, or there is an error status with a
// serial interface model.
// With a parallel interface model, this command can’t be
// executed when the printer is busy.
func (p *Printer) RecoverAndRestartPrint() error {
	_, err := p.rwc.Write([]byte{ESC, ENQ, 0x01})
	return err
}

// Recovers from a recoverable error after clearing the receive and print buffers
// This command is effective only when an auto cutter error, a
// BM detecting error or a platen-open error occurs.
// This command is executed even when the printer is offline,
// the receive buffer is full, or there is an error status with a
// serial interface model.
// With a parallel interface model, this command can’t be
// executed when the printer is busy.
func (p *Printer) RecoverAndCancelPrint() error {
	_, err := p.rwc.Write([]byte{ESC, ENQ, 0x02})
	return err
}

// Generate a pulse at real-time to either pin 2 or pin 5
// The pulse width is 100 ms
// m = false: Pin 2
// m = true: Pin 5
// t = time X 100ms
func (p *Printer) SendPulseToPin(m bool, t uint8) error {
	var pin byte
	if m {
		pin = 0x05
	} else {
		pin = 0x02
	}

	if (t > 0x08) || (t < 0x01) {
		return ErrInvalidPulseTime
	}

	_, err := p.rwc.Write([]byte{ESC, BANG, pin, t})
	return err
}

// Set beep prompt
// Only for page mode and general 347
// n: number of beeps (1 <= n <= 9)
// t: time of each beep (1 <= t <= 9)
func (p *Printer) SetBeepPrompt(n, t uint8) error {
	if n < 1 || n > 9 {
		return ErrInvalidNumberOfBeeps
	}

	if t < 1 || t > 9 {
		return ErrInvalidBeepTime
	}

	_, err := p.rwc.Write([]byte{ESC, 'B', n, t})
	return err
}

// Generate pulse
// m = 0: Drawer kick out pin 2
// m = 1: Drawer kick out pin 5
// on time = t1 X 2ms
// off time = t2 X 2ms
func (p *Printer) GeneratePulse(m bool, t1, t2 uint8) error {
	var pin byte
	if m {
		pin = 0x05
	} else {
		pin = 0x02
	}

	_, err := p.rwc.Write([]byte{ESC, 'p', pin, t1, t2})
	return err
}

// Disable/Enable pannel buttons.
// When the LSB of n is 0, the panel buttons are enabled.
// When the LSB of n is 1, the panel buttons are disabled.
func (p *Printer) DisablePanelButtons(n uint8) error {
	_, err := p.rwc.Write([]byte{ESC, 'c', '5', n})
	return err
}

// Cut paper (only partial is supported)
func (p *Printer) Cut() error {
	_, err := p.rwc.Write([]byte{ESC, 'i'})
	return err
}

// Transmit printer ID
// n = 1: Transmit the printer ID
// n = 2: Transmit the printer ID and the firmware version
// n = 65: Firmware version
// n = 66: Printer ID
func (p *Printer) TransmitPrinterID(n PrinterIDInfo) ([]byte, error) {
	buf := make([]byte, 1024)

	switch n {
	case PrinterModelID | PrinterTypeID:
		_, err := p.rwc.Write([]byte{ESC, 'i', 1})
		if err != nil {
			return []byte{}, err
		}

		// Read the response
		bytesRead, err := p.rwc.Read(buf)
		if err != nil {
			return []byte{}, err
		}

		return buf[:bytesRead], nil

	case FirmwareVersion | ManufacturerID | PrinterName | SerialNumber:
		bytesRead, err := p.rwc.Write([]byte{ESC, 'i', 2})
		if err != nil {
			return []byte{}, err
		}

		return buf[:bytesRead], nil
	default:
		return []byte{}, ErrInvalidTypePrinterID
	}
}

// Toggle macro definition
func (p *Printer) ToggleMacroDefinition() error {
	_, err := p.rwc.Write([]byte{GS, ':'})
	return err
}

// Execute macro
// r: number of times to execute the macro (1 <= r <= 255)
// t: specificies the waiting time for executing the macro (0 <= t <= 255)
// m: specifies the macro number (0 <= m <= 1)
// After waiting for the period specified by t, the PAPER OUT LED
// indicators blink and the printer waits for the FEED button to be
// pressed. After the button is pressed, the printer executes the
// macro once. The printer repeats the operation r times.
// The waiting time is t x 100ms.
func (p *Printer) ExecuteMacro(r, t, m uint8) error {
	_, err := p.rwc.Write([]byte{GS, '^', r, t, m})
	return err
}

// Toggle ASB
// Bit 0: Undefined
// Bit 1: Undefined
// Bit 2: 0 = Error status disabled, 4: Error status enabled
// Bit 3: 0 = Paper sensor disabled, 8: Paper sensor enabled
// Bit 4-7: Undefined
// TODO: Define types for the bits
func (p *Printer) ToggleASB(n uint8) error {
	_, err := p.rwc.Write([]byte{GS, 'a', n})
	return err
}

// Transmit status
func (p *Printer) TransmitStatus() (PaperStatus, error) {
	_, err := p.rwc.Write([]byte{GS, 'r', 1})
	if err != nil {
		return PaperStatusLow, err
	}

	// Read status
	buf := make([]byte, 1)
	_, err = p.rwc.Read(buf)

	return PaperStatus(buf[0] & 0x0C), err
}
