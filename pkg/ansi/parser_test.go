package ansi

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {

	parser := parser{
		buffer:   bytes.Buffer{},
		escaping: false,
	}

	//input := "\x1b[J"
	//input := "\x1b[2J"
	input := "abc\x1b[2Jdef\x1b[0mghi"
	expect := []control{
		{code: Skip},
		{code: Skip},
		{code: Skip},
		{code: EraseEntireScreen, params: []int{2}},
	}

	result := make([]control, 0)
	for _, i := range input {
		result = append(result, parser.parse(i))
	}
	fmt.Printf("result: %v\n", result)

	assert.Equal(t, len(expect), len(result))
	for i := range expect {
		assert.Equal(t, expect[i], result[i])
	}
}
