package commands

import "errors"

type Font bool
type Underline uint8
type Justify uint8

const (
	FontA Font = false
	FontB Font = true

	UnderlineNone  Underline = 0
	UnderlineThin  Underline = 1
	UnderlineThick Underline = 2

	JustifyLeft   Justify = 0
	JustifyCenter Justify = 1
	JustifyRight  Justify = 2
)

type PrintMode struct {
	Font           Font
	IsEmphasized   bool
	IsDoubleHeight bool
	IsDoubleWidth  bool
	IsUnderline    bool
}

var (
	// Errors
	ErrInvalidCharWidth      = errors.New("invalid character width")
	ErrInvalidCharHeight     = errors.New("invalid character height")
	ErrInvalidPrintDirection = errors.New("invalid print direction")
)

// Write string to printer buffer
func (p *Driver) WriteStringToBuffer(s string) error {
	_, err := p.rwc.Write([]byte(s))
	return err
}

// Set the right-side character spacing to n X 0.125mm
func (p *Driver) SetRightSideChar(n uint8) error {
	_, err := p.rwc.Write([]byte{ESC, SP, n})
	return err
}

// Set print mode(s)
// Character font A is 12x24 dots, font B is 9x17 dots
func (p *Driver) SetPrintMode(pm *PrintMode) error {
	uint8Mode := uint8(0)

	if pm.Font == FontB {
		uint8Mode |= 0x01
	}

	if pm.IsEmphasized {
		uint8Mode |= 0x08
	}

	if pm.IsDoubleHeight {
		uint8Mode |= 0x10
	}

	if pm.IsDoubleWidth {
		uint8Mode |= 0x20
	}

	if pm.IsUnderline {
		uint8Mode |= 0x80
	}

	_, err := p.rwc.Write([]byte{ESC, BANG, uint8Mode})
	return err
}

// Set absolute print position
// nL, nH = (nL + nH * 256) X 0.125mm
func (p *Driver) SetAbsolutePrintPosition(nL, nH uint8) error {
	_, err := p.rwc.Write([]byte{ESC, '$', nL, nH})
	return err
}

// Select/Cancel user-defined character
// When the LSB of n is 0, the user-defined character set is canceled
// When the LSB of n is 1, the user-defined character set is selected
// Note: When the user-defined character set is canceled, the resident
// character set is automatically selected.
func (p *Driver) SelectUserDefinedCharacter(n uint8) error {
	_, err := p.rwc.Write([]byte{ESC, '%', n})
	return err
}

// [Unimplemented] TODO: Implement this
// Define user-defined characters
func (p *Driver) DefineUserDefinedCharacters(n uint8, data []uint8) error {
	_, err := p.rwc.Write(append([]byte{ESC, '&', n}, data...))

	panic("TODO: Implement DefineUserDefinedCharacters")
	return err
}

// Turns underline mode on or off using n
// n = 0: Turns off underline mode
// n = 1: Turns on underline mode (1-dot thick)
// n = 2: Turns on underline mode (2-dot thick)
func (p *Driver) SetUnderline(u Underline) error {
	underlineBit := uint8(0)

	switch u {
	case UnderlineNone:
		underlineBit = 0
	case UnderlineThin:
		underlineBit = 1
	case UnderlineThick:
		underlineBit = 2
	}

	_, err := p.rwc.Write([]byte{ESC, DASH, underlineBit})
	return err
}

// Select default line spacing
func (p *Driver) SetDefaultLineSpacing() error {
	_, err := p.rwc.Write([]byte{ESC, '2'})
	return err
}

// Set line spacing
// Line spacing = n X 0.125mm
func (p *Driver) SetLineSpacing(n uint8) error {
	_, err := p.rwc.Write([]byte{ESC, '3', n})
	return err
}

// Initialize the printer
func (p *Driver) Initialize() error {
	_, err := p.rwc.Write([]byte{ESC, '@'})
	return err
}

// Set horizontal tab positions
// n = 0, 1, 2, ..., 255
// k = 0, 1, 2, ..., 32
func (p *Driver) SetHorizontalTabPositions(n, k uint8) error {
	_, err := p.rwc.Write([]byte{ESC, 'D', n, k})
	return err
}

// Set emphasized mode
// When the LSB of n is 0, emphasized mode is turned off.
// When the LSB of n is 1, emphasized mode is turned on.
func (p *Driver) SetEmphasizedMode(n uint8) error {
	_, err := p.rwc.Write([]byte{ESC, 'E', n})
	return err
}

// Set double-strike mode
// When the LSB of n is 0, double-strike mode is turned off.
// When the LSB of n is 1, double-strike mode is turned on.
func (p *Driver) SetDoubleStrikeMode(n uint8) error {
	_, err := p.rwc.Write([]byte{ESC, 'G', n})
	return err
}

// Print and feed n lines
func (p *Driver) PrintAndFeedNLines(n uint8) error {
	_, err := p.rwc.Write([]byte{ESC, 'd', n})
	return err
}

// Set character font
func (p *Driver) SetCharacterFont(f Font) error {
	n := uint8(0)
	if f == FontB {
		n = 1
	}

	_, err := p.rwc.Write([]byte{ESC, 'M', n})
	return err
}

// Rotate clockwise 90 degrees mode
// True: Rotate 90 degrees clockwise
// False: Cancel rotate 90 degrees clockwise
func (p *Driver) RotateClockwise90Degrees(r bool) error {
	bit := uint8(0)

	if r {
		bit = 0x01
	}

	_, err := p.rwc.Write([]byte{ESC, 'V', bit})
	return err
}

