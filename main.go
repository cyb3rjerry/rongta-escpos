package rongta

import "io"

// ESC/POS Command Set as defined in
// https://www.manualslib.com/manual/3423402/Rongta-Technology-Rp325.html

const (
	HT        = 0x09 // Horizontal Tab
	LF        = 0x0A // Line Feed
	CR        = 0x0D // Print and carriage return
	CAN       = 0x18 // Cancel print data in page mode
	DLE       = 0x10 // Data link escape
	EOT       = 0x04 // End of transmission
	ENQ       = 0x05 // Enquiry
	SP        = 0x20 // Space
	BANG      = 0x21 // !
	DOLLAR    = 0x24 // $
	PERCENT   = 0x25 // %
	DASH      = 0x2D // -
	AMPERSAND = 0x26 // &
	ASTERISK  = 0x2A // *
	SLASH     = 0x2F // /
	BACKSLASH = 0x5C // \

	ESC = 0x1B // Escape
	GS  = 0x1D // Group separator
	NUL = 0x00 // Null
	DC2 = 0x12 // Device control 2
	FF  = 0x0C // Form feed
)

type Printer struct {
	rwc io.ReadWriteCloser
}
