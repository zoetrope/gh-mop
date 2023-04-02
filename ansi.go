package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"regexp"
	"strings"
)

const AnsiEscapeSequenceRegex = `\x1B([87MH>])|\[(\?\d+[hl])|([0-2]?[KJ])|(\d*[ABCDDEFGgimnPSsTu@])|(1000D\d+)|(\d*;\d*[fHrm])|(\d+;\d+;\d+m)|(=\d+h)`

//                                   (1         (2          (3           (4                        (5         (6              (7
// 1. \x1B8 - cursor control
// 2. \x1B[?25l - common private mode
// 3. \x1B[2J - erase
// 4. \x1B[1C - cursor control
// 5. \x1B[1000D - move cursor to beginning of line
// 6. \x1B[1;1H - move cursor to position
// 7. \x1B[1;1m - set color
// 8. \x1B[=1h - set screen mode

func main() {
	inputFile := "./core/testdata/input.txt"

	err := processFile(inputFile)
	if err != nil {
		fmt.Println("Error:", err)
	}
}

func processFile(inputFile string) error {
	file, err := os.Open(inputFile)
	if err != nil {
		return err
	}
	defer file.Close()

	ansiRegex, err := regexp.Compile(AnsiEscapeSequenceRegex)
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		processLine(line, ansiRegex)
	}

	if scanner.Err() != nil {
		return scanner.Err()
	}

	return nil
}

const (
	// General ASCII Codes
	BEL = 0x07 // Bell
	BS  = 0x08 // Backspace
	HT  = 0x09 // Horizontal Tab
	LF  = 0x0A // Line Feed
	VT  = 0x0B // Vertical Tab
	FF  = 0x0C // Form Feed
	CR  = 0x0D // Carriage Return
	ESC = 0x1B // Escape
	DEL = 0x7F // Delete
)

const (
	SKIP = iota
	CHARACTER
	BACKSPACE
	LINEFEED

	CLEAR
	INSERT_SPACE

	MOVE_RIGHT
	DELETE
)

type termInfo struct {
	buffer    bytes.Buffer
	escape    bool
	ansiRegex *regexp.Regexp
}

func (t *termInfo) parse(r rune) int {
	if t.escape {
		t.buffer.WriteRune(r)
		currentString := t.buffer.String()
		groups := t.ansiRegex.FindStringSubmatch(currentString)
		//match := ansiRegex.FindString(currentString)
		if len(groups) > 0 {
			code := SKIP
			if len(groups[1]) > 0 {

			} else if len(groups[2]) > 0 {

			} else if len(groups[3]) > 0 {
				code = CLEAR
			} else if len(groups[4]) > 0 {
				if strings.HasSuffix(groups[4], "@") {
					code = INSERT_SPACE
				}
				if strings.HasSuffix(groups[4], "C") {
					code = MOVE_RIGHT
				}
				if strings.HasSuffix(groups[4], "P") {
					code = DELETE
				}
			}

			t.buffer.Reset()
			t.escape = false
			return code
		}
		return SKIP
	}
	switch r {
	case ESC:
		t.escape = true
		return SKIP
	case BS:
		return BACKSPACE
	case LF:
		return LINEFEED
	}

	return CHARACTER
}

func processLine(line string, ansiRegex *regexp.Regexp) {
	reader := strings.NewReader(line)
	result := make([]byte, 0, len(line))
	cur := 0
	term := termInfo{
		buffer:    bytes.Buffer{},
		escape:    false,
		ansiRegex: ansiRegex,
	}

	for {
		r, _, err := reader.ReadRune()
		if err != nil {
			break
		}
		code := term.parse(r)

		switch code {
		case SKIP:
			break
		case BACKSPACE:
			if cur > 0 {
				cur--
			}
		case LINEFEED:
			result = make([]byte, 0, len(line))
			cur = 0
		case CLEAR:
			result = result[:cur]
		case INSERT_SPACE:
			result = append(result[:cur+1], result[cur:]...)
			result[cur] = ' '
		case MOVE_RIGHT:
			cur++
		case DELETE:
			result = append(result[:cur], result[cur+1:]...)
		case CHARACTER:
			if len(result) > cur {
				result[cur] = byte(r)
			} else {
				result = append(result, byte(r))
			}
			cur++
		default:

		}
	}
	fmt.Printf("%s\n", string(result))
}
