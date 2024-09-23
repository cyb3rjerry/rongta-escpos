package rongta

// Real-time status transmission commands
// https://www.manualslib.com/manual/3423402/Rongta-Technology-Rp325.html

type PrinterStatus string

const (
	// Printer status information bitmasks
	DRAWER_OPEN_CLOSE_STATUS_MASK uint8 = 0x04

	// Offline status information bitmasks
	COVER_STATUS_MASK       uint8 = 0x04 // 0 = Closed, 4 = Open
	FEED_BUTTON_STATUS_MASK uint8 = 0x08 // 0 = Pressed, 8 = Released

	// Error status information bitmasks
	AUTOCUTER_STATUS_MASK             uint8 = 0x08
	UNRECOVERABLE_ERROR_STATUS_MASK   uint8 = 0x20
	AUTORECOVERABLE_ERROR_STATUS_MASK uint8 = 0x40

	// Continuous paper detector status information bitmasks
	PAPER_PRESENT_STATUS_MASK uint8 = 0x60
)

// Get the status of the printer cover
// Returns true if the cover is pin 3 is HIGH, false if it's LOW
func (p *Printer) GetDrawerStatus() (bool, error) {
	status, err := p.getPrinterStatus()
	if err != nil {
		return false, err
	}

	return status&DRAWER_OPEN_CLOSE_STATUS_MASK == 0, nil
}

// Get the status of the printer cover
// Returns true if the cover is open, false if it's closed
func (p *Printer) GetCoverStatus() (bool, error) {
	status, err := p.getOfflineStatus()
	if err != nil {
		return false, err
	}

	return status&COVER_STATUS_MASK == 0, nil
}

// Get the status of the feed button
// Returns true if the feed button is pressed, false if it's released
func (p *Printer) GetFeedButtonStatus() (bool, error) {
	status, err := p.getOfflineStatus()
	if err != nil {
		return false, err
	}

	return status&FEED_BUTTON_STATUS_MASK == 0, nil
}

// Get the status of the autocutter
// Returns true if the autocutter is jammed, false if it's not
func (p *Printer) GetAutocutterStatus() (bool, error) {
	status, err := p.getErrorStatus()
	if err != nil {
		return false, err
	}

	return status&AUTOCUTER_STATUS_MASK != 0, nil
}

// Get the status of the unrecoverable error
// Returns true if there is an unrecoverable error, false if there isn't
func (p *Printer) GetUnrecoverableErrorStatus() (bool, error) {
	status, err := p.getErrorStatus()
	if err != nil {
		return false, err
	}

	return status&UNRECOVERABLE_ERROR_STATUS_MASK != 0, nil
}

// Get the status of the autorecoverable error
// Returns true if there is an autorecoverable error, false if there isn't
func (p *Printer) GetAutorecoverableErrorStatus() (bool, error) {
	status, err := p.getErrorStatus()
	if err != nil {
		return false, err
	}

	return status&AUTORECOVERABLE_ERROR_STATUS_MASK != 0, nil
}

// Transmit printer status
func (p *Printer) getPrinterStatus() (uint8, error) {
	return p.getTransmitStatus(0x01)
}

// Offline status
func (p *Printer) getOfflineStatus() (uint8, error) {
	return p.getTransmitStatus(0x02)
}

// Transmit error status
func (p *Printer) getErrorStatus() (uint8, error) {
	return p.getTransmitStatus(0x03)
}

func (p *Printer) getTransmitStatus(statusType uint8) (uint8, error) {
	status := make([]byte, 1)

	_, err := p.rwc.Write([]byte{DLE, EOT, statusType})
	if err != nil {
		return 0, err
	}

	_, err = p.rwc.Read(status)
	return status[0], err
}
