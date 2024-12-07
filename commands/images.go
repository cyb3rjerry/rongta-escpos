package commands

import (
	"bytes"
	"errors"
	"fmt"
)

var (
	ErrInvalidBitImageModevalue = errors.New("invalid m value")
)

// SetPageMode switches the printer to page mode
func (p *Driver) SetPageMode() error {
	_, err := p.rwc.Write([]byte{0x1B, 0x4C}) // ESC L
	return err
}

// SetStandardMode switches the printer to standard mode
func (p *Driver) SetStandardMode() error {
	_, err := p.rwc.Write([]byte{0x1B, 0x53}) // ESC S
	return err
}

// PrintPage prints the content in page mode
func (p *Driver) PrintPage() error {
	_, err := p.rwc.Write([]byte{0x0C}) // FF (Form Feed)
	return err
}

// PrintRasterBitImage prints a raster bit image using the GS v 0 command
//
// m = 0, 1, 32, 33
// m = 0: Normal mode
// m = 1: Double width mode
// m = 2: Double height mode
// m = 3: Quadruple mode
func (p *Driver) PrintRasterBitImage(m, width, height int, data []byte) error {
	if m != 0 && m != 1 && m != 2 && m != 3 {
		return fmt.Errorf("invalid raster bit-image mode value: %d", m)
	}

	expectedDataLength := (width + 7) / 8 * height
	if len(data) != expectedDataLength {
		return fmt.Errorf("data length does not match width and height")
	}

	var buffer bytes.Buffer

	// GS v 0 m xL xH yL yH d1...dk
	buffer.WriteByte(0x1D)    // GS
	buffer.WriteByte('v')     // v
	buffer.WriteByte('0')     // 0
	buffer.WriteByte(byte(m)) // m

	// Width in bytes (xL and xH)
	xL := byte((width + 7) / 8 % 256)
	xH := byte((width + 7) / 8 / 256)
	buffer.WriteByte(xL)
	buffer.WriteByte(xH)

	// Height in bytes (yL and yH)
	yL := byte(height % 256)
	yH := byte(height / 256)
	buffer.WriteByte(yL)
	buffer.WriteByte(yH)

	// Write image data
	buffer.Write(data)

	// Send the command to the printer
	_, err := p.rwc.Write(buffer.Bytes())
	if err != nil {
		return fmt.Errorf("failed to write to printer: %w", err)
	}

	return nil
}