// Set relative print position
func (p *Driver) SetRelativePrintPosition(nL, nH uint8) error {
	_, err := p.rwc.Write([]byte{ESC, BACKSLASH, nL, nH})
	return err
}

// Set justification
func (p *Driver) SetJustification(j Justify) error {
	_, err := p.rwc.Write([]byte{ESC, 'a', uint8(j)})
	return err
}

// Print and feed n lines
func (p *Driver) PrintAndFeedNDotsLines(n uint8) error {
	_, err := p.rwc.Write([]byte{ESC, 'J', n})
	return err
}

// Select character size
// 1 = normal size, 2 = double, 3 = quadruple, ...
func (p *Driver) SelectCharacterSize(w, h uint8) error {
	// TODO: There's gotta be a better way to do this
	charSizeBit := uint8(0)

	switch w {
	case 1:
		charSizeBit = 0
	case 2:
		charSizeBit = 0x10
	case 3:
		charSizeBit = 0x20
	case 4:
		charSizeBit = 0x30
	case 5:
		charSizeBit = 0x40
	case 6:
		charSizeBit = 0x50
	case 7:
		charSizeBit = 0x60
	case 8:
		charSizeBit = 0x70
	default:
		return errors.New("invalid character width")
	}

	switch h {
	case 1:
		charSizeBit |= 0
	case 2:
		charSizeBit |= 0x01
	case 3:
		charSizeBit |= 0x02
	case 4:
		charSizeBit |= 0x03
	case 5:
		charSizeBit |= 0x04
	case 6:
		charSizeBit |= 0x05
	case 7:
		charSizeBit |= 0x06
	case 8:
		charSizeBit |= 0x07
	}

	_, err := p.rwc.Write([]byte{GS, '!', charSizeBit})
	return err
}

// Turn white/black reverse printing mode
// When the LSB of n is 0, white/black reverse printing mode is turned off.
// When the LSB of n is 1, white/black reverse printing mode is turned on.
func (p *Driver) SetWhiteBlackReversePrintingMode(n uint8) error {
	_, err := p.rwc.Write([]byte{GS, 'B', n})
	return err
}

// Set left margin
// nL, nH = (nL + nH * 256) X 0.125mm
func (p *Driver) SetLeftMargin(nL, nH uint8) error {
	_, err := p.rwc.Write([]byte{GS, 'L', nL, nH})
	return err
}

// Select cut mode and cut paper to cutting position n
// Feeds paper (cutting position + [n x 0.125mm])
func (p *Driver) SelectCutModeAndCutPaper(n uint8) error {
	_, err := p.rwc.Write([]byte{GS, 'V', 0x66, n})
	return err
}

// Set printing area width
// nL, nH = (nL + nH x 256) x 0.125mm
func (p *Driver) SetPrintingAreaWidth(nL, nH uint8) error {
	_, err := p.rwc.Write([]byte{GS, 'W', nL, nH})
	return err
}

// Prints the data in the print buffer collectively
// and returns to standard mode.
func (p *Driver) PrintBufferAndReturnToStandardMode() error {
	_, err := p.rwc.Write([]byte{FF})
	return err
}

// When in page mode, all data in the print buffer is printed
// Command is only effective in page mode
// After printing, the printer does not delete the set value of
// ESC T and ESC W
func (p *Driver) PrintBufferInPageMode() error {
	_, err := p.rwc.Write([]byte{ESC, FF})
	return err
}

// Selects page mode
func (p *Driver) SelectPageMode() error {
	_, err := p.rwc.Write([]byte{ESC, 'L'})
	return err
}

// Selects standard mode
func (p *Driver) SelectStandardMode() error {
	_, err := p.rwc.Write([]byte{ESC, 'S'})
	return err
}

// Select print direction in page mode
// a: 0 <= a <= 3
// 0: left ro tight, starting upper left corner
// 1: bottom to top, starting lower left corner
// 2: right to left, starting lower right corner
// 3: top to bottom, starting upper right corner
func (p *Driver) SelectPrintDirectionInPageMode(a uint8) error {
	if a > 3 {
		return ErrInvalidPrintDirection
	}
	_, err := p.rwc.Write([]byte{ESC, 'T', a})
	return err
}

// Set print area in page mode
// xL, xH: Horizontal starting position
// yL, yH: Vertical starting position
// dxL, dxH: Horizontal printing area
// dyL, dyH: Vertical printing area
// x0 = ((xL + xH x 256) x 0.125mm)
// y0 = ((yL + yH x 256) x 0.125mm)
// dx = ((dxL + dxH x 256) x 0.125mm)
// dy = ((dyL + dyH x 256) x 0.125mm)
func (p *Driver) SetPrintAreaInPageMode(xL, xH, yL, yH, dxL, dxH, dyL, dyH uint8) error {
	// TODO: Handle error on dL, dH = 0
	_, err := p.rwc.Write([]byte{ESC, 'W', xL, xH, yL, yH, dxL, dxH, dyL, dyH})
	return err
}

// Set absolute vertical print position in page mode
// nL, nH = (nL + nH x 256) x 0.125mm
func (p *Driver) SetAbsoluteVerticalPrintPositionInPageMode(nL, nH uint8) error {
	_, err := p.rwc.Write([]byte{GS, DOLLAR, nL, nH})
	return err
}

// Set relative vertical print position in page mode
func (p *Driver) SetRelativeVerticalPrintPositionInPageMode(nL, nH uint8) error {
	_, err := p.rwc.Write([]byte{GS, BACKSLASH, nL, nH})
	return err
}
