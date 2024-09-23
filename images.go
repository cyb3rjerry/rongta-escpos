package rongta

import "errors"

var (
	ErrInvalidBitImageModevalue = errors.New("invalid m value")
)

// Select bit-image mode
//
// Selects a bit-image mode using m for the number of dots specified
// by nL and nH.
// m = 0, 1, 32, 33,
// m = 0: 8-dot single-density mode,
// m = 1: 8-dot double-density mode,
// m = 32: 24-dot single-density mode,
// m = 33: 24-dot double-density mode,
// 0 <= nH <= 3
func (p *Printer) SelectBitImageMode(m, nL, nH uint8, d []uint8) error {
	if m != 0 && m != 1 && m != 32 && m != 33 {
		return ErrInvalidBitImageModevalue
	}

	_, err := p.rwc.Write(append([]byte{ESC, '*', m, nL, nH}, d...))
	return err
}

// Print NV bit image
// m: bit image mode
// m = 0: Normal mode (vDensity = 203.2dpi, hDensity = 203.2dpi)
// m = 1: Double width mode (vDensity = 203.2dpi, hDensity = 101.6dpi)
// m = 2: Double height mode (vDensity = 101.6dpi, hDensity = 203.2dpi)
// m = 3: Quadruple mode (vDensity = 101.6dpi, hDensity = 101.6dpi)
// n = the number of the NV bit image (defined using the FS q command)
// This command is not effective when the specified NV bit image
// has not been defined
func (p *Printer) PrintNVBitImage(n, m uint8) error {
	_, err := p.rwc.Write([]byte{ESC, 'p', n, m})
	return err
}

// Define NV bit image
// n: specifies the number of the defined NV bit image
// xL, xH (xL + yH x 256) x 8 dots: specify the horizontal size of the bit image
// yL, yH (yL + yH x 256) x 8 dots: specify the vertical size of the bit image
//
// 1 <= n <= 255,
// 0 <= xL <= 255,
// 0 <= xH <= 3 (when 1 <= (xL + xH x 256) <= 1023),
// 0 <= yL <= 255,
// 0 <= yH <= 3 (when 1 <= (yL + yH x 256) <= 288),
// 0 <= d <= 255,
// Total defined data area = 192K bytes.
//
// Frequent write command executions may damage the NV
// memory. Therefore, it is recommended to write the NV memory
// 10 times or less a day.
func (p *Printer) DefineNVBitImage(n, xL, xH, yL, yH uint8, d []uint8) error {
	// TODO: Validate all of this
	panic("unimplemented")
	_, err := p.rwc.Write(append([]byte{ESC, 'q', n, xL, xH, yL, yH}, d...))
	return err
}

// Define downloaded bit images
// x specifies the number of dots in the horizontal direction
// y specifies the number of dots in the vertical direction
// d specifies the bit image data
//
// The downloaded bit image definition is cleared when:
// 1) ESC @ is executed.
// 2) ESC & is executed.
// 3) Printer is reset or the power is turned off.
func (p *Printer) DefineDownloadedBitImage(x, y uint8, d []uint8) error {
	panic("unimplemented")
	_, err := p.rwc.Write(append([]byte{GS, 'v', 0, x, y}, d...))
	return err
}

// Prints a downloaded bit image using the mode specified by
// m. m selects a mode from the table below:
// m = 0: Normal mode (vDensity = 203.2dpi, hDensity = 203.2dpi),
// m = 1: Double width mode (vDensity = 203.2dpi, hDensity = 101.6dpi),
// m = 2: Double height mode (vDensity = 101.6dpi, hDensity = 203.2dpi),
// m = 3: Quadruple mode (vDensity = 101.6dpi, hDensity = 101.6dpi)
func (p *Printer) PrintDownloadedBitImage(m uint8) error {
	panic("unimplemented")
	if m > 3 {
		return ErrInvalidBitImageModevalue
	}

	_, err := p.rwc.Write([]byte{GS, SLASH, m})
	return err
}

// Prints NV bit image n using the mode specified by m.
// m selects a mode from the table below:
// m = 0: Normal mode (vDensity = 203.2dpi, hDensity = 203.2dpi),
// m = 1: Double width mode (vDensity = 203.2dpi, hDensity = 101.6dpi),
// m = 2: Double height mode (vDensity = 101.6dpi, hDensity = 203.2dpi),
// m = 3: Quadruple mode (vDensity = 101.6dpi, hDensity = 101.6dpi)
func (p *Printer) PrintNVBitImageMode(n, m uint8) error {
	if m > 3 {
		return ErrInvalidBitImageModevalue
	}

	_, err := p.rwc.Write([]byte{GS, 'p', n, m})
	return err
}

// Print raster bit image
// m: bit image mode
// d: bit image data
// xL, xH: select the number of data bytes (xL + xH x 256) in the horizontal direction
// yL, yH: select the number of data bytes (yL + yH x 256) in the vertical direction
// m = 0: Normal mode (vDensity = 203.2dpi, hDensity = 203.2dpi),
// m = 1: Double width mode (vDensity = 203.2dpi, hDensity = 101.6dpi),
// m = 2: Double height mode (vDensity = 101.6dpi, hDensity = 203.2dpi),
// m = 3: Quadruple mode (vDensity = 101.6dpi, hDensity = 101.6dpi)
// 0 <= m <= 3, 0 <= d <= 255
// xL <= 255, yL <= 255
// 0 <= xH <= 255 where 1 <= (xL + xH x 256) <= 128
// 0 <= yH <= 8 where 1 <= (yL + yH x 256) <= 4095
func (p *Printer) PrintRasterBitImage(m, xL, xH, yL, yH uint8, d []uint8) error {
	if m > 3 {
		return ErrInvalidBitImageModevalue
	}

	_, err := p.rwc.Write(append([]byte{GS, 'v', '0', m, xL, xH, yL, yH}, d...))
	return err
}
