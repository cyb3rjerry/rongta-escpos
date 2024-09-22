package rongta

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
)

// Select font for Human Readable Interpretation (HRI) characters
// n = 0: Font A
// n = 1: Font B
func (p *Printer) SelectFontForHRICharacters(n uint8) error {
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
func (p *Printer) SelectHRICharacterPrintPosition(n uint8) error {
	_, err := p.rwc.Write([]byte{GS, 'H', n})
	return err
}

// Print bar code mode 0
// m = bar code system
// n = the number of bar code data bytes
// d1...dn = bar code data
func (p *Printer) PrintBarCode(n uint8, m BARCODESYSTEM, d []uint8) error {

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

	panic("TODO: Implement special characters")
	return err
}
