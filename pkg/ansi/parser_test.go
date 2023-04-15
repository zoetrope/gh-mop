package ansi

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		expect []control
	}{
		{
			name:  "erase screen",
			input: "\x1b[2J",
			expect: []control{
				{code: EraseScreen, params: []int{2}},
			},
		},
		{
			name:  "remove 4 characters",
			input: "\b\x1b[K\b\x1b[K\b\x1b[K\b\x1b[K",
			expect: []control{
				{code: MoveLeft},
				{code: EraseLine, params: []int{}},
				{code: MoveLeft},
				{code: EraseLine, params: []int{}},
				{code: MoveLeft},
				{code: EraseLine, params: []int{}},
				{code: MoveLeft},
				{code: EraseLine, params: []int{}},
			},
		},
		{
			name:  "characters",
			input: "abc\x1b[Kdef\x1b[0mghi",
			expect: []control{
				{code: Character}, // a
				{code: Character}, // b
				{code: Character}, // c
				{code: EraseLine, params: []int{}},
				{code: Character}, // d
				{code: Character}, // e
				{code: Character}, // f
				// skip \x1b[0m
				{code: Character}, // g
				{code: Character}, // h
				{code: Character}, // i
			},
		},
		{
			name:  "multiple parameters",
			input: "\x1b[10;20;30@",
			expect: []control{
				{code: InsertSpace, params: []int{10, 20, 30}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := parser{
				buffer:   bytes.Buffer{},
				escaping: false,
			}
			result := make([]control, 0)
			for _, i := range tt.input {
				c := parser.parse(i)
				if c.code != Skip {
					result = append(result, c)
				}
			}
			require.Equal(t, len(tt.expect), len(result))
			for i := range tt.expect {
				assert.Equal(t, tt.expect[i], result[i])
			}
		})
	}
}
