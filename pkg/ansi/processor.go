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

		switch es.code {
		case Skip:
			break
		case Backspace:
			if cur > 0 {
				cur--
			}
		//case LINEFEED:
		//	result = make([]byte, 0, len(line))
		//	cur = 0
		case CarriageReturn:
			if len(result) > cur {
				result[cur] = '\n'
			} else {
				result = append(result, '\n')
			}
			cur++
		case EraseLine:
			result = result[:cur]
		case EraseEntireScreen:
			result = make([]byte, 0, len(line))
			cur = 0
		case InsertSpace:
			result = append(result[:cur+1], result[cur:]...)
			result[cur] = ' '
		case MoveRight:
			cur++
		case Delete:
			result = append(result[:cur], result[cur+1:]...)
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
