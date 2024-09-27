package commands

import "errors"

type CharacterCode uint8
type CharacterSet uint8

const (
	// Character code table
	CP437      CharacterCode = 0
	Katakana   CharacterCode = 1
	CP850      CharacterCode = 2
	CP860      CharacterCode = 3
	CP863      CharacterCode = 4
	CP865      CharacterCode = 5
	WPCP1251   CharacterCode = 6
	CP866      CharacterCode = 7
	MIK        CharacterCode = 8
	CP755      CharacterCode = 9
	Iran       CharacterCode = 10
	CP862      CharacterCode = 15
	WCP1252    CharacterCode = 16
	WCP1253    CharacterCode = 17
	CP852      CharacterCode = 18
	CP858      CharacterCode = 19
	Iran2      CharacterCode = 20
	Latvian    CharacterCode = 21
	CP864      CharacterCode = 22
	ISO8859_1  CharacterCode = 23
	CP737      CharacterCode = 24
	WCP1257    CharacterCode = 25
	Thai       CharacterCode = 26
	CP720      CharacterCode = 27
	CP855      CharacterCode = 28
	CP857      CharacterCode = 29
	WCP1250    CharacterCode = 30
	CP775      CharacterCode = 31
	WCP1254    CharacterCode = 32
	WCP1255    CharacterCode = 33
	WCP1256    CharacterCode = 34
	WCP1258    CharacterCode = 35
	ISO8859_2  CharacterCode = 36
	ISO8859_3  CharacterCode = 37
	ISO8859_4  CharacterCode = 38
	ISO8859_5  CharacterCode = 39
	ISO8859_6  CharacterCode = 40
	ISO8859_7  CharacterCode = 41
	ISO8859_8  CharacterCode = 42
	ISO8859_9  CharacterCode = 43
	ISO8859_15 CharacterCode = 44
	Thai2      CharacterCode = 45
	CP856      CharacterCode = 46
	CP874      CharacterCode = 47

	// Character sets
	USA          CharacterSet = 0
	France       CharacterSet = 1
	German       CharacterSet = 2
	UK           CharacterSet = 3
	DenmarkI     CharacterSet = 4
	Sweden       CharacterSet = 5
	Italy        CharacterSet = 6
	SpainI       CharacterSet = 7
	Japan        CharacterSet = 8
	Norway       CharacterSet = 9
	DenmarkII    CharacterSet = 10
	SpainII      CharacterSet = 11
	LatinAmerica CharacterSet = 12
	Korea        CharacterSet = 13
	Slovenia     CharacterSet = 14
	China        CharacterSet = 15
)

var (
	ErrInvalidCharacterCode             = errors.New("invalid character code")
	ErrInvalidInternationalCharacterSet = errors.New("invalid international character set")

	ErrInvalidCancelCharacterCode = errors.New("invalid character code being canceled")
)

// Select international character set
// n = 0, 1, 2, ..., 15
// 0: USA, 1: France, 2: Germany, 3: UK, 4: Denmark I, 5: Sweden, 6: Italy, 7: Spain I,
// 8: Japan, 9: Norway, 10: Denmark II, 11: Spain II, 12: Latin America, 13: Korea, 14: Slovenia,
// 15: China
func (p *Driver) SelectInternationalCharacterSet(c CharacterSet) error {
	_, err := p.rwc.Write([]byte{ESC, 'R', byte(c)})
	return err
}

// Select international character code
func (p *Driver) SelectInternationalCharacterCode(n CharacterCode) error {
	_, err := p.rwc.Write([]byte{ESC, 'R', uint8(n)})
	return err
}

// Cancel user-defined characters
// 32 <= n <= 126
func (p *Driver) CancelUserDefinedCharacters(n uint8) error {
	if n < 32 || n > 126 {
		return ErrInvalidCancelCharacterCode
	}

	_, err := p.rwc.Write([]byte{ESC, '-', n})
	return err
}
