package commands

import "errors"

type BARCODESYSTEM uint8

const (
	UPCA BARCODESYSTEM = iota + 64
	UPCE
	EAN13
	EAN8
	CODE39
	ITF
	CODABAR
	CODE93
	CODE128
)

var (
	ErrinvalidHRICharacterFont = errors.New("invalid HRI character font")
	ErrInvalidBarCodeMode      = errors.New("invalid bar code mode")
	ErrInvalidBarCodeLength    = errors.New("invalid bar code length")
	ErrInvalidBarCodeChar      = errors.New("invalid bar code data")
	ErrBarcodeLengthMismatch   = errors.New("barcode length mismatch")
	ErrInvalidBarCodeWidth     = errors.New("invalid bar code width")
)

// Select font for Human Readable Interpretation (HRI) characters
// n = 0: Font A
// n = 1: Font B
func (p *Driver) SelectFontForHRICharacters(n uint8) error {
	if n != 0 && n != 1 {
		return ErrinvalidHRICharacterFont
	}
	_, err := p.rwc.Write([]byte{GS, 'f', n})
	return err
}

// Selects the printing position of HRI characters when printing
// a bar code. n selects the printing position as follows:
// n = 0: Not printed
// n = 1: Above the bar code
// n = 2: Below the bar code
// n = 3: Both above and below the bar code
func (p *Driver) SelectHRICharacterPrintPosition(n uint8) error {
	_, err := p.rwc.Write([]byte{GS, 'H', n})
	return err
}

// [Incomplete] Currently doesn't handle special characters
// Print bar code mode 0
// m = bar code system
// n = the number of bar code data bytes
// d1...dn = bar code data
func (p *Driver) PrintBarCode(n uint8, m BARCODESYSTEM, d []uint8) error {

	switch m {
	case UPCA:
		if n < 11 || n > 12 {
			return ErrInvalidBarCodeLength
		}
		if len(d) != int(n) {
			return ErrBarcodeLengthMismatch
		}
		for _, v := range d {
			if v < 48 || v > 57 {
				return ErrInvalidBarCodeChar
			}
		}

	case UPCE:
		if n < 11 || n > 12 {
			return ErrInvalidBarCodeLength
		}
		if len(d) != int(n) {
			return ErrBarcodeLengthMismatch
		}
		for _, v := range d {
			if v < 48 || v > 57 {
				return ErrInvalidBarCodeChar
			}
		}

	case EAN13:
		if n < 12 || n > 13 {
			return ErrInvalidBarCodeLength
		}
		if len(d) != int(n) {
			return ErrBarcodeLengthMismatch
		}
		for _, v := range d {
			if v < 48 || v > 57 {
				return ErrInvalidBarCodeChar
			}
		}

	case EAN8:
		if n < 7 || n > 8 {
			return ErrInvalidBarCodeLength
		}
		if len(d) != int(n) {
			return ErrBarcodeLengthMismatch
		}
		for _, v := range d {
			if v < 48 || v > 57 {
				return ErrInvalidBarCodeChar
			}
		}

	case CODE39:
		if n < 1 {
			return ErrInvalidBarCodeLength
		}
		if len(d) != int(n) {
			return ErrBarcodeLengthMismatch
		}
		for _, v := range d {
			// Find a better way to comply to these conditions
			// This fucking sucks
			if (v < 45 && v != 32 && v != 36 && v != 37 && v != 43) || (v > 57 && v != 58 && v < 65) || v > 90 {
				return ErrInvalidBarCodeChar
			}
		}

	case ITF:
		if n < 1 && (n%2 != 0) {
			return ErrInvalidBarCodeLength
		}
		if len(d) != int(n) {
			return ErrBarcodeLengthMismatch
		}
		for _, v := range d {
			if v < 48 || v > 57 {
				return ErrInvalidBarCodeChar
			}
		}

	case CODABAR:
		if n < 1 {
			return ErrInvalidBarCodeLength
		}
		if len(d) != int(n) {
			return ErrBarcodeLengthMismatch
		}
		for _, v := range d {
			// Find a better way to comply to these conditions
			// This fucking sucks
			if (v < 45 && v != 43 && v != 36) || (v > 57 && v != 58 && v < 68) {
				return ErrInvalidBarCodeChar
			}
		}

	case CODE93:
		if n < 1 {
			return ErrInvalidBarCodeLength
		}
		if len(d) != int(n) {
			return ErrBarcodeLengthMismatch
		}
		for _, v := range d {
			if v > 127 {
				return ErrInvalidBarCodeChar
			}
		}

	case CODE128:
		if n < 1 {
			return ErrInvalidBarCodeLength
		}
		if len(d) != int(n) {
			return ErrBarcodeLengthMismatch
		}
		for _, v := range d {
			if v > 127 {
				return ErrInvalidBarCodeChar
			}
		}

	default:
		return ErrInvalidBarCodeMode
	}

	command := []byte{GS, 'k', uint8(m), n}
	command = append(command, d...)
	_, err := p.rwc.Write(command)

	return err
}

// Sets the horizontal size of the bar code. n
// specifies the bar code width as follows:
// n | Module width (mm) for Multi-level bar code | thin element width | thick element width (for binary level bar code)
//
// 2 | 0.250 | 0.250 | 0.625
// 3 | 0.375 | 0.375 | 1.000
// 4 | 0.560 | 0.500 | 1.250
// 5 | 0.625 | 0.625 | 1.625
// 6 | 0.750 | 0.750 | 2.000
//
// Multi-level bar codes: UPC-A, UPC-E, EAN13, EAN8, CODE93, CODE128
// Binary level bar codes: CODE39, ITF, CODABAR
// The default value is 3.
func (p *Driver) SetBarcodeWidth(n uint8) error {
	if n < 2 || n > 6 {
		return ErrInvalidBarCodeWidth
	}
	_, err := p.rwc.Write([]byte{GS, 'w', n})
	return err
}

// Sets the printing position of the bar code.
// The print bar code starting position is: 0->255
func (p *Driver) SetBarCodePrintPosition(n uint8) error {
	_, err := p.rwc.Write([]byte{GS, 'x', n})
	return err
}

// Print 2D barcode
// IMPORTANT: This command relies on the previously set 2D barcode mode via
// `GS, Z, m` command
//
// PDF417 (mode 0)
// m: specifies the column number (1 <= m <= 30), when m = 0, the column number is automatically set
// n: specifies security level to restore when barcode image is damaged
// k: specifies the horizontal and vertical ratio (2 <= k <= 5)
// dL: Lower number
// dH: Higher number
// d1..dn: the data to be printed
//
// QR Code (mode 1)
// m: specifies version (1 <= m <= 40, 0 = autosize)
// n: specifies error correction level (n = 'L' | 'M' | 'Q' | 'H')
// k: specifies module size (1 <= k <= 8)
// dL: Lower number
// dH: Higher number
// d1..dn: the data to be printed
func (p *Driver) PrintQRBarcode(m, n, k, dL, dH uint8, d []uint8) error {
	_, err := p.rwc.Write(append([]byte{ESC, 'Z', m, n, k, dL, dH}, d...))

	return err
}

// Select 2D barcode mode
// m = 0: PDF417
// m = 1: QR code
// Default: 0
func (p *Driver) Select2DBarcodeMode(m uint8) error {
	if m != 0 && m != 1 {
		return ErrInvalidBarCodeMode
	}
	_, err := p.rwc.Write([]byte{GS, 'Z', m})
	return err
}
