package ansi

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strings"
)

func ProcessFile(inputFile string) (string, error) {
	file, err := os.Open(inputFile)
	if err != nil {
		return "", err
	}
	defer file.Close()

	result := ""
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		result += processLine(line) + "\n"
	}

	if scanner.Err() != nil {
		return "", scanner.Err()
	}

	return result, nil
}

func processLine(line string) string {
	reader := strings.NewReader(line)
	result := make([]byte, 0, len(line))
	cur := 0
	parser := parser{
		buffer:   bytes.Buffer{},
		escaping: false,
	}

	for {
		r, _, err := reader.ReadRune()
		if err != nil {
			break
		}
		es := parser.parse(r)

		// for debug
		//if es.code != Skip {
		//	fmt.Printf("code: %v, params: %v, cur: %v, result: %v\n", es.code, es.params, cur, string(result))
		//}

		switch es.code {
		case Skip:
			break
		case MoveLeft:
			p := 1
			if len(es.params) != 0 {
				p = es.params[0]
			}
			cur -= p
			if cur < 0 {
				cur = 0
			}
		case MoveRight:
			p := 1
			if len(es.params) != 0 {
				p = es.params[0]
			}
			cur += p
			if cur > len(result) {
				cur = len(result)
			}
		case CarriageReturn:
			cur = 0
		case EraseLine:
			p := 0
			if len(es.params) != 0 {
				p = es.params[0]
			}
			switch p {
			case 0: // erase from cursor to end of line
				result = result[:cur]
				cur = len(result)
			case 1: // erase start of line to the cursor
				result = result[cur:]
				cur = 0
			case 2: // erase the entire line
				result = make([]byte, 0, len(line))
				cur = 0
			}
		case EraseScreen:
			p := 0
			if len(es.params) != 0 {
				p = es.params[0]
			}
			switch p {
			case 0: // erase from cursor to end of screen
				result = result[:cur]
				cur = len(result)
			case 1: // erase from cursor to beginning of screen
				result = result[cur:]
				cur = 0
			case 2: // erase entire screen
				result = make([]byte, 0, len(line))
				cur = 0
			case 3: // erase saved lines
				result = make([]byte, 0, len(line))
				cur = 0
			}
		case InsertSpace:
			p := 1
			if len(es.params) != 0 {
				p = es.params[0]
			}
			for i := 0; i < p; i++ {
				result = append(result[:cur+1], result[cur:]...)
				result[cur] = ' '
			}
		case Delete:
			p := 1
			if len(es.params) != 0 {
				p = es.params[0]
			}
			for i := 0; i < p; i++ {
				result = append(result[:cur], result[cur+1:]...)
			}
		case Character:
			if len(result) > cur {
				result[cur] = byte(r)
			} else {
				result = append(result, byte(r))
			}
			cur++
		default:

		}
	}
	if len(parser.buffer.String()) > 0 {
		fmt.Printf("unprocessed: %s\n", parser.buffer.String())
	}
	return string(result)
}
