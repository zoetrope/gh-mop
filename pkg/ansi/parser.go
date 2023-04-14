package ansi

import (
	"bytes"
	"regexp"
	"strconv"
	"strings"
)

var escapeSequenceRegex *regexp.Regexp

func init() {
	escapeSequenceRegex = regexp.MustCompile(`^([0-9A-Z>])|(]\d+;)|\[([?=]?)(\d+)?((;\d+)*)([a-zA-Z@])$`)
}

const (
	AsciiBell = 0x07 // Bell
	AsciiBS   = 0x08 // Backspace
	AsciiHT   = 0x09 // Horizontal Tab
	AsciiLF   = 0x0A // Line Feed
	AsciiVT   = 0x0B // Vertical Tab
	AsciiFF   = 0x0C // Form Feed
	AsciiCR   = 0x0D // Carriage Return
	AsciiEsc  = 0x1B // Escape
	AsciiDell = 0x7F // Delete
)

type Code string

const (
	Skip           = Code("Skip")
	Character      = Code("Character")
	MoveLeft       = Code("MoveLeft")
	Linefeed       = Code("Linefeed")
	CarriageReturn = Code("CarriageReturn")

	EraseLine   = Code("EraseLine")
	EraseScreen = Code("EraseScreen")
	InsertSpace = Code("InsertSpace")

	MoveRight = Code("MoveRight")
	Delete    = Code("Delete")
)

type parser struct {
	buffer   bytes.Buffer
	escaping bool
}

type control struct {
	code   Code
	params []int
}

func (t *parser) parse(r rune) control {
	if !t.escaping {
		switch r {
		case AsciiEsc:
			t.escaping = true
			return control{code: Skip}
		case AsciiBS:
			return control{code: MoveLeft}
		case AsciiCR:
			return control{code: CarriageReturn}
		case AsciiLF:
			fallthrough
		case AsciiBell:
			return control{code: Skip}
		}
		return control{code: Character}
	}

	t.buffer.WriteRune(r)
	currentString := t.buffer.String()
	groups := escapeSequenceRegex.FindStringSubmatch(currentString)
	if len(groups) == 0 {
		return control{code: Skip}
	}

	code := Skip
	if len(groups[1]) > 0 {
		if groups[1] == ">" {
			code = EraseScreen
		}
	} else if len(groups[7]) > 0 {
		switch groups[7] {
		case "K":
			code = EraseLine
		case "J":
			code = EraseScreen
		case "C":
			code = MoveRight
		case "P":
			code = Delete
		case "@":
			code = InsertSpace
		}
	}
	params := make([]int, 0)
	if len(groups[4]) > 0 {
		p, err := strconv.Atoi(groups[4])
		if err != nil {
			panic(err)
		}
		params = append(params, p)
	}
	if len(groups[5]) > 0 {
		for _, p := range strings.Split(groups[5], ";") {
			if len(p) > 0 {
				p, err := strconv.Atoi(p)
				if err != nil {
					panic(err)
				}
				params = append(params, p)
			}
		}
	}

	t.buffer.Reset()
	t.escaping = false

	return control{code: code, params: params}
}
